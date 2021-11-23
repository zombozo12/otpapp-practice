package config

import "os"

func (cfg *GlobalCfg) ApplyConfig() {
	// idk why it looks like this
	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Pass = os.Getenv("DB_PASS")

	cfg.Redis.URL = os.Getenv("REDIS_URL")
	cfg.Redis.Port = os.Getenv("REDIS_PORT")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
}
