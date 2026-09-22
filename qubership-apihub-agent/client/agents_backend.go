package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/secctx"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
	"gopkg.in/resty.v1"
)

const keepaliveHTTPTimeout = time.Second * 60

type AgentsBackendClient interface {
	SendKeepaliveMessage(pathPrefix string, msg view.AgentKeepaliveMessage) (string, error)
}

func NewAgentsBackendClient(apihubUrl string, accessToken string) (AgentsBackendClient, error) {
	httpClient, err := utils.CreateSecureHTTPClientCloseAfterRequest(keepaliveHTTPTimeout)
	if err != nil {
		return nil, err
	}
	restyClient := resty.NewWithClient(httpClient)
	restyClient.SetPreRequestHook(func(_ *resty.Client, req *resty.Request) error {
		if req.RawRequest != nil {
			req.RawRequest.Close = true
		}
		return nil
	})
	return &agentsBackendClientImpl{
		apihubUrl:   apihubUrl,
		accessToken: accessToken,
		httpClient:  httpClient,
		restyClient: restyClient,
	}, nil
}

type agentsBackendClientImpl struct {
	apihubUrl   string
	accessToken string
	httpClient  *http.Client
	restyClient *resty.Client
}

func (a agentsBackendClientImpl) SendKeepaliveMessage(pathPrefix string, msg view.AgentKeepaliveMessage) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), keepaliveHTTPTimeout)
	defer cancel()

	req := a.makeRequest(secctx.CreateSystemContext())
	req.SetContext(ctx)
	req.SetHeader("Connection", "close")
	req.SetBody(msg)

	resp, err := req.Post(fmt.Sprintf("%s/%s/api/v2/agents", a.apihubUrl, pathPrefix))
	if err != nil {
		closeRestyResponseBody(resp)
		a.httpClient.CloseIdleConnections()
		return "", err
	}
	if resp.StatusCode() != http.StatusOK {
		if authErr := checkUnauthorized(resp); authErr != nil {
			return "", authErr
		}
		return "", fmt.Errorf("failed to send registration message with error code %d", resp.StatusCode())
	}
	body := resp.Body()
	if len(body) > 0 {
		type agentVersion struct {
			Version string `json:"version"`
		}
		var version agentVersion
		err = json.Unmarshal(body, &version)
		if err != nil {
			return "", err
		}
		return version.Version, nil
	}

	return "", nil
}

func (a agentsBackendClientImpl) makeRequest(ctx secctx.SecurityContext) *resty.Request {
	req := a.restyClient.R()
	if ctx.GetUserToken() != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", ctx.GetUserToken()))
	} else {
		req.SetHeader("api-key", a.accessToken)
	}
	return req
}

func closeRestyResponseBody(resp *resty.Response) {
	if resp == nil || resp.RawResponse == nil || resp.RawResponse.Body == nil {
		return
	}
	_ = resp.RawResponse.Body.Close()
}
