package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cybenc/cryptLocal/models"
	"github.com/sirupsen/logrus"
)

type Alist struct {
	Url      string
	UserName string
	Password string
	OptCode  string
	Token    string
}

type AListLoginResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

type AListDriverInfoResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Common     []AListDriverInfo    `json:"common"`
		Additional []DriverAdditionInfo `json:"additional"`
		Config     DriverConfig         `json:"config"`
	} `json:"data"`
}

type AListDriverInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  string `json:"default"`
	Options  string `json:"options"`
	Required bool   `json:"required"`
	Help     string `json:"help"`
}

type DriverAdditionInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Default  string `json:"default"`
	Options  string `json:"options"`
	Required bool   `json:"required"`
	Help     string `json:"help"`
}

type DriverConfig struct {
	Name        string `json:"name"`
	LocalSort   bool   `json:"local_sort"`
	OnlyLocal   bool   `json:"only_local"`
	OnlyProxy   bool   `json:"only_proxy"`
	NoCache     bool   `json:"no_cache"`
	NoUpload    bool   `json:"no_upload"`
	NeedMs      bool   `json:"need_ms"`
	DefaultRoot string `json:"default_root"`
	Alert       string `json:"alert"`
}

type StorageListResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Content []StorageInfoContent `json:"content"`
		Total   int                  `json:"total"`
	}
}

type StorageInfoContent struct {
	Id              int    `json:"id"`
	MountPath       string `json:"mount_path"`
	Order           int    `json:"order"`
	Driver          string `json:"driver"`
	CacheExpiration int    `json:"cache_expiration"`
	Status          string `json:"status"`
	Addition        string `json:"addition"`
	Remark          string `json:"remark"`
	Modified        string `json:"modified"`
	Disabled        bool   `json:"disabled"`
	EnableSign      bool   `json:"enable_sign"`
	OrderBy         string `json:"order_by"`
	OrderDirection  string `json:"order_direction"`
	ExtractFolder   string `json:"extract_folder"`
	WebProxy        bool   `json:"web_proxy"`
	WebdavPolicy    string `json:"webdav_policy"`
	DownProxyURL    string `json:"down_proxy_url"`
}

type StorageInfo struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    StorageInfoContent `json:"data"`
}

func (a *Alist) HttpGet(url string, resp interface{}) error {
	// 发送请求 headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", a.Token)
	// 发送请求
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer r.Body.Close()
	// 读取响应体并打印
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	err = json.Unmarshal(bodyBytes, resp)
	if err != nil {
		return err
	}
	return nil
}

func (a *Alist) Login() (string, error) {
	url := a.Url
	// 判断是否以http或者https开头
	if !strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		fmt.Println("url必须以http://或者https://开头")
		return "", fmt.Errorf("url必须以http://或者https://开头")
	}
	url = url + "/api/auth/login"
	// 构建body
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, a.UserName, a.Password)
	if a.OptCode != "" {
		body = fmt.Sprintf(`{"username":"%s","password":"%s","optcode":"%s"}`, a.UserName, a.Password, a.OptCode)
	}
	// 发送请求
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer resp.Body.Close()
	// 读取响应体并打印
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}
	var loginResp AListLoginResp
	err = json.Unmarshal(bodyBytes, &loginResp)
	if err != nil {
		return "", err
	}
	if loginResp.Code != 200 {
		return "", fmt.Errorf("Error:%s", loginResp.Message)
	}
	a.Token = loginResp.Data.Token
	return loginResp.Data.Token, nil
}

func (a *Alist) GetCryptDriverInfo() (*AListDriverInfoResp, error) {
	url := a.Url + "/api/admin/driver/info?driver=Crypt"
	var driverInfoResp AListDriverInfoResp
	err := a.HttpGet(url, &driverInfoResp)
	if err != nil {
		return nil, err
	}
	if driverInfoResp.Code != 200 {
		return nil, fmt.Errorf("Error:%s", driverInfoResp.Message)
	}
	return &driverInfoResp, nil
}

func (a *Alist) GetStorageList() (*StorageListResp, error) {
	url := a.Url + "/api/admin/storage/list"
	var storageListResp StorageListResp
	err := a.HttpGet(url, &storageListResp)
	if err != nil {
		return nil, err
	}
	if storageListResp.Code != 200 {
		return nil, fmt.Errorf("Error:%s", storageListResp.Message)
	}
	return &storageListResp, nil
}

func (a *Alist) GetCryptSorageId() (int, error) {
	// 获取存储列表
	storageListResp, err := a.GetStorageList()
	if err != nil {
		return 0, err
	}
	// 遍历存储列表，找到crypt存储的id
	for _, storage := range storageListResp.Data.Content {
		if storage.Driver == "Crypt" {
			return storage.Id, nil
		}
	}
	return 0, nil
}

func (a *Alist) GetStorageIdByName(name string) (int, error) {
	// 获取存储列表
	storageListResp, err := a.GetStorageList()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// 遍历存储列表，找到指定名称的存储的id
	for _, storage := range storageListResp.Data.Content {
		if storage.Driver == name {
			return storage.Id, nil
		}
	}
	return 0, nil
}

func (a *Alist) GetStorageInfo(storageId int) (*StorageInfo, error) {
	url := a.Url + fmt.Sprintf("/api/admin/storage/get?id=%d", storageId)
	var storageInfoResp StorageInfo
	err := a.HttpGet(url, &storageInfoResp)
	if err != nil {
		return nil, err
	}
	if storageInfoResp.Code != 200 {
		return nil, fmt.Errorf("Error:%s", storageInfoResp.Message)
	}
	return &storageInfoResp, nil
}

func (a *Alist) GetStorageInfoByName(name string) (*StorageInfo, error) {
	// 获取存储id
	storageId, err := a.GetStorageIdByName(name)
	if err != nil || storageId <= 0 {
		return nil, err
	}
	// 获取存储信息
	storageInfo, err := a.GetStorageInfo(storageId)
	if err != nil || storageInfo == nil {
		return nil, err
	}
	return storageInfo, nil
}

func (a *Alist) GenerateConfigByApi() (*models.EncrtptConfig, error) {
	info, err := a.GetStorageInfoByName("Crypt")
	if err != nil || info == nil {
		return nil, err
	}
	logrus.Info("获取crypt驱动信息成功")
	var config models.EncrtptConfig
	addition := info.Data.Addition
	// 转换为map
	var additionMap map[string]interface{}
	err = json.Unmarshal([]byte(addition), &additionMap)
	if err != nil {
		return nil, err
	}
	// 解析配置
	config.FileNameEncryption = additionMap["filename_encryption"].(string)
	config.DirectoryNameEncryption = additionMap["directory_name_encryption"].(string)
	config.Password = additionMap["password"].(string)
	config.Password2 = additionMap["salt"].(string)
	config.Suffix = additionMap["encrypted_suffix"].(string)
	config.FileNameEncoding = additionMap["filename_encoding"].(string)
	config.PassBadBlocks = ""
	return &config, nil
}
