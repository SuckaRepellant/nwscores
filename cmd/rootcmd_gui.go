//go:build gui

package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "nwscores",
	Short: "Extract your Neon White scores from your save file",
	Run: func(cmd *cobra.Command, args []string) {
		RunGui()
	},
}
