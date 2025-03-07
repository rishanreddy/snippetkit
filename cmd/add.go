package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"snippetkit/internal"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Flags
var addPath string
var addForce bool
var addSilent bool

func init() {
	rootCmd.AddCommand(addCmd)

	// Define CLI flags
	addCmd.Flags().StringVarP(&addPath, "path", "p", "", "Specify install path for the snippet")
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Force overwrite if file exists")
	addCmd.Flags().BoolVarP(&addSilent, "silent", "s", false, "Suppress output")
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [snippet ID]",
	Short: "Fetch and add a snippet to your project",
	Long:  `Fetch a snippet from SnippetKit API and install it into your project.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snippetID := args[0]

		internal.LoadConfig()
		// Load API token from config file
		apiToken, err := internal.GetAPIToken()
		if err != nil {
			fmt.Println("❌", err)
			return
		}

		// Fetch snippet from API
		snippet, err := internal.FetchSnippet(snippetID, apiToken)
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

			fmt.Println("📂 Default install path:", defaultPath)

			// Prompt user if they want to change the default install path
			prompt := promptui.Select{
				Label: "Do you want to change the install path?",
				Items: []string{"No", "Yes"},
			}

			_, choice, err := prompt.Run()
			if err != nil {
				fmt.Println("❌ Error with prompt:", err)
				return
			}

			if choice == "Yes" {
				pathPrompt := promptui.Prompt{
					Label:   "Enter new install path",

					Validate: func(input string) error {
						if len(input) == 0 {
							if len(input) == 0 {
								return fmt.Errorf("path cannot be empty")
							}
							cwd, err := os.Getwd()
							if err != nil {
								return fmt.Errorf("error getting current working directory: %w", err)
							}
							if !filepath.IsAbs(input) && !filepath.IsAbs(filepath.Join(cwd, input)) {
								return fmt.Errorf("invalid path format")
							}
						}
						return nil
					},
				}
				installPath, err = pathPrompt.Run()
				if err != nil {
					fmt.Println("❌ Error entering path:", err)
					return
				}
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
