package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nwscores",
	Short: "Extract your Neon White scores from your save file",
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

		viper.SetDefault("savefile", `C:\xyz\abc\savedata.dat`)
		viper.SetDefault("sheets.apikey", "")
		viper.SetDefault("sheets.sheet_id", "1rG5WNRp4XBGxImwF4c0cj5oYbdIC4yMTpx45BU3cOLU")
		viper.SetDefault("sheets.ranges", `Rebirth!F5:F14;Killer Inside!F5:F14;Only Shallow!F5:F14;Boss Chapters!F5:F7;The Burn That Cures!F5:F14;Covenant!F5:F14;Reckoning!F5:F14;Benediction!F5:F14;Apocrypha!F5:F14;Boss Chapters!F18:F19;Thousand Pound Butterfly!F5:F14;Boss Chapters!F29:F30;Sidequests!E5:E12;Sidequests!U5:U12;Sidequests!M5:M12`)

		viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
