package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/armory-io/dinghy/pkg/dinghyfile"
	"github.com/armory-io/dinghy/pkg/cache"
	"github.com/armory-io/arm/pkg"
	"github.com/spf13/viper"
)

// dinghyCmd represents the dinghy command
var dinghyCmd = &cobra.Command{
	Use:   "dinghy",
	Short: "Run dinghy subcommands",
}

var renderCmd = &cobra.Command{
	Use: "render [dinghyfile]",
	Short: "render a dinghyfile",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetOutput(ioutil.Discard)
		fmt.Println(args)

		file := "dinghyfile"
		if len(args) > 0 {
			file = args[0]
		}

		downloader := pkg.LocalDownloader{}
		builder := dinghyfile.PipelineBuilder{
			Downloader:   downloader,
			Depman:       cache.NewMemoryCache(),
			TemplateRepo: viper.GetString("modules"),
			TemplateOrg:  "",
		}

		out, err := builder.Render("", "", file, nil)
		if err != nil {
			log.Print(err)
                }

		fmt.Println(out.String())
	},
}

func init() {
	renderCmd.Flags().String("modules", "", "local path to the dinghy modules repository")
	dinghyCmd.AddCommand(renderCmd)
	rootCmd.AddCommand(dinghyCmd)
	viper.BindPFlags(renderCmd.Flags())
}
