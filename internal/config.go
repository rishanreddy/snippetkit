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

func SetAPIKey(apiKey string) (bool) {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/snippetkit/config.yaml")
	if valid, err := VerifyToken(apiKey); !valid || err != nil {
		Error("API token is invalid or expired. Please run 'snippetkit login' to authenticate", err, nil)
		return false
	}
	viper.Set("api_key", apiKey)
	viper.WriteConfigAs(configPath)
	return true
}

func RemoveAPIKey() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/snippetkit/config.yaml")
	viper.Set("api_key", "")
	return viper.WriteConfigAs(configPath)
}