package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/armory-io/arm/pkg"
	"github.com/armory/dinghy/pkg/cache"
	"github.com/armory/dinghy/pkg/dinghyfile"
	"github.com/armory/plank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"path/filepath"
)

// dinghyCmd represents the dinghy command
var dinghyCmd = &cobra.Command{
	Use:   "dinghy",
	Short: "Run dinghy subcommands",
}

var renderCmd = &cobra.Command{
	Use:   "render [dinghyfile]",
	Short: "render a dinghyfile",
	Run: func(cmd *cobra.Command, args []string) {

		log := initLog()
		log.Debug("Checking dinghyfile")

		file := "dinghyfile"
		if len(args) > 0 {
			file = args[0]
		}

		downloader := pkg.LocalDownloader{}
		builder := &dinghyfile.PipelineBuilder{
			Downloader:      downloader,
			Depman:          cache.NewMemoryCache(),
			TemplateRepo:    viper.GetString("modules"),
			TemplateOrg:     "/",
			Logger:          log.WithField("arm-cli-test", ""),
			Client:          plank.New(),
			EventClient:     &dinghyfile.EventsTestClient{},
			Parser:          &dinghyfile.DinghyfileParser{},
			DinghyfileName:  filepath.Base(file),
		}

		rawDataPath := viper.GetString("rawdata")
		if rawDataPath != "" {
			log.Debug("Reading rawdata file")
			rawData, err := ioutil.ReadFile(rawDataPath)
			if err != nil {
				log.Error("Invalid rawdata json file path")
			} else {
				log.Debug("Parsing rawdata json")
				rawPushData := make(map[string]interface{})
				err = json.Unmarshal([]byte(rawData), &rawPushData)
				if err != nil {
					log.Error("Invalid rawData json file")
				}
				builder.PushRaw = rawPushData
			}
		}

		log.Debug("Parsing dinghyfile")

		builder.Parser.SetBuilder(builder)

		out, err := builder.Parser.Parse("", "", file, "", nil)

		if err != nil {
			log.Errorf("Parsing dinghyfile failed: %s", err )
		} else {
			log.Debug("Parsed dinghyfile")
			log.Info("Output:\n")
			fmt.Println(out.String())

			log.Debug("Validating output json.")
			if !json.Valid(out.Bytes()){
				log.Fatal("The result is not a valid JSON object, please fix your dinghyfile")
			}
		}
	},
}

func init() {
	renderCmd.Flags().String("rawdata", "", "optional rawdata json in case is needed")
	renderCmd.Flags().String("modules", "", "local path to the dinghy modules repository")
	dinghyCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(dinghyCmd)
	viper.BindPFlags(renderCmd.Flags())
}
