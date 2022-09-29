package cmd

import (
	"fmt"
	"log"
	"nwscores/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// xorCmd represents the xor command
var xorCmd = &cobra.Command{
	Use:   "xor",
	Short: "xor a file with the save key",
	Run: func(cmd *cobra.Command, args []string) {
		saveContent, err := lib.GetPlainSave(viper.GetString("savefile"))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(saveContent))
	},
}

func init() {
	rootCmd.AddCommand(xorCmd)
}
