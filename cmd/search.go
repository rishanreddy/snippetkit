package cmd

import (
	"fmt"
	"snippetkit/internal"
	"sync"

	"github.com/leaanthony/spinner"
	"github.com/spf13/cobra"
)

// Flags
var (
    langFilter string
    tagFilter  string
    limit      int
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
    Use:   "search [query]",
    Short: "Search for snippets by name or tags",
    Long:  `Search for snippets in the SnippetKit repository by name or tags.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        query := args[0]

        // Run the search with a spinner
        searchWithSpinner(query)
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)

    // Define search flags
    searchCmd.Flags().StringVarP(&langFilter, "lang", "l", "", "Filter snippets by programming language (e.g., typescript)")
    searchCmd.Flags().StringVarP(&tagFilter, "tag", "t", "", "Filter snippets by tag (e.g., ui, shadcnui)")
    searchCmd.Flags().IntVarP(&limit, "limit", "n", 10, "Limit the number of search results")
}

// fetchSnippetsMsg is used to pass search results into the UI
type fetchSnippetsMsg struct {
    snippets []internal.Snippet
    err      error
}

// searchWithSpinner runs the search command with a spinner
func searchWithSpinner(query string) {
    var wg sync.WaitGroup
    wg.Add(1)

    myspinner := spinner.New()
    myspinner.Start("Searching for snippets...")

    // Run search in a separate goroutine
    go func() {
        defer wg.Done()

        // Load API token from config file
        apiToken, err := internal.GetAPIKey()
        if err != nil {
            internal.Error("❌ Error getting API token:", err, nil)
            myspinner.Error("❌ Error getting API token")
            return
        }
        snippets, err := internal.SearchSnippets(query, langFilter, tagFilter, limit, apiToken)
        if err != nil {
            internal.Error("❌ Error searching snippets:", err, nil)
            myspinner.Error("❌ Error searching snippets")
            return
        }

        myspinner.Success("Search completed successfully")

        // Print search results
        if len(snippets) == 0 {
            internal.Info(fmt.Sprintf("🔍 No snippets found for: %s", query), nil)
            return
        }

        output := "\n🔍 Search Results:\n"
        for _, snippet := range snippets {
            output += fmt.Sprintf("📌 %s - %s [%s]\n", snippet.ShortID, snippet.Title, snippet.Language)
        }

        output += "\n💡 To install a snippet, use:\n"
        output += "   snippetkit add <snippet_id>\n"

        internal.Info(output, nil)
    }()

    // Wait for the search to complete
    wg.Wait()
}