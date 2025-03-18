package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"snippetkit/internal"
	"strings"

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
		// Show spinner while checking API status
		checkSpinner := spinner.New()
		checkSpinner.Start("Checking auth status...")
		// Load API token from config file
		apiToken, err := internal.GetAPIKey()
		if err != nil {
			checkSpinner.Error("Failed to authenticate")
			fmt.Print("Error authenticating. Run 'snippetkit login' to authenticate.")
			internal.Error("Failed to get API key", err, nil)
			return
		}

		checkSpinner.Success("Authenticated successfully")

		// Show fancy spinner while fetching snippet
		myspinner := spinner.New()
		myspinner.Start(fmt.Sprintf("Fetching snippet %s...", snippetID))

		// Fetch snippet from API
		snippet, err := internal.FetchSnippet(snippetID, apiToken)
		if err != nil {
			myspinner.Error(fmt.Sprintf("Failed to fetch snippet with ID: %s", snippetID))
			internal.Error("Failed to fetch snippet", err, nil)
			return
		}
		myspinner.Success(fmt.Sprintf("Snippet %s fetched successfully", snippetID))

		// Show snippet info
		if !addSilent {
			fmt.Println(titleStyle.Render("> Snippet Details:"))
			fmt.Println(divider)
			fmt.Println(labelStyle.Render("Title: ") + snippet.Title)
			fmt.Println(labelStyle.Render("Language: ") + snippet.Language)
			fmt.Println(labelStyle.Render("Tags: ") + strings.Join(snippet.Tags, ", "))
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

			if !addSilent {
				fmt.Printf("%s %s\n", titleStyle.Render("Default install path:"), defaultPath)
			}

			// Skip prompt in silent mode
			if addSilent {
				installPath = defaultPath
			} else {
				// Prompt user if they want to change the default install path
				prompt := promptui.Select{
					Label: "Do you want to change the install path?",
					Items: []string{"No", "Yes"},
				}

				_, choice, err := prompt.Run()

				if err == promptui.ErrInterrupt {
					fmt.Println(errorStyle.Render("\n Operation cancelled by user."))
					internal.Error("Prompt cancelled by user", nil, nil)
					os.Exit(1)
				}

				if err != nil {
					fmt.Println(errorStyle.Render("\nError reading input."))
					return
				}

				if choice == "Yes" {
					pathPrompt := promptui.Prompt{
						Label: "Enter new install path",
						Validate: func(input string) error {
							if len(input) == 0 {
								return fmt.Errorf("path cannot be empty")
							}
							return nil
						},
					}
					installPath, err = pathPrompt.Run()
					if err == promptui.ErrInterrupt {
						fmt.Println(errorStyle.Render("\n Operation cancelled by user."))
						os.Exit(1)
					}
					if err != nil {
						fmt.Println(errorStyle.Render("Error: Could not read input."))
						return
					}
				} else {
					installPath = defaultPath
				}
			}
		}

		// Ensure parent directory exists
		parentDir := filepath.Dir(installPath)
		if parentDir != "." && !internal.FileExists(parentDir) {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				fmt.Println(errorStyle.Render(fmt.Sprintf("Failed to create directory: %s", parentDir)))
				internal.Error(fmt.Sprintf("Error creating directory: %s", parentDir), err, nil)
				return
			}
		}

		// Handle overwrite
		if !addForce && internal.FileExists(installPath) {
			fmt.Println(warningStyle.Render(fmt.Sprintf("\n Skipping installation. File %s already exists. Use --force to overwrite.", installPath)))
			internal.Warn("File already exists. Use --force to overwrite.", nil)
			return
		}

		// Write snippet to file
		if err := internal.WriteToFile(installPath, snippet.Code); err != nil {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Failed to write snippet to file: %s", installPath)))
			internal.Error("Error writing snippet", err, nil)
			return
		}

		// Show success message
		if !addSilent {
			fmt.Println(successStyle.Render("\n Snippet installed successfully!"))
		}
		internal.Info(fmt.Sprintf("Snippet installed successfully at %s", installPath), nil)
	},
}
