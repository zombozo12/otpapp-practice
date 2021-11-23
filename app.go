package main

import (
	"otpapp-native/config"
	"otpapp-native/redis"
)

func InitApp(c *config.GlobalCfg) error {
	config.InitDB(&c.Database)

	rd, err := redis.Init(c.Redis)
	if err != nil {
		return err
	}

	defer rd.Close()

	return err
}
