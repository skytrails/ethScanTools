package cmd

import (
	"errors"
	"eth-scan/cmd/app"
	"eth-scan/cmd/scan"
	"eth-scan/common/global"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"os"

	"github.com/spf13/cobra"

	"eth-scan/cmd/api"
	"eth-scan/cmd/config"
	"eth-scan/cmd/migrate"
	"eth-scan/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:          "eth-scan",
	Short:        "eth-scan",
	SilenceUsage: true,
	Long:         `eth-scan`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`eth-scan `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	usageStr1 := `也可以参考 https://doc.go-admin.dev/guide/ksks 的相关内容`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(scan.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(app.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
