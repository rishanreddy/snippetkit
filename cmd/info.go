package cmd

import (
	"encoding/json"
	"fmt"
	"snippetkit/internal"

	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"

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

		// Show loading spinner
		// s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		// s.Start()

		// Fetch snippet details from API
		snippet, err := internal.FetchSnippet(snippetID)
		// s.Stop()
		if err != nil {
			fmt.Println("❌ Error fetching snippet info:", err)
			return
		}

		// Output JSON if requested
		if jsonOutput {
			jsonData, _ := json.MarshalIndent(snippet, "", "  ")
			fmt.Println(string(jsonData))
			return
		}

		// Print snippet metadata
		fmt.Println("📜 Snippet Details:")
		fmt.Println("----------------------------")
		fmt.Println("🆔 ID:       ", snippetID)
		fmt.Println("📌 Title:    ", snippet.Title)
		fmt.Println("📝 Language: ", snippet.Language)
		fmt.Println("📂 Path:     ", snippet.Path)
		fmt.Println("----------------------------")
		fmt.Println("📄 Code Preview:")
		fmt.Println("----------------------------")

		// Trim or show full code
		code := snippet.Code
		if !fullOutput {
			code = formatCodePreview(code)
		}

		// Print syntax-highlighted code
		printHighlightedCode(code, snippet.Language)

		fmt.Println("----------------------------")
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
		fmt.Println(code) // Fallback to plain text if highlighting fails
	}
}
