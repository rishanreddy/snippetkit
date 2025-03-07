package cmd

import (
	"fmt"
	"snippetkit/internal"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// searchModel handles the spinner state
type searchModel struct {
	spinner  spinner.Model
	query    string
	snippets []internal.Snippet
	err      error
	loading  bool
	done     bool
}

// Init starts the spinner animation
func (m searchModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles events and data fetching
func (m searchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if m.loading {
			return m, m.spinner.Tick
		}

	case fetchSnippetsMsg:
		m.snippets = msg.snippets
		m.err = msg.err
		m.loading = false
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI
func (m searchModel) View() string {
	if m.loading {
		return fmt.Sprintf("\n %s Searching for '%s'...\n", m.spinner.View(), m.query)
	}

	if m.err != nil {
		return fmt.Sprintf("❌ Error searching snippets: %v\n", m.err)
	}

	if len(m.snippets) == 0 {
		return fmt.Sprintf("🔍 No snippets found for: %s\n", m.query)
	}

	output := "\n🔍 Search Results:\n"
	for _, snippet := range m.snippets {
		output += fmt.Sprintf("📌 %s - %s [%s]\n", snippet.ShortID, snippet.Title, snippet.Language)
	}

	output += "\n💡 To install a snippet, use:\n"
	output += "   snippetkit add <snippet_id>\n"

	return output
}

// fetchSnippetsMsg is used to pass search results into the UI
type fetchSnippetsMsg struct {
	snippets []internal.Snippet
	err      error
}

// searchWithSpinner runs the search command with a spinner
func searchWithSpinner(query string) {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")) // Fix style reference

	m := searchModel{
		spinner: s,
		query:   query,
		loading: true,
	}

	p := tea.NewProgram(m)

	internal.LoadConfig()

	// Run search in a separate goroutine
	go func() {
		token, tokenErr := internal.GetAPIToken()
		if tokenErr != nil {
			p.Send(fetchSnippetsMsg{snippets: nil, err: tokenErr})
			return
		}
		snippets, err := internal.SearchSnippets(query, langFilter, tagFilter, limit, token)
		p.Send(fetchSnippetsMsg{snippets: snippets, err: err})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("❌ Error running search:", err)
	}
}
