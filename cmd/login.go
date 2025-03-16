package cmd

import (
	"snippetkit/internal"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var apiKey string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Authenticate and store your API token",
    Long:  "Use this command to set and save your API token in config.yaml for authentication.",
    Run: func(cmd *cobra.Command, args []string) {
        var apiToken string
        var err error

        if apiKey != "" {
            apiToken = apiKey
        } else {
            // Prompt user for API token
            prompt := promptui.Prompt{
                Label: "Enter your API token",
                Mask:  '*', // Hides input for security
            }

            apiToken, err = prompt.Run()
            if err == promptui.ErrInterrupt {
                internal.Info("Login cancelled.", nil)
                return
            }
            if err != nil {
                internal.Error("❌ Error reading input", err, nil)
                return
            }
        }

        // Save the token using internal function
        success := internal.SetAPIKey(apiToken)
   
        if err != nil {
            internal.Error("❌ Error saving API token", err, nil)
            return
        }

        if success {
            internal.Info("✅ API token saved successfully!", nil)
        }
   
    },
}

func init() {
    rootCmd.AddCommand(loginCmd)
    loginCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API token to authenticate")
}