package mcp

import (
	"fmt"
	"sync"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/api_type/generic"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
	log "github.com/sirupsen/logrus"
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

	/*
		Url      string
		XApiKind string
		Name     string
		ApiType  ApiType
		Required bool
		Timeout  time.Duration
	*/

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
