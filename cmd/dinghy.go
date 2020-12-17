package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armory-io/arm/pkg"
	dinghy_yaml "github.com/armory-io/dinghy/pkg/parsers/yaml"
	"github.com/armory/dinghy/pkg/cache"
	"github.com/armory/dinghy/pkg/dinghyfile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		_, err := dinghyRender(args)
		if err != nil {
			os.Exit(1)
		}
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
	renderCmd.Flags().String("local_modules", "", "local path to the dinghy local_module repository, if not specified local_module will be base path from dinghyfile")
	renderCmd.Flags().String("modules", "", "local path to the dinghy modules repository")
	renderCmd.Flags().String("rawdata", "", "optional rawdata json in case is needed")
	renderCmd.Flags().String("output", "", "output document rendered as it will be rendered by dinghy")
	renderCmd.Flags().String("type", "", "type of document to be rendered, supported types are: ['json','yaml']. Default value: json")
	dinghyCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(dinghyCmd)
	viper.BindPFlags(renderCmd.Flags())
}


func dinghyRender(args []string) (string, error) {

	//Init log
	log := initLog()
	log.Info("Checking dinghyfile")

	//Check that dinghyfile is entered
	var file string
	if len(args) > 0 {
		file = args[0]
	} else {
		log.Error("No dinghy file was entered, please refer to documentation or execute this command with --help")
		os.Exit(1)
	}

	//Get absolute path for dinghyfile
	absFile, err := filepath.Abs(file)
	if err != nil {
		log.Errorf("Invalid path for dinghyfile: %v", err)
		return "", fmt.Errorf("Invalid path for dinghyfile: %v", err)
	}

	//Separate directory path and filename
	repoFolder := fmt.Sprint(filepath.Dir(absFile))
	fileName := fmt.Sprint(filepath.Base(absFile))

	//Downloader will have the original dinghyfile directory and the dinghyfile name so we can read the dinghyfile
	//from the original location and then change it to the local_modules if some input parameter is there
	downloader := pkg.LocalDownloader{
		LocalModule: absFile,
		RepoFolder: repoFolder,
		DinghyfileName: fileName,
	}

	viper.SetDefault("type", "json")
	docType := viper.GetString("type")

	var unmarshaller dinghyfile.Unmarshaller
	var parser dinghyfile.Parser
	if docType == "json" {
		unmarshaller = &dinghyfile.DinghyJsonUnmarshaller{}
		parser = &dinghyfile.DinghyfileParser{}
	} else if docType == "yaml" {
		unmarshaller = &dinghy_yaml.DinghyYaml{}
		parser = &dinghy_yaml.DinghyfileYamlParser{}
	} else {
		log.Fatal(fmt.Sprintf("Invalid document type: %v", docType))
	}

	builder := &dinghyfile.PipelineBuilder{
		Downloader:      downloader,
		Depman:          cache.NewMemoryCache(),
		TemplateRepo:    viper.GetString("modules"),
		TemplateOrg:     "templateOrg",
		Logger:          log,
		Client:          PlankMock{},
		EventClient:     &dinghyfile.EventsTestClient{},
		Parser:          parser,
		DinghyfileName:  fileName,
		Ums:             []dinghyfile.Unmarshaller{unmarshaller},
	}

	//Process rawData and add it to the builder
	builder.PushRaw = processRawData(viper.GetString("rawdata"))

	log.Info("Parsing dinghyfile")
	builder.Parser.SetBuilder(builder)

	//If local_modules is entered we will change the repoFolder so all the files are read from that location
	//The dinghyfile is the only file that will not be read from there, downloader already has the information from
	//where to read
	var localModulesPath = viper.GetString("local_modules")
	if localModulesPath != "" {
		absFile = filepath.Dir(localModulesPath) + string(filepath.Separator) + filepath.Base(file)
	}

	repoFolder = fmt.Sprint(filepath.Dir(absFile))
	fileName = fmt.Sprint(filepath.Base(absFile))

	//Parse dinghyfile
	out, err := builder.Parser.Parse( "templateOrg", repoFolder, fileName, "", nil)

	if err != nil {
		log.Errorf("Parsing dinghyfile failed: %s", err )
		return "", fmt.Errorf("Parsing dinghyfile failed: %s", err )
	} else {
		log.Info("Parsed dinghyfile")
		
		isValid, errval := validateType(docType, out.Bytes())
		if errval != nil {
			log.Error(fmt.Sprintf("Validation for %v failed: %v", docType, errval))
		}

		if !isValid {
			log.Info("Output:\n")
			fmt.Println(out.String())
			errorMsg := fmt.Sprintf("The result is not a valid %v Object, please fix your dinghyfile", strings.ToUpper(docType))
			log.Error(errorMsg)
			return "", errors.New(errorMsg)
		}

		log.Info("Parsing final dinghyfile to struct for validation")
		d, err := dinghyfileStruct(builder, out)
		if err != nil {
			log.Errorf("Parsing to struct failed: %v", err)
			return "", fmt.Errorf("Parsing to struct failed: %v", err)
		}

		errValidation := builder.ValidatePipelines(d, out.Bytes())
		if errValidation != nil {
			log.Errorf("Final Dinghyfile failed validations, please correct them and retry. %v", errValidation)
			return "", fmt.Errorf("Final Dinghyfile failed validations, please correct them and retry. %v", errValidation)
		}

		errValidation = builder.ValidateAppNotifications(d, out.Bytes())
		if errValidation != nil {
			log.Errorf("Final Dinghyfile failed application notification validations, please correct them and retry. %v", errValidation)
			return "", fmt.Errorf("Final Dinghyfile failed application notification validations, please correct them and retry. %v", errValidation)
		}

		//Save file if output exists
		//Log output
		outIndent := indentType(docType, out.Bytes())
		buff:=bytes.NewBuffer(outIndent)
		saveOutputFile(viper.GetString("output"), buff)
		log.Info("Output:\n")
		fmt.Println(buff.String())
		log.Info(fmt.Sprintf("Final dinghyfile is a valid %v Object.", strings.ToUpper(docType)))

		return buff.String(), nil

	}
	return "", nil
}

func indentType(docType string, i []byte) []byte {
	switch  docType {
	case "json":
		var outIndent bytes.Buffer
		json.Indent(&outIndent, i, "", "  ")
		return outIndent.Bytes()
	}
	return i
}

func validateType(docType string, i []byte) (bool, error) {
	switch  docType {
	case "json":
		return json.Valid(i), nil
	case "yaml":
		var validate map[string]interface{}
		erryaml := yaml.Unmarshal(i, &validate)
		if erryaml != nil {
			return false, erryaml
		}
		return true, nil
	default:
		return false, errors.New(fmt.Sprintf("Invalid docType: %v", docType))
	}
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

func saveOutputFile(outputPath string, content *bytes.Buffer) {
	if outputPath != "" {
		log.Info("Saving output file")
		err := ioutil.WriteFile(outputPath, content.Bytes(), 0644)
		if err != nil {
			log.Error("Failed to save output file")
		}
	}
}

