package cmd

import (
	"encoding/json"
	"fmt"
	"snippetkit/internal"

	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/leaanthony/spinner"

	"github.com/spf13/cobra"
)

var jsonOutput bool
var fullOutput bool

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info [snippet ID]",
	Short: "Preview a snippet before adding it",
	Long:  `The 'info' command retrieves metadata and a preview of the snippet code from SnippetKit's API.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snippetID := args[0]

		checkSpinner := spinner.New()
		checkSpinner.Start("Checking auth status...")
		apiToken, err := internal.GetAPIKey()
		if err != nil {
			checkSpinner.Error("Failed to authenticate")
			fmt.Print("Error authenticating. Run 'snippetkit login' to authenticate.")
			internal.Error("Failed to get API key", err, nil)
			return
		}

		checkSpinner.Success("Authenticated successfully")

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

		// Output JSON if requested
		if jsonOutput {
			jsonData, _ := json.MarshalIndent(snippet, "", "  ")
			fmt.Println(string(jsonData))
			return
		}

		// Print snippet metadata
		fmt.Println(titleStyle.Render("> Snippet Details:"))
		fmt.Println(divider)
		fmt.Println(labelStyle.Render("ID: ") + snippetID)
		fmt.Println(labelStyle.Render("Title: ") + snippet.Title)
		fmt.Println(labelStyle.Render("Language: ") + snippet.Language)
		fmt.Println(labelStyle.Render("Tags: ") + strings.Join(snippet.Tags, ", "))
		fmt.Println(labelStyle.Render("Path: ") + snippet.Path)
		fmt.Println(divider)
		fmt.Println(titleStyle.Render("Code Preview"))
		fmt.Println(divider)

		// Trim or show full code
		code := snippet.Code
		if !fullOutput {
			code = formatCodePreview(code)
		}

		// Print syntax-highlighted code
		printHighlightedCode(code, snippet.Language)

		fmt.Println(divider)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output snippet info as JSON")
	infoCmd.Flags().BoolVarP(&fullOutput, "full", "f", false, "Show full snippet instead of a preview")
}

// formatCodePreview trims and formats the snippet code for preview
func formatCodePreview(code string) string {
	lines := strings.Split(code, "\n")
	maxLines := 10 // Limit preview to 10 lines

	if len(lines) > maxLines {
		return strings.Join(lines[:maxLines], "\n") + "\n... (preview truncated)"
	}

	return strings.Join(lines, "\n")
}

// printHighlightedCode prints the snippet with syntax highlighting
func printHighlightedCode(code, language string) {
	// Determine lexer (detect from language or analyze content)
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Analyse(code) // Auto-detect language
	}
	if lexer == nil {
		lexer = lexers.Fallback // Default to plaintext if no match
	}

	// Highlight and print the code
	err := quick.Highlight(os.Stdout, code, lexer.Config().Name, "terminal16m", "onedark")
	if err != nil {
		internal.Warn("Error highlighting code, falling back to plain text", nil)
	}
}
