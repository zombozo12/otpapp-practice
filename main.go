package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"otpapp-native/config"
	"otpapp-native/router"
)

func main() {
	var cfg *config.GlobalCfg

	dot := godotenv.Load(".env")
	if dot != nil {
		log.Fatal("Error loading .env file")
	}

	cfg = &config.GlobalCfg{}

	cfg.ApplyConfig()

	err := InitApp(cfg)
	if err != nil {
		log.Fatalf("Failed to init: error: %+v", err)
	}

	rt := router.Init()
	err = http.ListenAndServe("localhost:8080", rt)
	if err != nil {
		log.Fatalf("Failed to server. Error:  %+v", err)
		return
	}

	fmt.Println("Server available at localhost:8080")
}
