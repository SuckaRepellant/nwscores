package cmd

import (
	"embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"nwscores/lib"
	"text/template"
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

//go:embed templates
var template_files embed.FS

func ServeCompare(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		tmpl, err := template.ParseFS(template_files, "templates/compare.html.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, "missing template", 500)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to execute template", 500)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		return
	case "POST":
		if current_wrs == nil {
			http.Error(w, "records not yet initialized", 500)
			return
		}

		file, _, err := r.FormFile("save")
		if err != nil {
			log.Println(err)
			http.Error(w, "error getting your save as a file", 500)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "error reading your save", 500)
			return
		}

		pbs, err := lib.RetrievePBsFromSave(fileBytes)
		if err != nil {
			http.Error(w, "error unmarshaling your save data", 500)
			return
		}

		var data [][]string
		for i, k := range lib.GetLevels() {
			wrTime := current_wrs.WRs[i]
			pbTime := pbs[i]
			diff := pbTime - wrTime
			data = append(data, []string{k, pbTime.String(), wrTime.String(), diff.String()})
		}

		tmpl, err := template.ParseFS(template_files, "templates/compare.html.tmpl")
		if err != nil {
			log.Println(err)
			http.Error(w, "missing template", 500)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to execute template", 500)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		return
	}
}

var wrsCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve as a web api",
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
		http.HandleFunc("/", ServeCompare)
		log.Println("Listening")
		log.Fatal(http.ListenAndServe(":8000", nil))
	},
}

func init() {
	rootCmd.AddCommand(wrsCmd)
}
