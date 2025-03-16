package cmd

import (
	"fmt"
	"os"
	"snippetkit/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "snippetkit",
	Short: "SnippetKit - Easily manage reusable code snippets",
	Long: `SnippetKit CLI allows you to search, add, and manage code snippets.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		internal.LoadConfig() // Load config before executing commands
		internal.InitLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := "internal.GetVersion() "// Pull the version dynamically
		internal.Info(fmt.Sprintf("SnippetKit CLI v%s", version), nil)
		cmd.Help() // Display the help command
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global Persistent Flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "Specify config file (default is $HOME/.snippetkit/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}
