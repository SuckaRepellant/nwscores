//go:build gui

package cmd

import (
	"io/ioutil"
	"log"
	"nwscores/lib"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RunGui() {
	a := app.New()

	dw := DiffWindow(a)

	if viper.GetBool("gui.firstrun") {
		// Attempt automatic discovery of savefile for windows users
		if runtime.GOOS == "windows" {
			var potentials []string

			nw := filepath.Join(os.Getenv("APPDATA"), "..", "LocalLow") + `\Little Flag Software, LLC\Neon White`

			items, err := ioutil.ReadDir(nw)
			if err == nil {
				for _, item := range items {
					if item.IsDir() {
						potentials = append(potentials, nw+`\`+item.Name()+`\savedata.dat`)
					}
				}
			}

			log.Println(potentials)

			frw := FirstRunWindow(a, potentials)
			frw.Show()
			frw.SetOnClosed(func() {
				if !viper.GetBool("gui.firstrun") {
					dw.Show()
				}
			})
		} else {
			frw := FirstRunWindow(a, []string{})
			frw.Show()
			frw.SetOnClosed(func() {
				if !viper.GetBool("gui.firstrun") {
					dw.Show()
				}
			})
		}
	} else {
		dw.Show()
	}

	a.Run()

}

func FirstRunWindow(a fyne.App, potentials []string) fyne.Window {
	w := a.NewWindow("nwscores - First Run")
	w.Resize(fyne.NewSize(800, 300))

	pathField := widget.NewEntry()
	pathField.SetPlaceHolder(DEFAULT_SAVE_PATH)

	openButton := widget.NewButton("Open", func() {
		viper.Set("pbs.savefile", pathField.Text)
		viper.Set("gui.firstrun", false)
		viper.WriteConfig()
		w.Close()
	})

	discoverSelect := widget.NewSelect(potentials, func(new string) { pathField.SetText(new) })

	w.SetContent(container.NewBorder(
		widget.NewLabel("Which one of these steam IDs is yours?\nYou can also enter a savefile.dat path manually if none are."),
		container.NewVBox(pathField, openButton),
		nil,
		nil,
		discoverSelect,
	))

	return w
}

func DiffWindow(a fyne.App) fyne.Window {
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

	var data = [][]string{
		{"Name", "PB", "WR", "Diff"},
	}

	var state_PBs, state_WRs []time.Duration

	for _, lv := range lib.GetLevels() {
		//                          Name PB       WR       Diff
		data = append(data, []string{lv, "0.000s", "0.000s", "0ms"})
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
		state_PBs, err = lib.RetrievePBsFromDisk(viper.GetString("pbs.savefile"))
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

	if !viper.GetBool("gui.firstrun") {
		go func() {
			wrRefresh.Disable()
			pbRefresh.Disable()
			pbRefresh.OnTapped()
			wrRefresh.OnTapped()
			refreshDiffs()
			wrRefresh.Enable()
			pbRefresh.Enable()
		}()
	}

	w.SetContent(container.NewBorder(container.NewHBox(pbRefresh, wrRefresh, openOpts), nil, nil, nil, diffs))

	return w
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Show GUI",
	Run: func(cmd *cobra.Command, args []string) {
		RunGui()
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
