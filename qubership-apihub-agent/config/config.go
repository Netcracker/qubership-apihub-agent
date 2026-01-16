package config

type Config struct {
	TechnicalParameters TechnicalParameters
	Security            SecurityConfig
	Discovery           DiscoveryConfig
}

type TechnicalParameters struct {
	BasePath      string
	ListenAddress string `validate:"required"`
	Version       string
	Apihub        ApihubConfig
	AgentUrl      string `validate:"required"`
	CloudName     string `validate:"required,slug_only_characters"`
	Namespace     string `validate:"required,slug_only_characters"`
	AgentName     string `validate:"required,slug_only_characters"`
	PaasPlatform  string `validate:"required"`
}

type ApihubConfig struct {
	URL         string `validate:"required"`
	AccessToken string `validate:"required" sensitive:"true"`
}

type SecurityConfig struct {
	AllowedOrigins []string
	InsecureProxy  bool
}

type DiscoveryConfig struct {
	ExcludeLabels  []string
	GroupingLabels []string
	TimeoutSec     int
}
