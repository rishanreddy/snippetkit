package cmd

import (
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
        if err != nil {
            if err == promptui.ErrAbort {
                internal.Info("Logout cancelled.", nil)
                return
            }
            internal.Error("❌ Error reading input", err, nil)
            return
        }

        if result != "y" && result != "Y" {
            internal.Info("Logout cancelled.", nil)
            return
        }

        // Remove the API token
        if err := internal.RemoveAPIKey(); err != nil {
            internal.Error("❌ Error removing API token", err, nil)
            return
        }
        
        internal.Info("✅ Successfully logged out and removed API token", nil)
        
    },
}

func init() {
    rootCmd.AddCommand(logoutCmd)
}