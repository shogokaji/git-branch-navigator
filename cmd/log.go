package cmd

import (
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Displays up to 10 branches from Git switch history and allows selection to switch.",
	Long: `This command displays up to 10 branches from the 'git reflog' command (excluding deleted branches).
By selecting a branch from the displayed list, you can switch to the chosen branch.
	`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
