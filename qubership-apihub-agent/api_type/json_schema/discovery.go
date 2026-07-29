package json_schema

import (
	"context"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/api_type/generic"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
)

func NewJsonSchemaDiscoveryRunner() generic.DiscoveryRunner {
	return &jsonSchemaDiscoveryRunner{}
}

type jsonSchemaDiscoveryRunner struct {
}

func (j jsonSchemaDiscoveryRunner) DiscoverDocuments(ctx context.Context, baseUrl string, urls view.DocumentDiscoveryUrls, timeout time.Duration) ([]view.Document, []view.EndpointCallInfo, error) {
	// No default paths for this type
	return []view.Document{}, nil, nil
}

func (j jsonSchemaDiscoveryRunner) GetDocumentsByRefs(ctx context.Context, baseUrl string, refs []view.DocumentRef, configPath string) ([]view.Document, []view.EndpointCallInfo, error) {
	return generic.GetAnyDocsByRefs(ctx, baseUrl, j.FilterRefsForApiType(refs), configPath)
}

func (j jsonSchemaDiscoveryRunner) FilterRefsForApiType(refs []view.DocumentRef) []view.DocumentRef {
	return utils.FilterRefsForApiType(refs, view.ATJsonSchema)
}

func (j jsonSchemaDiscoveryRunner) GetName() string {
	return "json-schema"
}
