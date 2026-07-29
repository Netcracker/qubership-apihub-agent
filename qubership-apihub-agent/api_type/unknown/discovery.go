package unknown

import (
	"context"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/api_type/generic"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
)

func NewUnknownDiscoveryRunner() generic.DiscoveryRunner {
	return &unknownDiscoveryRunner{}
}

type unknownDiscoveryRunner struct {
}

func (m unknownDiscoveryRunner) DiscoverDocuments(ctx context.Context, baseUrl string, urls view.DocumentDiscoveryUrls, timeout time.Duration) ([]view.Document, []view.EndpointCallInfo, error) {
	// No default paths for this type
	return []view.Document{}, nil, nil
}

func (m unknownDiscoveryRunner) GetDocumentsByRefs(ctx context.Context, baseUrl string, refs []view.DocumentRef, configPath string) ([]view.Document, []view.EndpointCallInfo, error) {
	return generic.GetAnyDocsByRefs(ctx, baseUrl, m.FilterRefsForApiType(refs), configPath)
}

func (m unknownDiscoveryRunner) FilterRefsForApiType(refs []view.DocumentRef) []view.DocumentRef {
	return utils.FilterRefsForApiType(refs, view.ATUnknown)
}

func (m unknownDiscoveryRunner) GetName() string {
	return "unknown"
}
