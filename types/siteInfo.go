package types

type SiteInfo struct { // Origin as key
	SiteKey        string   `yaml:"site_key"`
	SecretKey      string   `yaml:"secret_key"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}
