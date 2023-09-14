package scan

import (
	"eth-scan/cmd/scan/models"
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
		Use:     "scan",
		Short:   "scan eth balance",
		Example: "eth-scan scan -c config/settings.yml",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "scan", "c", "config/settings.yml", "Start server with provided configuration file")
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
	ethBlockScanRecord := models.EthBlockScanRecord{}
	var data []models.EthBlockScanRecord
	fmt.Println("test1")
	for _, db := range dbList {
		fmt.Println("test2")
		err := ethBlockScanRecord.GetToDoList(db, &data)
		if err != nil {
			return
		}
		for _, d := range data {
			fmt.Printf("id: %d, start: %d, cur: %d, end: %d\n", d.Id, d.StartBlock, d.CurBlock, d.EndBlock)
			go scan(db, d.Id, d.CurBlock, d.EndBlock)
		}

		//go scanBalance(db)
	}

	fmt.Println("All goroutines completed.")
	select {}

}
