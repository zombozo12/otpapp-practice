package config

type DatabaseCfg struct {
	Host string
	Port string
	Name string
	User string
	Pass string
}

type RedisCfg struct {
	URL      string `json:"url"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type GlobalCfg struct {
	Database DatabaseCfg
	Redis    RedisCfg
}
