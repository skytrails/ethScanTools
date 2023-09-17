package copy

import (
	"eth-scan/common/database"
	"eth-scan/common/storage"
	ext "eth-scan/config"
	"fmt"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/spf13/cobra"
	"log"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "copy",
		Short:   "copy db",
		Example: "eth-scan copy -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "copy", "c", "config/settings.yml", "Start server with provided configuration file")
}

func setup() {
	// 注入配置扩展项
	config.ExtendConfig = &ext.ExtConfig
	//1. 读取配置
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		storage.Setup,
	)

	usageStr := `starting eth scan...`
	log.Println(usageStr)
}

func run() {
	dbList := sdk.Runtime.GetDb()
	for _, db := range dbList {
		go process(db)
	}

	fmt.Println("All goroutines completed.")
	select {}

}
