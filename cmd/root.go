package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"net/http"
	"errors"
	"encoding/json"
)


const ENABLE_FLAG="check"
const VERSION_URL="https://get.armory.io/arm/version.json"

var (
	currentVersion = "HEAD"
	enableVersionCheck = "offbydefault"
	UPGRADE_VERSION_ERROR=errors.New("there is a new version available")
	LogLevel string
	log *logrus.Logger
)

type VersionDesc struct {
	Version string `json:"version"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "arm",
	Short: "The Armory Platform CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if enableVersionCheck == ENABLE_FLAG {
			if err, version := checkVersion(); err != nil {
				if err == UPGRADE_VERSION_ERROR {
					log.Warn(`Client version is %s but a newer version (%s) is available. Please upgrade to the latest version!`, currentVersion, version)
				} else {
					log.Warn("There was a problem verifying the arm version number. Your client may be out of date.")
				}
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkVersion() (error, string){

	resp, err := http.Get(VERSION_URL)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()
	v := VersionDesc{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return err, ""
	}

	if currentVersion != v.Version {
		return UPGRADE_VERSION_ERROR, v.Version
	}

	return nil, ""

}

func initLog() *logrus.Logger {
	log = logrus.New()
	var level logrus.Level
	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	return log
}


func init() {
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "loglevel", "l", "info", "log level")
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}
