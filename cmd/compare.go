package cmd

import (
	"log"
	"os"

	"nwscores/lib"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare your times to the IL WRs",
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Retrieving PB times from savefile...")
		pbs, err := lib.RetrievePBs(viper.GetString("savefile"))
		if err != nil {
			log.Fatalln("Error reading PBs", err)
		}

		log.Println("Retrieving WR times from Google Sheets...")
		wrs, err := lib.RetrieveWRsFromGoogle(
			viper.GetString("sheets.apikey"),
			viper.GetString("sheets.sheet_id"),
			viper.GetString("sheets.ranges"),
		)
		if err != nil {
			log.Fatalln("Error reading WRs", err)
		}

		var data [][]string
		for i, k := range lib.GetLevels() {
			wrTime := wrs[i]
			pbTime := pbs[i]
			diff := pbTime - wrTime
			data = append(data, []string{k, pbTime.String(), wrTime.String(), diff.String()})
		}

		switch output {
		case "pretty":
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "PB", "WR", "Diff"})
			table.AppendBulk(data)
			table.Render()
		default:
			log.Fatalln("Unknown output type")
		}
	},
}

var output string

func init() {
	rootCmd.AddCommand(compareCmd)
	compareCmd.Flags().StringVarP(&output, "output", "o", "pretty", "Output format (pretty,plain)")
}
