package utils

import (
	"strings"

	"github.com/cybenc/cryptLocal/common"
	"github.com/cybenc/cryptLocal/models"
	rcCrypt "github.com/rclone/rclone/backend/crypt"
	"github.com/rclone/rclone/fs/config/configmap"
)

var CipherInstance *Cipher

type Cipher struct {
	RcCipher *rcCrypt.Cipher
}

func (c *Cipher) NewCiper(config *models.EncrtptConfig) (*rcCrypt.Cipher, error) {
	p, _ := strings.CutPrefix(config.Password, common.ObfuscatedPrefix)
	p2, _ := strings.CutPrefix(config.Password2, common.ObfuscatedPrefix)

	cmap := configmap.Simple{
		"password":                  p,
		"password2":                 p2,
		"filename_encryption":       config.FileNameEncryption,
		"directory_name_encryption": config.DirectoryNameEncryption,
		"filename_encoding":         config.FileNameEncoding,
		"suffix":                    config.Suffix,
		"pass_bad_blocks":           "",
	}
	c.RcCipher, _ = rcCrypt.NewCipher(cmap)
	return c.RcCipher, nil
}
