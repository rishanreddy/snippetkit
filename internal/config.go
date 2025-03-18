package internal

import (
	"fmt"

	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".config/snippetkit"))
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("logging_enabled", true)
	if viper.GetBool("logging_enabled") {
		fmt.Println("\033[1;32m✔ Logging is enabled\033[0m")
	} else {
		fmt.Println("\033[1;31m✘ Logging is disabled\033[0m")
	}

	if err := viper.ReadInConfig(); err != nil {
		Warn("No config file found. Using default settings.", nil)
	}
}

func GetAPIKey() (string, error) {
	apiToken := viper.GetString("api_key")
	if valid, err := VerifyToken(apiToken); !valid || err != nil {
		return "", fmt.Errorf("API token is invalid or expired. Please run 'snippetkit login' to authenticate")
	}
	if apiToken == "" {
		return "", fmt.Errorf("API token is missing. Please run 'snippetkit login' to authenticate")
	}
	return viper.GetString("api_key"), nil
}

func SetAPIKey(apiKey string) (bool, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/snippetkit/config.yaml")
	if valid, err := VerifyToken(apiKey); !valid || err != nil {
		Error("API token is invalid or expired. Please run 'snippetkit login' to authenticate", err, nil)
		return false, fmt.Errorf("API token is invalid or expired")
	}
	viper.Set("api_key", apiKey)
	viper.WriteConfigAs(configPath)
	return true, nil
}

func RemoveAPIKey() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/snippetkit/config.yaml")
	viper.Set("api_key", "")
	return viper.WriteConfigAs(configPath)
}

func GetConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config/snippetkit/config.yaml")
}
