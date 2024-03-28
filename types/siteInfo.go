package types

type SiteInfo struct { // Origin as key
	SiteKey        string   `yaml:"site_key"`
	SiteSecret     string   `yaml:"site_secret"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}
