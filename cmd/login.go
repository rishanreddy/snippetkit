package cmd

import (
	"fmt"
	"snippetkit/internal"

	"github.com/leaanthony/spinner"
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
				fmt.Println(errorStyle.Render("\n Login cancelled by user."))
				internal.Info("Login cancelled.", nil)
				return
			}
			if err != nil {
				fmt.Println(errorStyle.Render("\n Error reading input"))
				internal.Error("Error reading API token input", err, nil)
				return
			}
		}

		myspinner := spinner.New()
		myspinner.Start("Saving API token...")
		// Save the token using internal function
		success, apiErr := internal.SetAPIKey(apiToken)

		if apiErr != nil {
			myspinner.Error("API token invalid or expired")
			internal.Error("Error saving API token", err, nil)
			return
		}

		// Confirmation message
		if success {
			myspinner.Success("API token saved successfully!")
			internal.Info("A new API token was saved successfully!", nil)

			// Show where the API key is stored
			configPath := internal.GetConfigPath()
			fmt.Println(infoStyle.Render(fmt.Sprintf("\n API key stored in: %s", configPath)))
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API token to authenticate")
}
