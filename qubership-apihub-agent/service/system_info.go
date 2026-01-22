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
	"fmt"
	"os"
	"time"

	"github.com/Netcracker/qubership-apihub-agent/config"
	"github.com/Netcracker/qubership-apihub-agent/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SystemInfoService interface {
	GetBackendVersion() string
	GetApihubUrl() string
	GetAgentUrl() string
	GetAccessToken() string
	GetCloudName() string
	GetAgentNamespace() string
	GetExcludeLabels() []string
	GetGroupingLabels() []string
	GetAgentName() string
	GetDiscoveryTimeout() time.Duration
	InsecureProxyEnabled() bool //TODO: remove this after deprecated proxy path is removed
	GetBasePath() string
	GetPaasPlatform() string
	GetDiscoveryUrls() config.ApiTypeUrlsConfig
}

func NewSystemInfoService() (SystemInfoService, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(getConfigFolder())
	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	utils.PrintConfig(cfg)
	if err := utils.ValidateConfig(cfg); err != nil {
		return nil, err
	}
	return &systemInfoServiceImpl{
		config: cfg}, nil
}

func setDefaults() {
	viper.SetDefault("technicalParameters.basePath", ".")
	viper.SetDefault("technicalParameters.listenAddress", ":8080")
	viper.SetDefault("technicalParameters.version", "unknown")
	viper.SetDefault("technicalParameters.apihub.url", "http://localhost:8090")
	viper.SetDefault("technicalParameters.cloudName", "unknown")
	viper.SetDefault("technicalParameters.namespace", "unknown")
	viper.SetDefault("technicalParameters.paasPlatform", "KUBERNETES")
	viper.SetDefault("security.allowedOrigins", []string{})
	viper.SetDefault("security.insecureProxy", false)
	viper.SetDefault("discovery.excludeLabels", []string{})
	viper.SetDefault("discovery.groupingLabels", []string{})
	viper.SetDefault("discovery.timeoutSec", 15)
	viper.SetDefault("discovery.urls.openapi.config-urls", []string{"/v3/api-docs/swagger-config", "/swagger-resources"})
	viper.SetDefault("discovery.urls.openapi.doc-urls", []string{"/q/openapi?format=json", "/v3/api-docs?format=json", "/v2/api-docs", "/swagger-ui/swagger.json"})
	viper.SetDefault("discovery.urls.apihub-config.config-urls", []string{"/v3/api-docs/apihub-swagger-config"})
}

func getConfigFolder() string {
	folder := os.Getenv("AGENT_CONFIG_FOLDER")
	if folder == "" {
		log.Warn("AGENT_CONFIG_FOLDER is not set, using default value: '.'")
		folder = "."
	}
	return folder
}

type systemInfoServiceImpl struct {
	config config.Config
}

func (g systemInfoServiceImpl) GetBackendVersion() string {
	return g.config.TechnicalParameters.Version
}

func (g systemInfoServiceImpl) GetApihubUrl() string {
	return g.config.TechnicalParameters.Apihub.URL
}

func (g systemInfoServiceImpl) GetAgentUrl() string {
	return g.config.TechnicalParameters.AgentUrl
}

func (g systemInfoServiceImpl) GetAccessToken() string {
	return g.config.TechnicalParameters.Apihub.AccessToken
}

func (g systemInfoServiceImpl) GetCloudName() string {
	return g.config.TechnicalParameters.CloudName
}

func (g systemInfoServiceImpl) GetAgentNamespace() string {
	return g.config.TechnicalParameters.Namespace
}

func (g systemInfoServiceImpl) GetExcludeLabels() []string {
	return g.config.Discovery.ExcludeLabels
}

func (g systemInfoServiceImpl) GetGroupingLabels() []string {
	return g.config.Discovery.GroupingLabels
}

func (g systemInfoServiceImpl) GetAgentName() string {
	return g.config.TechnicalParameters.AgentName
}

func (g systemInfoServiceImpl) GetDiscoveryTimeout() time.Duration {
	return time.Duration(g.config.Discovery.TimeoutSec) * time.Second
}

func (g systemInfoServiceImpl) InsecureProxyEnabled() bool {
	return g.config.Security.InsecureProxy
}

func (g systemInfoServiceImpl) GetBasePath() string {
	return g.config.TechnicalParameters.BasePath
}

func (g systemInfoServiceImpl) GetPaasPlatform() string {
	return g.config.TechnicalParameters.PaasPlatform
}

func (g systemInfoServiceImpl) GetDiscoveryUrls() config.ApiTypeUrlsConfig {
	return g.config.Discovery.Urls
}
