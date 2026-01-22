// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package asyncapi

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/api_type/generic"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
	log "github.com/sirupsen/logrus"
)

const DefaultAsyncAPISpecName = "AsyncAPI specification"

var asyncapi3Regexp = regexp.MustCompile(`^3\.0`)

func NewAsyncAPIDiscoveryRunner() generic.DiscoveryRunner {
	return &asyncAPIDiscoveryRunner{}
}

type asyncAPIDiscoveryRunner struct {
}

func (r asyncAPIDiscoveryRunner) DiscoverDocuments(baseUrl string, urls view.DocumentDiscoveryUrls, timeout time.Duration) ([]view.Document, error) {
	refs := utils.MakeDocumentRefsFromUrls(urls.AsyncAPI, view.ATAsyncAPI, false, timeout)
	return r.GetDocumentsByRefs(baseUrl, refs)
}

func (r asyncAPIDiscoveryRunner) GetDocumentsByRefs(baseUrl string, refs []view.DocumentRef) ([]view.Document, error) {
	filteredRefs := r.FilterRefsForApiType(refs) // take only appropriate api type
	if len(filteredRefs) == 0 {
		return nil, nil
	}

	result := make([]view.Document, len(filteredRefs))
	errors := make([]string, len(filteredRefs))

	wg := sync.WaitGroup{}
	wg.Add(len(filteredRefs))

	fileIds := sync.Map{}

	for it, ref := range filteredRefs {
		i := it
		currentSpecRef := ref
		currentSpecUrl := ref.Url

		utils.SafeAsync(func() {
			defer wg.Done()

			spec := view.Document{
				Path: currentSpecUrl,
			}

			url := baseUrl + currentSpecUrl

			specVersion, specTitle, specFormat, err := getAsyncAPISpecInfo(url, ref.Timeout)
			if err != nil {
				log.Debugf("Failed to read asyncapi spec from %s: %s", url, err)
				return
			}
			log.Debugf("Got valid asyncapi spec from: %v", url)
			if currentSpecRef.Name != "" {
				spec.Name = currentSpecRef.Name
			} else if specTitle != "" {
				spec.Name = specTitle
			} else {
				spec.Name = DefaultAsyncAPISpecName
			}
			spec.Format = specFormat
			spec.FileId = utils.GenerateFileId(&fileIds, spec.Name, specFormat)
			spec.Type = specVersion

			spec.XApiKind = currentSpecRef.XApiKind

			result[i] = spec
		})
	}

	wg.Wait()

	return utils.FilterResultDocuments(result), utils.FilterResultErrors(errors)
}

func getAsyncAPISpecInfo(specUrl string, timeout time.Duration) (string, string, string, error) {
	var spec view.JsonMap
	spec, specFormat, err := generic.GetGenericObjectFromUrl(specUrl, timeout)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get specification from '%v': %v", specUrl, err.Error())
	}
	if spec == nil {
		return "", "", "", fmt.Errorf("specification from '%v' is invalid", specUrl)
	}

	asyncapiVersion := spec.GetValueAsString("asyncapi")
	if asyncapiVersion == "" {
		return "", "", "", fmt.Errorf("not an asyncapi spec at '%v': missing 'asyncapi' field", specUrl)
	}

	infoObject := spec.GetObject("info")
	title := infoObject.GetValueAsString("title")
	version := infoObject.GetValueAsString("version")
	specTitle := title
	if version != "" {
		specTitle = title + " " + version
	}

	if asyncapi3Regexp.MatchString(asyncapiVersion) {
		return view.AsyncAPI30Type, specTitle, specFormat, nil
	}

	return "", "", "", fmt.Errorf("failed to determine asyncapi version from spec at '%v': version '%s'", specUrl, asyncapiVersion)
}

func (r asyncAPIDiscoveryRunner) FilterRefsForApiType(refs []view.DocumentRef) []view.DocumentRef {
	return utils.FilterRefsForApiType(refs, view.ATAsyncAPI)
}

func (r asyncAPIDiscoveryRunner) GetName() string {
	return "asyncapi"
}
