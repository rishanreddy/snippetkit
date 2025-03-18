package cmd

import (
	"fmt"
	"snippetkit/internal"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and remove your API token",
	Long:  "Use this command to log out and remove your API token from the configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt user for confirmation
		prompt := promptui.Prompt{
			Label:     "Are you sure you want to log out? (y/N)",
			IsConfirm: true,
		}

		result, err := prompt.Run()
		if err == promptui.ErrInterrupt {
			fmt.Println(errorStyle.Render("\n Logout cancelled by user."))
			internal.Info("Logout cancelled.", nil)
			return
		}

		if result != "y" && result != "Y" {
			fmt.Println(errorStyle.Render("\n Error reading logout confirmation input (y/N)"))
			internal.Error("Error reading logout confirmation input", err, nil)
			return
		}
		if err != nil {
			fmt.Println(errorStyle.Render("\n Error reading input"))
			internal.Error("Error reading logout confirmation input", err, nil)
			return
		}

		// Remove the API token
		if err := internal.RemoveAPIKey(); err != nil {
			internal.Error("Error removing API token", err, nil)
			fmt.Println(errorStyle.Render("\n Error removing API token"))
			return
		}

		internal.Info("Successfully logged out and removed API token", nil)
		fmt.Println(successStyle.Render("\n Successfully logged out and removed API token"))
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
