/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/cybenc/cryptLocal/utils"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理app相关配置",
	Long:  `清理app相关配置(危险操作⚠️)`,
	Run: func(cmd *cobra.Command, args []string) {
		// 确认是否继续
		prompt := promptui.Prompt{
			Label:       "该操作会删除所有app相关配置,请选择是否继续?(y/n)",
			HideEntered: true,
			Validate: func(input string) error {
				if input != "y" && input != "n" {
					return fmt.Errorf("输入错误,请输入y或n")
				}
				return nil
			},
		}
		result, err := prompt.Run()
		if err != nil {
			return
		}
		if result != "y" {
			return
		}
		// 执行清理操作
		logrus.Info("cleaning...,delete app dir:", utils.CIns.AppDir)
		utils.CIns.DeleteAppDir()
		logrus.Info("clean ", utils.CIns.AppDir, " success")
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
