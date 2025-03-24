package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snippet from your terminal (coming soon)",
	Long:  "This feature is under development and will be available in a future release.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(titleStyle.Render("\nSnippet Creation Coming Soon"))
		fmt.Println(divider)
		fmt.Println(infoStyle.Render("We're working on a powerful feature that will allow you to:"))
		fmt.Println()
		fmt.Println("  • Create new snippets directly from the CLI")
		fmt.Println("  • Upload them to your SnippetKit account")
		fmt.Println()
		fmt.Println(infoStyle.Render("In the meantime, manage your snippets at:"))
		fmt.Println("  " + urlStyle.Render("https://snippetkit.vercel.app/snippet"))
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
