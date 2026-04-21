package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/api_type/generic"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
	log "github.com/sirupsen/logrus"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gosimple/slug"
)

func NewMcpDiscoveryRunner() generic.DiscoveryRunner {
	return &mcpDiscoveryRunner{}
}

type mcpDiscoveryRunner struct {
}

func (r mcpDiscoveryRunner) DiscoverDocuments(baseUrl string, urls view.DocumentDiscoveryUrls, timeout time.Duration) ([]view.Document, []view.EndpointCallInfo, error) {
	// connect to streamable http endpoint
	// send initialization request and get response
	// collect capabilities and generate fake documents based on capabilities, probably need to send specific requests like get tools
	// (extra document for init response itself)

	// include capability type to name
	// most probably it should be like mcp_tools_api-v1-mcp, mcp_prompts_api-v1-mcp, etc...

	refs := utils.MakeDocumentRefsFromUrls(urls.Mcp, view.ATMcp, false, timeout)
	return r.GetDocumentsByRefs(baseUrl, refs, "")

}

type LoggingTransport struct {
	mcp.Transport
}

func (lt LoggingTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	return lt.Transport.Connect(ctx)
}

func (r mcpDiscoveryRunner) GetDocumentsByRefs(baseUrl string, refs []view.DocumentRef, configPath string) ([]view.Document, []view.EndpointCallInfo, error) {
	// TODO: not sure if applicable to MCP in real life

	filteredRefs := r.FilterRefsForApiType(refs) // take only appropriate api type
	if len(filteredRefs) == 0 {
		return nil, nil, nil
	}

	result := make([]view.Document, len(filteredRefs))
	failedCalls := make([]view.EndpointCallInfo, len(filteredRefs))
	errors := make([]string, len(filteredRefs))

	// for MCP the implementation is different: connect to an endpoint(could be multiple MCP endpoints!!)

	for i, ref := range filteredRefs {
		ctx, cancelFunc := context.WithTimeout(context.Background(), ref.Timeout)

		// Create a new client, with no features.
		client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

		transport := &mcp.StreamableClientTransport{
			Endpoint: ref.Url,
		}

		lt := LoggingTransport{
			Transport: transport,
		}

		session, err := client.Connect(ctx, lt, nil)
		if err != nil {
			failedCalls[i] = view.EndpointCallInfo{
				Path:         ref.Url,
				ErrorSummary: fmt.Sprintf("Failed to connect to MCP endpoint: %s", err.Error()),
			}
			if ref.Required {
				errors[i] = fmt.Sprintf("Failed to connect to required MCP endpoint %s: %s", ref.Url, err.Error())
			}
			cancelFunc()
			continue
		}

		urlSlug := slug.Make(ref.Url) // TODO: use only path from the URL!!!

		var serverCaps *mcp.ServerCapabilities
		if ir := session.InitializeResult(); ir != nil {
			serverCaps = ir.Capabilities
		}
		if serverCaps != nil && serverCaps.Tools != nil {
			toolsRes, err := session.ListTools(ctx, nil) // TODO: handle paging(cursor)
			if err != nil {
				failedCalls[i] = view.EndpointCallInfo{
					Path:         ref.Url,
					ErrorSummary: fmt.Sprintf("Failed to list tools from MCP endpoint: %s", err.Error()),
				}
				if ref.Required {
					errors[i] = fmt.Sprintf("Failed to list tools from required MCP endpoint %s: %s", ref.Url, err.Error())
				}
				cancelFunc()
				continue
			}
			if len(toolsRes.Tools) > 0 {
				result = append(result, view.Document{
					Name:       "tools",
					Format:     "json",
					FileId:     "tools_" + urlSlug + ".json",
					Type:       view.McpType,
					XApiKind:   view.UnknownType,
					DocPath:    ref.Url,
					ConfigPath: "",
				})
			}
		}
		if serverCaps != nil && serverCaps.Resources != nil {
			resourcesRes, err := session.ListResources(ctx, nil) // TODO: handle paging(cursor)
			if err != nil {
				failedCalls[i] = view.EndpointCallInfo{
					Path:         ref.Url,
					ErrorSummary: fmt.Sprintf("Failed to list resources from MCP endpoint: %s", err.Error()),
				}
				if ref.Required {
					errors[i] = fmt.Sprintf("Failed to list resources from required MCP endpoint %s: %s", ref.Url, err.Error())
				}
				if err = session.Close(); err != nil {
					log.Warnf("Failed to close MCP session for endpoint %s: %s", ref.Url, err.Error())
				}
				cancelFunc()
				continue
			}
			if len(resourcesRes.Resources) > 0 {
				result = append(result, view.Document{
					Name:       "resources",
					Format:     "json",
					FileId:     "resources_" + urlSlug + ".json",
					Type:       view.McpType,
					XApiKind:   view.UnknownType,
					DocPath:    ref.Url,
					ConfigPath: "",
				})
			}
		}
		if serverCaps != nil && serverCaps.Prompts != nil {
			promptsRes, err := session.ListPrompts(ctx, nil) // TODO: handle paging(cursor)
			if err != nil {
				failedCalls[i] = view.EndpointCallInfo{
					Path:         ref.Url,
					ErrorSummary: fmt.Sprintf("Failed to list prompts from MCP endpoint: %s", err.Error()),
				}
				if ref.Required {
					errors[i] = fmt.Sprintf("Failed to list prompts from required MCP endpoint %s: %s", ref.Url, err.Error())
				}
				if err = session.Close(); err != nil {
					log.Warnf("Failed to close MCP session for endpoint %s: %s", ref.Url, err.Error())
				}
				cancelFunc()
				continue
			}
			if len(promptsRes.Prompts) > 0 {
				result = append(result, view.Document{
					Name:       "prompts",
					Format:     "json",
					FileId:     "prompts_" + urlSlug + ".json",
					Type:       view.McpType,
					XApiKind:   view.UnknownType,
					DocPath:    ref.Url,
					ConfigPath: "",
				})
			}
		}

		if err = session.Close(); err != nil {
			log.Warnf("Failed to close MCP session for endpoint %s: %s", ref.Url, err.Error())
		}
		cancelFunc()
	}

	// 1. iterate over refs to detect different MCP endpoints
	// 2 connect to MCP service
	// 3. send init anyway + initialized notification
	// 4. iterate over capabilities and collect raw responses along with server init response.
	// same as above ...

	return utils.FilterResultDocuments(result), utils.FilterFailedEndpointCalls(failedCalls), utils.FilterResultErrors(errors)
}

func (m mcpDiscoveryRunner) FilterRefsForApiType(refs []view.DocumentRef) []view.DocumentRef {
	return utils.FilterRefsForApiType(refs, view.ATMcp)
}

func (m mcpDiscoveryRunner) GetName() string {
	return "mcp"
}
