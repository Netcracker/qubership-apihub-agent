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
	"reflect"
	"time"
	"unicode"

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
	GetDiscoveryUrls() config.UrlsConfig
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
	//TODO: what are defaults for discovery URLs ?
	viper.SetDefault("discovery.urls.openapi.config-urls", []string{"/v3/api-docs/swagger-config", "/swagger-resources"})
	viper.SetDefault("discovery.urls.openapi.doc-urls", []string{"/q/openapi?format=json", "/v3/api-docs?format=json"})
	viper.SetDefault("discovery.urls.graphql.config-urls", []string{"/api/graphql-server/schema/domains"})
	viper.SetDefault("discovery.urls.graphql.doc-urls", []string{"/api/graphql-server/schema", "/graphql", "/graphql/introspection"})
	viper.SetDefault("discovery.urls.apihub-config.config-urls", []string{"/v3/api-docs/apihub-swagger-config"})
	viper.SetDefault("discovery.urls.smartplug.config-urls", []string{"/smartplug/v1/api/config"})
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

func (g systemInfoServiceImpl) GetDiscoveryUrls() config.UrlsConfig {
	return g.config.Discovery.Urls
}

func PrintConfig(config interface{}) {
	log.Info("Loaded configuration:")
	printStruct("", reflect.ValueOf(config))
}

func printStruct(prefix string, v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			continue
		}
		value := v.Field(i)

		runes := []rune(field.Name)
		if len(runes) > 0 {
			runes[0] = unicode.ToLower(runes[0])
		}
		fieldName := string(runes)

		key := fieldName
		if prefix != "" {
			key = prefix + "." + fieldName
		}

		_, isSensitive := field.Tag.Lookup("sensitive")

		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				log.Infof("%s=<nil>", key)
				continue
			}
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Struct:
			printStruct(key, value)
		case reflect.Slice:
			if value.Type().Elem().Kind() == reflect.Struct && value.Len() > 0 {
				for j := 0; j < value.Len(); j++ {
					printStruct(fmt.Sprintf("%s[%d]", key, j), value.Index(j))
				}
			} else {
				printValue(key, value, isSensitive)
			}
		default:
			printValue(key, value, isSensitive)
		}
	}
}

func isValueEmpty(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

func printValue(key string, value reflect.Value, isSensitive bool) {
	var valStr string
	if isSensitive && !isValueEmpty(value) {
		valStr = "*****"
	} else if value.IsValid() && value.CanInterface() {
		valStr = fmt.Sprintf("%v", value.Interface())
	}
	log.Infof("%s=%s", key, valStr)
}
