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
	Short: "Search for snippets by name, language, or tags",
	Long:  "Search for snippets in the SnippetKit repository based on name, programming language, or associated tags.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
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

// searchWithSpinner runs the search command with a spinner
func searchWithSpinner(query string) {
	var wg sync.WaitGroup
	wg.Add(1)

	// Start spinner

	// Run search in a separate goroutine
	go func() {
		defer wg.Done()

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

		myspinner := spinner.New()
		myspinner.Start("Searching for snippets...")
		// Fetch search results
		snippets, err := internal.SearchSnippets(query, langFilter, tagFilter, limit, apiToken)
		if err != nil {
			internal.Error("Error searching snippets", err, nil)
			myspinner.Error("Failed to fetch search results.")
			return
		}

		myspinner.Success("Search completed successfully")

		// Display results
		if len(snippets) == 0 {
			fmt.Println(infoStyle.Render("\n No snippets found for query: ") + labelStyle.Render(query))
			return
		}

		// Format search results
		var output string
		output += "\n" + titleStyle.Render("> Search Results:") + "\n"
		for _, snippet := range snippets {
			output += fmt.Sprintf(
				"* %s  %s [%s]\n",
				labelStyle.Render(snippet.ShortID),
				snippet.Title,
				snippet.Language,
			)
		}

		output += "\n" + infoStyle.Render("To install a snippet, run:") + "\n"
		output += "   " + infoStyle.Render("snippetkit add <snippet_id>") + "\n"

		fmt.Println(output)
	}()

	// Wait for the search to complete
	wg.Wait()
}
