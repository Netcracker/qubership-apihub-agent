package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type authTransport struct {
	key  string
	Base http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("api-key", t.key) // FIXME: just for testing

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clonedReq)
}

type captureTransport struct {
	base     mcp.Transport
	messages []jsonrpc.Message
	mu       sync.Mutex
}

func (t *captureTransport) Connect(ctx context.Context) (mcp.Connection, error) {
	conn, err := t.base.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &captureConnection{base: conn, transport: t}, nil
}

type captureConnection struct {
	base      mcp.Connection
	transport *captureTransport
}

func (c *captureConnection) Read(ctx context.Context) (jsonrpc.Message, error) {
	msg, err := c.base.Read(ctx)
	if err == nil {
		c.transport.mu.Lock()
		c.transport.messages = append(c.transport.messages, msg)
		c.transport.mu.Unlock()
	}
	return msg, err
}

func (c *captureConnection) Write(ctx context.Context, msg jsonrpc.Message) error {
	c.transport.mu.Lock()
	c.transport.messages = append(c.transport.messages, msg)
	c.transport.mu.Unlock()
	return c.base.Write(ctx, msg)
}

func (c *captureConnection) Close() error {
	return c.base.Close()
}

func (c *captureConnection) SessionID() string {
	return c.base.SessionID()
}

func GetMcpDocumentFromUrl(url string, documentName string, timeout time.Duration) ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	client := mcp.NewClient(&mcp.Implementation{Name: "apihub-agent", Version: "v0.0.1"}, nil)

	httpTransport := &mcp.StreamableClientTransport{
		Endpoint: url,
		HTTPClient: &http.Client{
			Transport: &authTransport{
				key: "ecd8a844aaa99e37d2a7c7c0432b4a755cf583ffc93f02a37b737b25661932d7", // FIXME: just for testing
			},
		},
	}

	capTransport := &captureTransport{
		base:     httpTransport,
		messages: make([]jsonrpc.Message, 0),
	}

	session, err := client.Connect(ctx, capTransport, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP endpoint: %w", err)
	}
	defer session.Close()

	var relevantMessages []jsonrpc.Message
	if documentName == "init" {
		// If documentName is "init", we need the messages from the connect phase and no need to send others
		capTransport.mu.Lock()
		relevantMessages = append(relevantMessages, capTransport.messages...)
		capTransport.mu.Unlock()
	} else {
		// Clear messages so we only get the ones for the specific call
		capTransport.mu.Lock()
		capTransport.messages = make([]jsonrpc.Message, 0)
		capTransport.mu.Unlock()

		// TODO: support paging?
		switch documentName {
		case "tools":
			_, err = session.ListTools(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to list tools: %w", err)
			}
		case "resources":
			_, err = session.ListResources(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to list resources: %w", err)
			}
		case "prompts":
			_, err = session.ListPrompts(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to list prompts: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown MCP document name: %s", documentName)
		}

		capTransport.mu.Lock()
		relevantMessages = append(relevantMessages, capTransport.messages...)
		capTransport.mu.Unlock()
	}

	// TODO: maybe better to start from the beginning: find the request and the response

	// We want to return the raw JSON-RPC response message.
	// The response message is the last one in relevantMessages that is a response.
	var responseMsg jsonrpc.Message
	for i := len(relevantMessages) - 1; i >= 0; i-- {
		msg := relevantMessages[i]
		if msg == nil {
			continue
		}

		// Encode to wire format to check if it's a response
		b, err := jsonrpc.EncodeMessage(msg)
		if err == nil {
			var m map[string]interface{}
			if err := json.Unmarshal(b, &m); err == nil {
				if _, hasResult := m["result"]; hasResult {
					responseMsg = msg
					break
				}
				if _, hasError := m["error"]; hasError {
					responseMsg = msg
					break
				}
			}
		}
	}

	if responseMsg == nil {
		return nil, fmt.Errorf("no JSON-RPC response captured for %s", documentName)
	}

	bytes, err := jsonrpc.EncodeMessage(responseMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode MCP result: %w", err)
	}

	// Format the JSON nicely
	var prettyJSON map[string]interface{}
	err = json.Unmarshal(bytes, &prettyJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal MCP response message: %w. Message: %s", err, string(bytes))
	}
	// TODO: process json and extract "result" value

	if prettyBytes, err := json.MarshalIndent(prettyJSON, "", "  "); err == nil {
		return prettyBytes, nil
	}

	return bytes, nil
}
