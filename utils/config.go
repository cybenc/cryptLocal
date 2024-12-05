package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cybenc/cryptLocal/common"
	"github.com/cybenc/cryptLocal/models"
	"github.com/sirupsen/logrus"
)

var CIns = Config{}

type Config struct {
	HomeDir       string
	AppSeriesDir  string
	AppDir        string
	AppConfigDir  string
	AppConfigFile string
	LogDir        string
}

func (c *Config) CheckConfig() bool {
	// TODO: initialize config
	configfile := c.GetAppConfigFile()
	configDir := filepath.Dir(configfile)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}
	// 判断配置文件是否存在,不存在则返回false
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *Config) ReadConfig() (*models.EncrtptConfig, error) {
	// 读取配置文件
	file, err := os.Open(c.AppConfigFile)
	if err != nil {
		logrus.Error("无法打开文件")
		return nil, err
	}
	defer file.Close()

	// 读取文件内容
	byteValue, err := io.ReadAll(file)
	if err != nil {
		logrus.Error("无法读取文件")
		return nil, err
	}
	var model models.EncrtptConfig
	err = json.Unmarshal(byteValue, &model)
	if err != nil {
		logrus.Error("配置文件格式错误")
		return nil, err
	}
	return &model, nil
}

func (c *Config) GetHomeDir() string {
	if c.HomeDir != "" {
		return c.HomeDir
	}
	usr, err := user.Current()
	if err != nil {
		panic(fmt.Sprintf("获取当前用户目录失败: %v", err))
	}
	c.HomeDir = usr.HomeDir
	return c.HomeDir
}

func (c *Config) GetAppSeriesDir() string {
	if c.AppSeriesDir != "" {
		return c.AppSeriesDir
	}
	c.AppSeriesDir = filepath.Join(c.GetHomeDir(), common.SeriesName)
	return c.AppSeriesDir
}

func (c *Config) GetAppDir() string {
	if c.AppDir != "" {
		return c.AppDir
	}
	c.AppDir = filepath.Join(c.GetAppSeriesDir(), common.AppName)
	return c.AppDir
}

func (c *Config) DeleteAppDir() error {
	// 递归删除目录及其子目录
	return os.RemoveAll(c.GetAppDir())
}

func (c *Config) GetAppConfigDir() string {
	if c.AppConfigDir != "" {
		return c.AppConfigDir
	}
	c.AppConfigDir = filepath.Join(c.GetAppDir(), "config")
	return c.AppConfigDir
}

func (c *Config) GetAppConfigFile() string {
	if c.AppConfigFile != "" {
		return c.AppConfigFile
	}
	c.AppConfigFile = filepath.Join(c.GetAppConfigDir(), "config.json")
	return c.AppConfigFile
}

func (c *Config) GetLogDir() string {
	if c.LogDir != "" {
		return c.LogDir
	}
	homeDir := filepath.Join(c.GetAppDir(), "logs")
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		os.MkdirAll(homeDir, 0755)
	}
	c.LogDir = homeDir
	return c.LogDir
}
