package cmd

import (
	"bytes"
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
	Run:   func(cmd *cobra.Command, args []string) {
		dinghyRender(args)
	},
}

func processRawData(rawDataPath string) map[string]interface{} {
	if rawDataPath != "" {
		log.Info("Reading rawdata file")
		rawData, err := ioutil.ReadFile(rawDataPath)
		if err != nil {
			log.Error("Invalid rawdata json file path")
		} else {
			log.Info("Parsing rawdata json")
			rawPushData := make(map[string]interface{})
			err = json.Unmarshal([]byte(rawData), &rawPushData)
			if err != nil {
				log.Error("Invalid rawData json file")
			}
			return rawPushData
		}
	}
	return nil
}

func init() {
	renderCmd.Flags().String("modules", "", "local path to the dinghy modules repository")
	renderCmd.Flags().String("rawdata", "", "optional rawdata json in case is needed")
	renderCmd.Flags().String("output", "", "output json rendered as it will be rendered by dinghy")
	dinghyCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(dinghyCmd)
	viper.BindPFlags(renderCmd.Flags())
}


func dinghyRender(args []string) string {

	log := initLog()
	log.Info("Checking dinghyfile")

	var file string
	if len(args) > 0 {
		file = args[0]
	} else {
		log.Fatal("No dinghy file was entered, please refer to documentation or execute this command with --help")
	}

	downloader := pkg.LocalDownloader{}
	builder := &dinghyfile.PipelineBuilder{
		Downloader:      downloader,
		Depman:          cache.NewMemoryCache(),
		TemplateRepo:    viper.GetString("modules"),
		TemplateOrg:     "templateOrg",
		Logger:          log.WithField("arm-cli-test", ""),
		Client:          plank.New(),
		EventClient:     &dinghyfile.EventsTestClient{},
		Parser:          &dinghyfile.DinghyfileParser{},
		DinghyfileName:  filepath.Base(file),
	}

	//Process rawData and add it to the builder
	builder.PushRaw = processRawData(viper.GetString("rawdata"))

	log.Info("Parsing dinghyfile")

	builder.Parser.SetBuilder(builder)

	out, err := builder.Parser.Parse("", "", file, "", nil)

	if err != nil {
		log.Fatalf("Parsing dinghyfile failed: %s", err )
	} else {

		log.Info("Parsed dinghyfile")

		if !json.Valid(out.Bytes()){
			log.Info("Output:\n")
			fmt.Println(out.String())
			log.Fatal("The result is not a valid JSON Object, please fix your dinghyfile")
		} else {
			var outIndent bytes.Buffer
			json.Indent(&outIndent, out.Bytes(), "", "  ")
			//Save file if output exists
			saveOutputFile(viper.GetString("output"), outIndent)
			//Log output
			log.Info("Output:\n")
			fmt.Println(outIndent.String())
			log.Info("Final dinghyfile is a valid JSON Object.")
			return outIndent.String()
		}
	}
	return out.String()
}

func saveOutputFile(outputPath string, content bytes.Buffer) {
	if outputPath != "" {
		log.Info("Saving output file")
		err := ioutil.WriteFile(outputPath, content.Bytes(), 0644)
		if err != nil {
			log.Error("Failed to save output file")
		}
	}
}