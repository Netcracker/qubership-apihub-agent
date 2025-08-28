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

package service

import (
	"strings"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/client"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	"github.com/Netcracker/qubership-apihub-agent/view"
	log "github.com/sirupsen/logrus"
)

type RegistrationService interface {
	RunAgentRegistrationProcess()
}

func NewRegistrationService(cloudName string, namespace string, agentUrl string, backendVersion string, agentName string, client client.AgentsBackendClient, disablingService DisablingService) RegistrationService {
	return &registrationServiceImpl{cloudName: cloudName, namespace: namespace, agentUrl: agentUrl, backendVersion: backendVersion, agentName: agentName, client: client, disablingService: disablingService}
}

type registrationServiceImpl struct {
	cloudName      string
	namespace      string
	agentUrl       string
	backendVersion string
	agentName      string

	client           client.AgentsBackendClient
	disablingService DisablingService
}

const AGENT_VERSION = "1.0.0"

func (r registrationServiceImpl) RunAgentRegistrationProcess() {
	utils.SafeAsync(func() {
		req := view.AgentKeepaliveMessage{
			AgentDeploymentCloud:     r.cloudName,
			AgentDeploymentNamespace: r.namespace,
			AgentUrl:                 r.agentUrl,
			BackendVersion:           r.backendVersion,
			AgentName:                r.agentName,
			AgentVersion:             AGENT_VERSION,
		}
		for range time.Tick(time.Second * 10) {
			utils.SafeAsync(func() {
				version, err := r.client.SendKeepaliveMessage(req)
				if err != nil {
					log.Errorf("Failed to send registration message: %s", err)
					return
				}

				r.disablingService.DisableServices(!isAgentCompatibleWithAPIHUB(version), version)
			})
		}
	})
}

func isAgentCompatibleWithAPIHUB(apihubAgentVersion string) bool {
	if AGENT_VERSION == apihubAgentVersion {
		return true
	}
	if apihubAgentVersion == "" {
		log.Errorf("Current version %s of Agent is incompatible with APIHUB (APIHUB is probably outdated). Please, contact your System Administrator to update this Agent instance to version %s.", AGENT_VERSION, apihubAgentVersion)
		return false
	}
	agentVersion := strings.Split(AGENT_VERSION, ".")
	apihubVersion := strings.Split(apihubAgentVersion, ".")
	if agentVersion[0] != apihubVersion[0] {
		log.Errorf("Current version %s of Agent is incompatible with APIHUB. Please, contact your System Administrator to update this Agent instance to version %s.", AGENT_VERSION, apihubAgentVersion)
		return false
	}
	if agentVersion[1] != apihubVersion[1] || agentVersion[2] != apihubVersion[2] {
		log.Warnf("Difference in minor/patch version of Agent detected. We recommend to contact your System Administrator to update this Agent instance to version %s.", apihubAgentVersion)
		return true
	}
	return true
}
