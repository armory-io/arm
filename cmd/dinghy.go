package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/armory-io/arm/pkg"
	"github.com/armory/dinghy/pkg/cache"
	"github.com/armory/dinghy/pkg/dinghyfile"
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
		Client:          nil,
		EventClient:     &dinghyfile.EventsTestClient{},
		Parser:          &dinghyfile.DinghyfileParser{},
		DinghyfileName:  filepath.Base(file),
		Ums:             []dinghyfile.Unmarshaller{&dinghyfile.DinghyJsonUnmarshaller{}},
	}

	//Process rawData and add it to the builder
	builder.PushRaw = processRawData(viper.GetString("rawdata"))

	log.Info("Parsing dinghyfile")

	builder.Parser.SetBuilder(builder)

	absFile, err := filepath.Abs(file)
	if err != nil {
		log.Fatal("Invalid path for dinghyfile: %v", err)
	}
	repoFolder := fmt.Sprint(filepath.Dir(absFile))
	fileName := fmt.Sprint(filepath.Base(absFile))
	out, err := builder.Parser.Parse( "templateOrg", repoFolder, fileName, "", nil)

	if err != nil {
		log.Fatalf("Parsing dinghyfile failed: %s", err )
	} else {

		log.Info("Parsed dinghyfile")

		if !json.Valid(out.Bytes()){
			log.Info("Output:\n")
			fmt.Println(out.String())
			log.Fatal("The result is not a valid JSON Object, please fix your dinghyfile")
		} else {
			log.Info("Parsing final dinghyfile to struct for validation")
			d, err := dinghyfileStruct(builder, out)
			if err != nil {
				log.Fatalf("Parsing to struct failed: %v", err)
			}

			errValidation := builder.ValidatePipelines(d, out.Bytes())
			if errValidation != nil {
				log.Fatal("Final Dinghyfile failed validations, please correct them and retry")
			}

			//Save file if output exists
			//Log output
			var outIndent bytes.Buffer
			json.Indent(&outIndent, out.Bytes(), "", "  ")
			saveOutputFile(viper.GetString("output"), outIndent)
			log.Info("Output:\n")
			fmt.Println(outIndent.String())
			log.Info("Final dinghyfile is a valid JSON Object.")

			return outIndent.String()
		}
	}
	return out.String()
}

func dinghyfileStruct(builder *dinghyfile.PipelineBuilder, out *bytes.Buffer) (dinghyfile.Dinghyfile, error) {
	d := dinghyfile.NewDinghyfile()
	parseErrs := 0
	var err error
	for _, ums := range builder.Ums {
		log.Info("Parsing Dinghyfile to struct for validation")
		if errmarshal := ums.Unmarshal(out.Bytes(), &d); errmarshal != nil {
			err = errmarshal
			log.Warnf("Cannot create Dinghyfile struct because of malformed syntax: %s", err.Error())
			parseErrs++
			continue
		}
	}
	return d, err
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

