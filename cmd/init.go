package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cybenc/cryptLocal/common"
	"github.com/cybenc/cryptLocal/models"
	"github.com/cybenc/cryptLocal/utils"
	"github.com/manifoldco/promptui"
	"github.com/rclone/rclone/fs/config/obscure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ApiCmdConfig struct {
	ApiUrl   string
	UserName string
	Password string
	OptCode  string
}

type InputCmdConfig struct {
	Password                string
	Salt                    string
	FileNameEncryption      string
	DirectoryNameEncryption string
	FileNameEncoding        string
	Suffix                  string
}

var apiCmdConfig = ApiCmdConfig{}
var inputCmdConfig = InputCmdConfig{}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化加密配置信息",
	Long:  `初始化加密配置信息,可以通过api或输入的方式进行初始化`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		flag := utils.CIns.CheckConfig()
		if flag {
			logrus.Error("请不要重复初始化,配置不一致会导致加密文件损坏。")
			// 退出程序
			os.Exit(0)
		}
	},
}

func api_promot() {
	if apiCmdConfig.ApiUrl == "" {
		prompt := promptui.Prompt{
			Label: "请输入api地址",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("api地址不能为空")
				}
				if !strings.HasPrefix(input, "http") {
					return fmt.Errorf("api地址必须以http或https开头")
				}
				return nil
			},
			HideEntered: true,
		}
		result, _ := prompt.Run()
		apiCmdConfig.ApiUrl = result
	}
	if apiCmdConfig.UserName == "" {
		prompt := promptui.Prompt{
			Label: "请输入用户名",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("用户名不能为空")
				}
				return nil
			},
			HideEntered: true,
		}
		result, _ := prompt.Run()
		apiCmdConfig.UserName = result
	}
	if apiCmdConfig.Password == "" {
		prompt := promptui.Prompt{
			Label: "请输入密码",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("密码不能为空")
				}
				return nil
			},
			Mask:        '*',
			HideEntered: true,
		}
		result, _ := prompt.Run()
		apiCmdConfig.Password = result
	}
	if apiCmdConfig.OptCode == "" {
		prompt := promptui.Prompt{
			Label:       "请输入两步验证码(可选)",
			HideEntered: true,
		}
		result, _ := prompt.Run()
		apiCmdConfig.OptCode = result
	}
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "以api方式初始化",
	Run: func(cmd *cobra.Command, args []string) {
		// 输入配置信息
		api_promot()
		alist := utils.Alist{
			Url:      apiCmdConfig.ApiUrl,
			Password: apiCmdConfig.Password,
			UserName: apiCmdConfig.UserName,
			OptCode:  apiCmdConfig.OptCode,
		}
		alist.Login()
		config, err := alist.GenerateConfigByApi()
		if err != nil {
			logrus.Error(err)
			return
		}
		if config == nil {
			fmt.Println("获取配置失败,请检查相关配置")
			return
		}

		GenerateConfig(*config)
		logrus.Info("cryptoLocal初始化成功")

	},
}

func input_promot() {
	if inputCmdConfig.Password == "" {
		prompt := promptui.Prompt{
			Label: "请输入密码",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("密码不能为空")
				}
				return nil
			},
			Mask:        '*',
			HideEntered: true,
		}
		result, _ := prompt.Run()
		inputCmdConfig.Password = result
	}
	if inputCmdConfig.Salt == "" {
		prompt := promptui.Prompt{
			Label:       "请输入盐(可选)",
			Mask:        '*',
			HideEntered: true,
		}
		result, _ := prompt.Run()
		inputCmdConfig.Salt = result
	}

	if inputCmdConfig.FileNameEncryption == "" {
		prompt := promptui.Select{
			Label:        "请选择文件名加密方式",
			Items:        []string{"off", "standard", "obfuscate"},
			HideSelected: true,
		}
		_, result, _ := prompt.Run()
		inputCmdConfig.FileNameEncryption = result
	}
	if inputCmdConfig.DirectoryNameEncryption == "" {
		prompt := promptui.Select{
			Label:        "请选择目录名是否加密",
			Items:        []string{"true", "false"},
			HideSelected: true,
		}
		_, result, _ := prompt.Run()
		inputCmdConfig.DirectoryNameEncryption = result
	}
	if inputCmdConfig.FileNameEncoding == "" {
		prompt := promptui.Select{
			Label:        "请选择文件名编码",
			Items:        []string{"base64", "base32", "base32768"},
			HideSelected: true,
		}
		_, result, _ := prompt.Run()
		inputCmdConfig.FileNameEncoding = result
	}

}

var inputCmd = &cobra.Command{
	Use:   "input",
	Short: "通过输入的方式初始化配置信息",
	Run: func(cmd *cobra.Command, args []string) {
		// 输入配置信息
		input_promot()
		if !strings.HasPrefix(inputCmdConfig.Password, common.ObfuscatedPrefix) {
			temp, err := obscure.Obscure(inputCmdConfig.Password)
			if err != nil {
				fmt.Println("密码加密失败,请检查相关配置")
				return
			}
			inputCmdConfig.Password = common.ObfuscatedPrefix + temp
		}
		if inputCmdConfig.Salt != "" && !strings.HasPrefix(inputCmdConfig.Salt, common.ObfuscatedPrefix) {
			temp, err := obscure.Obscure(inputCmdConfig.Salt)
			if err != nil {
				fmt.Println("盐加密失败,请检查相关配置")
				return
			}
			inputCmdConfig.Salt = common.ObfuscatedPrefix + temp
		}
		// 构建配置
		config := models.EncrtptConfig{
			Password:                inputCmdConfig.Password,
			Password2:               inputCmdConfig.Salt,
			FileNameEncryption:      inputCmdConfig.FileNameEncryption,
			DirectoryNameEncryption: inputCmdConfig.DirectoryNameEncryption,
			FileNameEncoding:        inputCmdConfig.FileNameEncoding,
			Suffix:                  inputCmdConfig.Suffix,
			PassBadBlocks:           "",
		}
		GenerateConfig(config)
		logrus.Info("cryptoLocal初始化成功")
	},
}

func GenerateConfig(config models.EncrtptConfig) {
	// 判断密码和盐是否已经加密

	// 转换为json字符串
	jsonstr, _ := utils.StructToPrettyJSON(config)
	// 获取配置文件路径
	configpath := utils.CIns.GetAppConfigFile()
	// 写入配置文件
	os.WriteFile(configpath, []byte(jsonstr), 0644)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(apiCmd)
	initCmd.AddCommand(inputCmd)

	apiCmd.Flags().StringVarP(&apiCmdConfig.ApiUrl, "url", "u", "", "api url")
	apiCmd.Flags().StringVarP(&apiCmdConfig.UserName, "username", "n", "", "用户名")
	apiCmd.Flags().StringVarP(&apiCmdConfig.Password, "password", "p", "", "密码")
	apiCmd.Flags().StringVarP(&apiCmdConfig.OptCode, "optcode", "o", "", "两步验证码,可选")

	inputCmd.Flags().StringVarP(&inputCmdConfig.Password, "password", "p", "", "密码")
	inputCmd.Flags().StringVarP(&inputCmdConfig.Salt, "salt", "s", "", "盐")
	inputCmd.Flags().StringVarP(&inputCmdConfig.FileNameEncryption, "filenameencryption", "f", "", "文件名加密算法")
	inputCmd.Flags().StringVarP(&inputCmdConfig.DirectoryNameEncryption, "directorynameencryption", "d", "", "目录名加密算法")
	inputCmd.Flags().StringVarP(&inputCmdConfig.FileNameEncoding, "filenameencoding", "e", "", "文件名编码")
	inputCmd.Flags().StringVarP(&inputCmdConfig.Suffix, "suffix", "x", "", "文件后缀")

}
