package cmd

import (
	"fmt"
	"snippetkit/internal"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate and store your API token",
	Long:  "Use this command to set and save your API token in config.yaml for authentication.",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt user for API token
		prompt := promptui.Prompt{
			Label: "Enter your API token",
			Mask:  '*', // Hides input for security
		}

		apiToken, err := prompt.Run()
		if err != nil {
			fmt.Println("❌ Error reading input:", err)
			return
		}

		// Save the token using internal function
		err = internal.SaveAPIToken(apiToken)
		if err != nil {
			fmt.Println("❌ Error saving API token:", err)
			return
		}

		fmt.Println("✅ API token saved successfully!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
