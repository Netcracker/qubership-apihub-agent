package unknown

import (
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

func (m unknownDiscoveryRunner) DiscoverDocuments(baseUrl string, urls view.DocumentDiscoveryUrls, timeout time.Duration) ([]view.Document, error) {
	// No default paths for this type
	return []view.Document{}, nil
}

func (m unknownDiscoveryRunner) GetDocumentsByRefs(baseUrl string, refs []view.DocumentRef) ([]view.Document, error) {
	return generic.GetAnyDocsByRefs(baseUrl, m.FilterRefsForApiType(refs))
}

func (m unknownDiscoveryRunner) FilterRefsForApiType(refs []view.DocumentRef) []view.DocumentRef {
	return utils.FilterRefsForApiType(refs, view.ATUnknown)
}

func (m unknownDiscoveryRunner) GetName() string {
	return "unknown"
}
