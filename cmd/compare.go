package cmd

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"nwscores/lib"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare your times to the IL WRs",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		srv, err := sheets.NewService(ctx, option.WithAPIKey(viper.GetString("sheets.apikey")))
		if err != nil {
			log.Fatalln("Error creating sheets service", err)
		}

		wrTimes := []string{}

		for _, rng := range strings.Split(viper.GetString("sheets.ranges"), ";") {
			log.Println("Retrieving WR times from Google Sheets...", rng)
			resp, err := srv.Spreadsheets.Values.Get(viper.GetString("sheets.sheet_id"), rng).Do()
			if err != nil {
				log.Fatalln("Error retrieving some WRs", err)
			}
			for _, row := range resp.Values {
				replacedTime := strings.Replace(row[0].(string), ":", "m", 1) + "s"
				wrTimes = append(wrTimes, replacedTime)
			}
		}

		log.Println("Retrieving PB times from savefile...")
		psd, err := lib.RetrievePSD(viper.GetString("savefile"))
		if err != nil {
			log.Fatalln("Error reading save file", err)
		}

		var data [][]string

		for i, k := range lib.GetLevels() {
			bestTime := (time.Duration(psd.LevelStats.Values[i].TimeBestMicroseconds) * time.Microsecond)
			bestTime = bestTime.Truncate(time.Microsecond * 1000)
			wrTime, err := time.ParseDuration(wrTimes[i])
			diff := bestTime - wrTime
			if err != nil {
				log.Fatalln(err)
			}
			data = append(data, []string{k, bestTime.String(), wrTime.String(), diff.String()})
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
