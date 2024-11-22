package deribit

import (
	"context"
	"os"
	"strconv"
)

const (
	RealBaseURL     = "wss://www.deribit.com/ws/api/v2/"
	TestBaseURL     = "wss://test.deribit.com/ws/api/v2/"
	RealRestBaseURL = "https://www.deribit.com/ws/api/v2/"
	TestRestBaseURL = "https://test.deribit.com/ws/api/v2/"
)

const (
	MaxTryTimes = 10000
)

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetConfig() *Configuration {
	autoReconnect, _ := strconv.ParseBool(getEnvWithDefault("DERIBIT_AUTO_RECONNECT", "true"))
	debugMode, _ := strconv.ParseBool(getEnvWithDefault("DERIBIT_DEBUG_MODE", "true"))
	realMode, _ := strconv.ParseBool(getEnvWithDefault("DERIBIT_REAL_MODE", "false"))

	wsBaseURL := TestBaseURL
	restBaseURL := TestRestBaseURL
	if realMode {
		wsBaseURL = RealBaseURL
		restBaseURL = RealRestBaseURL
	}

	return &Configuration{
		WsAddr:        wsBaseURL,
		RestAddr:      restBaseURL,
		ApiKey:        getEnvWithDefault("DERIBIT_API_KEY", ""),
		SecretKey:     getEnvWithDefault("DERIBIT_API_SECRET", ""),
		AutoReconnect: autoReconnect,
		DebugMode:     debugMode,
	}
}

type Configuration struct {
	Ctx           context.Context
	WsAddr        string `json:"ws_addr"`
	RestAddr      string `json:"rest_addr"`
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
	WSBaseURL     string `json:"ws_base_url"`
	RestBaseURL   string `json:"rest_base_url"`
}
