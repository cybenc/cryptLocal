package utils

import (
	"bytes"
	"encoding/json"
)

// 将结构体转换为JSON字符串的函数
func StructToJSON(v interface{}) (string, error) {
	// 将结构体编码为JSON格式
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// 将结构体转换为格式化的JSON字符串的函数
func StructToPrettyJSON(v interface{}) (string, error) {
	// 创建一个Buffer来存储JSON数据
	var buf bytes.Buffer

	// 创建一个Encoder，并设置SetEscapeHTML为false
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")

	// 将结构体编码为格式化的JSON字符串
	if err := encoder.Encode(v); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// 将JSON字符串转换为map的函数
func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
