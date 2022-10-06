package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	cobra.OnInitialize(initConfig)
}

const DEFAULT_SAVE_PATH = `C:\Users\user\AppData\LocalLow\Little Flag Software, LLC\Neon White\1234\savedata.dat`

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

		viper.SetDefault("pbs.savefile", DEFAULT_SAVE_PATH)
		viper.SetDefault("wrs.http.skip", false)
		viper.SetDefault("wrs.http.url", "https://nwscores.fuckhole.org/wrs.json")
		viper.SetDefault("wrs.sheets.apikey", "")
		viper.SetDefault("wrs.sheets.sheet_id", "1rG5WNRp4XBGxImwF4c0cj5oYbdIC4yMTpx45BU3cOLU")
		viper.SetDefault("wrs.sheets.ranges", `Rebirth!F5:F14;Killer Inside!F5:F14;Only Shallow!F5:F14;Boss Chapters!F5:F7;The Burn That Cures!F5:F14;Covenant!F5:F14;Reckoning!F5:F14;Benediction!F5:F14;Apocrypha!F5:F14;Boss Chapters!F18:F19;Thousand Pound Butterfly!F5:F14;Boss Chapters!F29:F30;Sidequests!E5:E12;Sidequests!U5:U12;Sidequests!M5:M12`)
		viper.SetDefault("gui.firstrun", true)

		confFolder, _ := os.UserConfigDir()
		confFolder = confFolder + "/nwscores"
		os.MkdirAll(confFolder, 0600)

		viper.SafeWriteConfig()
	}

	viper.SetEnvPrefix("NWSCORES")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
