package models

type EncrtptConfig struct {
	// 密码
	Password string `json:"password"`
	// 盐
	Password2 string `json:"password2"`
	// off standard obfuscate
	FileNameEncryption string `json:"filename_encryption"`
	// 密码
	DirectoryNameEncryption string `json:"directory_name_encryption"`
	// 文件名编码方式
	FileNameEncoding string `json:"filename_encoding"`
	// 文件后缀
	Suffix string `json:"suffix"`
	// 加密算法
	PassBadBlocks string `json:"pass_bad_blocks"`
}
