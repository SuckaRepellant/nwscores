package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"nwscores/lib"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var current_wrs *lib.RemoteWRs

func ServeWRs(w http.ResponseWriter, r *http.Request) {
	if current_wrs == nil {
		http.Error(w, "records not yet initialized", 500)
		return
	}

	output, err := json.Marshal(current_wrs)
	if err != nil {
		http.Error(w, "error marshaling current records", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	w.Write([]byte("\n"))
}

var wrsCmd = &cobra.Command{
	Use:   "wrs",
	Short: "Liberate WR times from Google Sheets",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Launching sheet scraper")
		go func() {
			for {
				var err error
				current_wrs, err = lib.RetrieveWRsFromGoogle(
					viper.GetString("wrs.sheets.apikey"),
					viper.GetString("wrs.sheets.sheet_id"),
					viper.GetString("wrs.sheets.ranges"),
				)
				if err != nil {
					log.Println("Error reading WRs", err)
				}
				time.Sleep(1 * time.Hour)
			}
		}()

		http.HandleFunc("/wrs.json", ServeWRs)
		log.Println("Listening")
		log.Fatal(http.ListenAndServe(":8000", nil))
	},
}

func init() {
	rootCmd.AddCommand(wrsCmd)
}
