package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"snippetkit/internal"

	"github.com/spf13/cobra"
)

// Flags
var addPath string
var addForce bool
var addSilent bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [snippet ID]",
	Short: "Fetch and add a snippet to your project",
	Long:  `Fetch a snippet from SnippetKit API and install it into your project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snippetID := args[0]

		// Fetch snippet from API
		snippet, err := internal.FetchSnippet(snippetID)
		if err != nil {
			fmt.Println("❌ Error fetching snippet:", err)
			return
		}

		// Determine install path
		var installPath string

		if addPath != "" {
			installPath = addPath // Use provided path
		} else {
			cwd, _ := os.Getwd()
			defaultPath := filepath.Join(cwd, snippet.Path)
			if snippet.Path == "" {
				defaultPath = filepath.Join(cwd, fmt.Sprintf("%s.%s", snippet.Title, snippet.Language))
			}

			// Ask user if they want to change the install path
			if internal.YesNoPrompt(fmt.Sprintf("📂 Default install path: %s\nWould you like to change the path?", defaultPath), false) {
				installPath = internal.GetUserInput("Enter new install path:", defaultPath)
			} else {
				installPath = defaultPath
			}
		}

		// Ensure directory exists
		if err := internal.EnsureDirExists(installPath); err != nil {
			fmt.Println("❌ Error creating directory:", err)
			return
		}

		// Handle overwrite
		if !addForce && internal.FileExists(installPath) {
			fmt.Println("⚠️ File already exists. Use --force to overwrite.")
			return
		}

		// Write snippet
		if err := internal.WriteToFile(installPath, snippet.Code); err != nil {
			fmt.Println("❌ Error writing snippet:", err)
			return
		}

		fmt.Println("✅ Snippet installed successfully at", installPath)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addPath, "path", "p", "", "Specify install path for the snippet (skips prompt)")
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Force overwrite if file already exists")
	addCmd.Flags().BoolVarP(&addSilent, "silent", "s", false, "Suppress output")
}
