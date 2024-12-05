/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/cybenc/cryptLocal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type FileInfo struct {
	SourceName  string
	DecryptName string
	IsDir       bool
}

// dirCmd represents the dir command
var dirCmd = &cobra.Command{
	Use:   "dir",
	Short: "目录相关操作",
	Long:  `目录相关操作`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flag := utils.CIns.CheckConfig()
		if !flag {
			logrus.Error("程序还没有初始化，请先初始化程序")
			// 退出程序
			os.Exit(0)
		}
		// 读取配置文件
		models, err := utils.CIns.ReadConfig()
		if err != nil {
			logrus.Error("读取配置文件失败")
			os.Exit(0)
		}
		utils.CipherInstance = &utils.Cipher{}
		utils.CipherInstance.NewCiper(models)
	},
}

var dirMkdirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "创建目录",
	Long:  `创建目录`,
	Run: func(cmd *cobra.Command, args []string) {
		// 判断参数是否为空
		if len(args) == 0 {
			logrus.Error("请输入目录名称")
			return
		}
		dirName := args[0]
		c := utils.CipherInstance.RcCipher
		targetDirName := c.EncryptDirName(dirName)
		// 创建目录
		os.Mkdir(targetDirName, 0755)
	},
}

var dirRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "重命名目录",
	Long:  `重命名目录`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logrus.Error("请输入源目录名称和目标目录名称")
			return
		}
		sourceDirName := args[0]
		targetDirName := args[1]
		logrus.Info("源目录名称：", sourceDirName, "目标目录名称：", targetDirName)
		// 解密源目录名称
		c := utils.CipherInstance.RcCipher
		sourceDcDirName, err := c.DecryptDirName(sourceDirName)
		if err != nil {
			logrus.Error("解密源目录名称失败")
			return
		}
		// 解密目标目录名称
		targetEcDirName := c.EncryptDirName(targetDirName)
		// 重命名目录
		os.Rename(sourceDcDirName, targetEcDirName)
		logrus.Info("目录重命名成功:", sourceDirName, "->", targetDirName)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出当前目录下的文件和目录",
	Long:  `列出当前目录下的文件和目录`,
	Run: func(cmd *cobra.Command, args []string) {
		// 列出当前目录下的文件和目录
		files, err := os.ReadDir(".")
		if err != nil {
			logrus.Error("读取当前目录失败")
			return
		}
		c := utils.CipherInstance.RcCipher
		file_info_list := make([]FileInfo, 0)
		for _, file := range files {
			if file.IsDir() {
				dir_name, err := c.DecryptDirName(file.Name())
				if err != nil {
					continue
				}
				file_info := FileInfo{
					SourceName:  file.Name(),
					DecryptName: dir_name,
					IsDir:       true,
				}
				file_info_list = append(file_info_list, file_info)
			} else {
				file_name, err := c.DecryptFileName(file.Name())
				if err != nil {
					continue
				}
				file_info := FileInfo{
					SourceName:  file.Name(),
					DecryptName: file_name,
					IsDir:       false,
				}
				file_info_list = append(file_info_list, file_info)
			}
		}
		for _, file_info := range file_info_list {
			if file_info.IsDir {
				fmt.Println("<DIR>\t", file_info.SourceName, "\t", file_info.DecryptName)
			} else {
				fmt.Println("     \t", file_info.SourceName, "\t", file_info.DecryptName)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(dirCmd)
	dirCmd.AddCommand(dirMkdirCmd)
	dirCmd.AddCommand(dirRenameCmd)
	dirCmd.AddCommand(listCmd)
}
