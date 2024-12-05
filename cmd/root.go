package cmd

import (
	"fmt"

	"github.com/cybenc/cryptLocal/common"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: common.Version,
	Short:   "alist本地加解密工具",
	Long:    `alist本地加解密工具`,
}

func Execute() {
	rootCmd.Use = common.ExecuteName
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		logrus.Error(err)
	}
}
