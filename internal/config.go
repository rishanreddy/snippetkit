package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig initializes Viper and loads the API token from config.yaml
func LoadConfig() {
	viper.SetConfigName("config") // Looks for "config.yaml"
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.snippetkit")  // call multiple times to add many search paths
	viper.AddConfigPath("/etc/snippetkit/")   // path to look for the config file in
	viper.AddConfigPath(".")               // optionally look for config in the working directory

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("❌ Error reading config file:", err)
		}
	}
}

// GetAPIToken retrieves the API token from Viper
func GetAPIToken() (string, error) {
	apiToken := viper.GetString("api_token")
	if valid, err := VerifyToken(apiToken); !valid || err != nil {
		return "", fmt.Errorf("API token is invalid or expired. Please run 'snippetkit login' to authenticate")
	}
	if apiToken == "" {
		return "", fmt.Errorf("API token is missing. Please run 'snippetkit login' to authenticate")
	}
	return apiToken, nil
}

// SaveAPIToken saves the API token in config.yaml
func SaveAPIToken(token string) error {
	viper.Set("api_token", token) // Store token in config
	if err := viper.WriteConfigAs("config.yaml"); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
