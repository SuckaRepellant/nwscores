package cmd

import (
	"nwscores/lib"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nwscores",
	Short: "Extract your Neon White scores from your save file",
	Run: func(cmd *cobra.Command, args []string) {
		a := app.New()
		w := a.NewWindow("nwscores")
		w.Resize(fyne.NewSize(600, 800))

		saveFile_edit := widget.NewEntry()
		saveFile_edit.SetText(viper.GetString("pbs.savefile"))

		pbsSettingsForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "Savefile", Widget: saveFile_edit},
			},
		}

		wrsHttpSkip_edit := widget.NewCheck("", func(bool) {})
		wrsHttpSkip_edit.SetChecked(viper.GetBool("wrs.http.skip"))

		wrsHttpUrl_edit := widget.NewEntry()
		wrsHttpUrl_edit.SetText(viper.GetString("wrs.http.url"))

		wrsHttpSettingsForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "Skip", Widget: wrsHttpSkip_edit},
				{Text: "URL", Widget: wrsHttpUrl_edit},
			},
		}

		apikey_edit := widget.NewEntry()
		apikey_edit.SetText(viper.GetString("wrs.sheets.apikey"))
		sheet_id_edit := widget.NewEntry()
		sheet_id_edit.SetText(viper.GetString("wrs.sheets.sheet_id"))
		ranges_edit := widget.NewEntry()
		ranges_edit.SetText(viper.GetString("wrs.sheets.ranges"))

		wrsSheetsSettingsForm := &widget.Form{
			Items: []*widget.FormItem{
				{Text: "API Key", Widget: apikey_edit},
				{Text: "Sheet ID", Widget: sheet_id_edit},
				{Text: "Ranges", Widget: ranges_edit},
			},
		}

		writeButton := widget.NewButton("Write", func() {
			viper.Set("pbs.savefile", saveFile_edit.Text)

			viper.Set("wrs.http.skip", wrsHttpSkip_edit.Checked)
			viper.Set("wrs.http.url", wrsHttpUrl_edit.Text)

			viper.Set("wrs.sheets.apikey", apikey_edit.Text)
			viper.Set("wrs.sheets.sheet_id", sheet_id_edit.Text)
			viper.Set("wrs.sheets.ranges", ranges_edit.Text)
			err := viper.WriteConfig()
			if err != nil {
				dialog.NewError(err, w).Show()
			}
		})
		opts := dialog.NewCustom("Options", "Close", container.NewVBox(
			widget.NewLabelWithStyle("PBs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			pbsSettingsForm,
			widget.NewSeparator(),

			widget.NewLabelWithStyle("WRs - HTTP", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			wrsHttpSettingsForm,
			widget.NewSeparator(),

			widget.NewLabelWithStyle("WRs - Google Sheets", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			wrsSheetsSettingsForm,
			widget.NewSeparator(),

			writeButton,
		), w)
		opts.Resize(fyne.NewSize(500, 200))

		//opts.MinSize().Add(fyne.NewSize(500, 100))

		var data = [][]string{
			[]string{"Name", "PB", "WR", "Diff"},
		}

		var state_PBs, state_WRs []time.Duration

		for _, lv := range lib.GetLevels() {
			//                          Name PB       WR       Diff
			data = append(data, []string{lv, "0.000", "0.000", "0ms"})
		}

		diffs := widget.NewTable(
			func() (int, int) {
				return len(data), len(data[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("5m39s.555")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(data[i.Row][i.Col])
			},
		)

		diffs.SetColumnWidth(0, 180)

		refreshDiffs := func() {
			if len(state_PBs) == 0 || len(state_WRs) == 0 {
				return
			}
			for i, _ := range lib.GetLevels() {
				data[i+1][3] = (state_PBs[i] - state_WRs[i]).String()
			}
		}

		pbRefresh := widget.NewButton("Refresh PBs", func() {
			var err error
			state_PBs, err = lib.RetrievePBs(viper.GetString("pbs.savefile"))
			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}
			for i, v := range state_PBs {
				data[i+1][1] = v.String()
			}
			refreshDiffs()
			diffs.Refresh()
		})

		wrRefresh := widget.NewButton("Refresh WRs", func() {
			var remoteWRs *lib.RemoteWRs
			var err error

			if !viper.GetBool("wrs.http.skip") {
				remoteWRs, err = lib.RetrieveWRsFromWeb(viper.GetString("wrs.http.url"))
				if err != nil {
					dialog.NewError(err, w).Show()
				}
			}

			if remoteWRs == nil {
				remoteWRs, err = lib.RetrieveWRsFromGoogle(
					viper.GetString("wrs.sheets.apikey"),
					viper.GetString("wrs.sheets.sheet_id"),
					viper.GetString("wrs.sheets.ranges"),
				)
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
			}

			state_WRs = remoteWRs.WRs

			for i, v := range state_WRs {
				data[i+1][2] = v.String()
			}
			refreshDiffs()
			diffs.Refresh()
		})

		openOpts := widget.NewButton("Settings", func() {
			opts.Show()
		})

		go func() {
			wrRefresh.Disable()
			pbRefresh.Disable()
			pbRefresh.OnTapped()
			wrRefresh.OnTapped()
			refreshDiffs()
			wrRefresh.Enable()
			pbRefresh.Enable()
		}()

		w.SetContent(container.NewBorder(container.NewHBox(pbRefresh, wrRefresh, openOpts), nil, nil, nil, diffs))

		w.ShowAndRun()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		conf, err := os.UserConfigDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(conf + "/nwscores")
		viper.AddConfigPath(".")

		viper.SetConfigType("toml")
		viper.SetConfigName("nwscores")

		viper.SetDefault("pbs.savefile", `C:\Users\user\AppData\LocalLow\Little Flag Software, LLC\Neon White\1234\savedata.dat`)
		viper.SetDefault("wrs.http.skip", false)
		viper.SetDefault("wrs.http.url", "https://nwscores.fuckhole.org/wrs.json")
		viper.SetDefault("wrs.sheets.apikey", "")
		viper.SetDefault("wrs.sheets.sheet_id", "1rG5WNRp4XBGxImwF4c0cj5oYbdIC4yMTpx45BU3cOLU")
		viper.SetDefault("wrs.sheets.ranges", `Rebirth!F5:F14;Killer Inside!F5:F14;Only Shallow!F5:F14;Boss Chapters!F5:F7;The Burn That Cures!F5:F14;Covenant!F5:F14;Reckoning!F5:F14;Benediction!F5:F14;Apocrypha!F5:F14;Boss Chapters!F18:F19;Thousand Pound Butterfly!F5:F14;Boss Chapters!F29:F30;Sidequests!E5:E12;Sidequests!U5:U12;Sidequests!M5:M12`)

		viper.SafeWriteConfig()
	}

	viper.SetEnvPrefix("NWSCORES")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
