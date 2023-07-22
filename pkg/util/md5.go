package util

import (
	"crypto/md5"
	"encoding/hex"
)

//传入一个字符串，就可以得到该字符串的 MD5 加密结果，以十六进制字符串的形式返回
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}