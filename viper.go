package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version int

var rootCmd = &cobra.Command{
	Use: "test",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\nVersion: %v\n\n", viper.Get("about.version"))
	},
}

func init() {
	rootCmd.PersistentFlags().IntVar(&version, "about-version", -1, "version number")
	// flag "about-version" maps to viper key "about.version"
	viper.BindPFlag("about.version", rootCmd.PersistentFlags().Lookup("about-version"))
	// load config
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	// ENV "ABOUT_VERSION" will map to viper key "about.version"
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func runViper() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}
}
