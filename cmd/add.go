package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"snippetkit/internal"

	"github.com/leaanthony/spinner"

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

        // Load API token from config file
        apiToken, err := internal.GetAPIKey()
        if err != nil {
            internal.Error("Failed to get API key", err, nil)
            return
        }
        myspinner := spinner.New()
        myspinner.Start("Fetching snippet...")
        fmt.Println()
        // Fetch snippet from API
        snippet, err := internal.FetchSnippet(snippetID, apiToken)

        if err != nil {
            internal.Error("Failed to fetch snippet", err, nil)
            myspinner.Error("Failed to fetch snippet")
            return
        }
        myspinner.Success("Snippet fetched successfully")

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

            internal.Info(fmt.Sprintf("📂 Default install path: %s", defaultPath), nil)

            // Prompt user if they want to change the default install path
            prompt := promptui.Select{
                Label: "Do you want to change the install path?",
                Items: []string{"No", "Yes", "Cancel"},
            }

            _, choice, err := prompt.Run()
            if err != nil {
                if err == promptui.ErrInterrupt {
                    internal.Info("Prompt cancelled", nil)
                } else {
                    internal.Error("Error with prompt", err, nil)
                }
                return
            }

            if choice == "Cancel" {
                return
            }

            if choice == "Yes" {
                pathPrompt := promptui.Prompt{
                    Label: "Enter new install path",
                    Validate: func(input string) error {
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
                        return nil
                    },
                }
                installPath, err = pathPrompt.Run()
                if err != nil {
                    internal.Error("Error entering path", err, nil)
                    return
                }
            } else {
                installPath = defaultPath
            }
        }

        // Ensure directory exists
        if err := internal.EnsureDirExists(installPath); err != nil {
            internal.Error("Error creating directory", err, nil)
            return
        }

        // Handle overwrite
        if !addForce && internal.FileExists(installPath) {
            internal.Warn("⚠️ File already exists. Use --force to overwrite.", nil)
            return
        }

        // Write snippet
        if err := internal.WriteToFile(installPath, snippet.Code); err != nil {
            internal.Error("Error writing snippet", err, nil)
            return
        }

        internal.Info(fmt.Sprintf("✅ Snippet installed successfully at %s", installPath), nil)
    },
}