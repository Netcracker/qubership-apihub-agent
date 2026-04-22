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

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

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

	// Clear messages before the specific call to only capture the relevant ones
	// Wait, initialize already happened during Connect.
	// If documentName is "init", we want the messages from Connect.
	var relevantMessages []jsonrpc.Message

	if documentName == "init" {
		capTransport.mu.Lock()
		relevantMessages = append(relevantMessages, capTransport.messages...)
		capTransport.mu.Unlock()
	} else {
		// Clear messages so we only get the ones for the specific call
		capTransport.mu.Lock()
		capTransport.messages = make([]jsonrpc.Message, 0)
		capTransport.mu.Unlock()

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

		// Wait a tiny bit for any trailing messages? No, the response is synchronous with the method return.
		capTransport.mu.Lock()
		relevantMessages = append(relevantMessages, capTransport.messages...)
		capTransport.mu.Unlock()
	}

	// We want to return the raw JSON-RPC response message.
	// The response message is the last one in relevantMessages that is a response.
	var responseMsg jsonrpc.Message
	for i := len(relevantMessages) - 1; i >= 0; i-- {
		msg := relevantMessages[i]
		if msg == nil {
			continue
		}
		
		// Marshal to map to check if it's a response
		b, err := json.Marshal(msg)
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

	bytes, err := json.MarshalIndent(responseMsg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	return bytes, nil
}
