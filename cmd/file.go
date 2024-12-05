/*
Copyright © 2024 NAME HERE <cybertm@qq.com>
*/
package cmd

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/cybenc/cryptLocal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "文件相关命令",
	Long:  `文件相关命令 加密、解密、添加、删除、恢复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flag := utils.CIns.CheckConfig()
		if !flag {
			logrus.Error("配置文件不存在")
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

var addCmd = &cobra.Command{
	Use:   "add [source_file_path]",
	Short: "添加文件",
	Long:  `添加文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("请输入文件路径")
			return
		}
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			logrus.Error("打开文件失败")
			return
		}
		defer file.Close()
		// 将 *os.File 转换为 io.Reader
		var reader io.Reader = file

		// 使用 bufio.NewReader 来包装 io.Reader，以便更方便地读取
		bufferedReader := bufio.NewReader(reader)
		//加密文件
		cipherData, err := utils.CipherInstance.RcCipher.EncryptData(bufferedReader)
		if err != nil {
			logrus.Error("加密文件失败")
			return
		}

		// 使用 filepath.Base 来获取文件名（包含扩展名）
		fileName := filepath.Base(filePath)
		cipherFileName := utils.CipherInstance.RcCipher.EncryptFileName(fileName)
		if cipherFileName == "" {
			logrus.Error("加密文件名为空")
			return
		}
		// 保存加密文件
		cipherFile, err := os.Create(cipherFileName)
		if err != nil {
			logrus.Error("创建加密文件失败")
			return
		}
		defer cipherFile.Close()
		// 将加密数据写入加密文件
		bytes, err := io.ReadAll(cipherData)
		if err != nil {
			logrus.Error("写入加密文件失败")
			return
		}
		// 保存文件
		cipherFile.Write(bytes)

	},
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "删除文件",
	Long:  `删除文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("请输入文件路径")
			return
		}
		filePath := args[0]
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			logrus.Error("文件不存在")
			return
		}
		// 解密文件名
		cipherFileName, err := utils.CipherInstance.RcCipher.DecryptFileName(filePath)
		if err != nil || cipherFileName == "" {
			logrus.Error("解密文件名失败")
			return
		}
		// 删除加密文件
		err = os.Remove(cipherFileName)
		if err != nil {
			logrus.Error("删除加密文件失败")
			return
		}
	},
}

var recoverCmd = &cobra.Command{
	Use:   "recover",
	Short: "恢复文件",
	Long:  `恢复文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Error("请输入文件路径")
			return
		}
		fp := args[0]
		FileInfo, err := os.Stat(fp)
		if os.IsNotExist(err) {
			logrus.Error("文件不存在")
			return
		}
		if FileInfo.IsDir() {
			logrus.Error("请输入文件路径")
			return
		}
		// 判断是否是文件
		file, err := os.Open(fp)
		if err != nil {
			logrus.Error("打开文件失败")
			return
		}
		defer file.Close()

		// 解密文件名
		cipherFileName, err := utils.CipherInstance.RcCipher.DecryptFileName(filepath.Base(fp))
		if err != nil || cipherFileName == "" {
			logrus.Error("解密文件名失败")
			return
		}
		logrus.Info("解密文件名:", cipherFileName)

		// 解密文件
		cipherData, err := utils.CipherInstance.RcCipher.DecryptData(file)
		if err != nil {
			logrus.Info("解密文件失败")
			return
		}
		// 保存解密文件
		output := cmd.Flags().Lookup("output").Value.String()
		if output == "" {
			output = filepath.Join(filepath.Dir(fp), filepath.Base(cipherFileName))
		} else {
			output, err = filepath.Abs(output)
			if err != nil {
				logrus.Info("指定输出路径错误")
				return
			}
		}
		logrus.Info("恢复文件保存路径:", output)
		outputFile, err := os.Create(output)
		if err != nil {
			logrus.Info("创建解密文件失败")
			return
		}
		defer outputFile.Close()
		// 将解密数据写入解密文件
		bytes, err := io.ReadAll(cipherData)
		if err != nil {
			logrus.Info("写入解密文件失败")
			return
		}
		outputFile.Write(bytes)

	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.AddCommand(addCmd)
	fileCmd.AddCommand(removeCmd)
	fileCmd.AddCommand(recoverCmd)

	recoverCmd.Flags().StringP("output", "o", "", "指定恢复文件保存路径")
}
