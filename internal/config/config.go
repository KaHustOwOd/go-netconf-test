package config

import (
	"log/slog"
	"os"
)

type AppConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	ConfigPath string
	Namespace  string
}

func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		Host:       os.Getenv("NETCONF_HOST"),
		Port:       os.Getenv("NETCONF_PORT"),
		Username:   os.Getenv("NETCONF_USERNAME"),
		Password:   os.Getenv("NETCONF_PASSWORD"),
		ConfigPath: os.Getenv("CONFIG_PATH"),
		Namespace:  os.Getenv("YANG_NAMESPACE"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" || cfg.ConfigPath == "" {
		slog.Error("Error: Missing of mandatory env variables (NETCONF_HOST, NETCONF_PORT, NETCONFSSH_USERNAME, CONFIG_PATH)")
		os.Exit(1)
	}
	return cfg
}
