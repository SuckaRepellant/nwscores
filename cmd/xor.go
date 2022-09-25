package cmd

import (
	"fmt"
	"log"
	"nwscores/lib"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// xorCmd represents the xor command
var xorCmd = &cobra.Command{
	Use:   "xor",
	Short: "xor a save file",
	Run: func(cmd *cobra.Command, args []string) {
		xoredSaveContent, err := os.ReadFile(viper.GetString("savefile"))
		if err != nil {
			log.Fatalln("Error reading save file", err)
		}

		saveContent := make([]byte, len(xoredSaveContent))

		xorKey := lib.GetXorKey()
		for i := 0; i < len(xoredSaveContent); i++ {
			saveContent[i] = xoredSaveContent[i] ^ xorKey[i%len(xorKey)]
		}

		fmt.Println(string(saveContent))
	},
}

func init() {
	rootCmd.AddCommand(xorCmd)
}
