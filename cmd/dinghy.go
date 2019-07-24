package cmd

import (
	"fmt"
	"github.com/armory-io/arm/pkg"
	"github.com/armory/dinghy/pkg/cache"
	"github.com/armory/dinghy/pkg/dinghyfile"
	"github.com/armory/plank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		log.Debug("checking dinghyfile")

		file := "dinghyfile"
		if len(args) > 0 {
			file = args[0]
		}

		downloader := pkg.LocalDownloader{}
		builder := &dinghyfile.PipelineBuilder{
			Downloader:   downloader,
			Depman:       cache.NewMemoryCache(),
			TemplateRepo: viper.GetString("modules"),
			TemplateOrg:  "",
			Logger:       log.WithField("arm-cli-test", ""),
			Client:       plank.New(),
			EventClient:  &dinghyfile.EventsTestClient{},
			Parser:       &dinghyfile.DinghyfileParser{},
		}
		builder.Parser.SetBuilder(builder)

		log.Debug("parsing dinghyfile")
		out, err := builder.Parser.Parse("", "", file, "", nil)
		if err != nil {
			log.Error(err)
		}
		log.Info("Output:\n")
		fmt.Println(out.String())
	},
}

func init() {
	renderCmd.Flags().String("modules", "", "local path to the dinghy modules repository")
	dinghyCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(dinghyCmd)
	viper.BindPFlags(renderCmd.Flags())
}
