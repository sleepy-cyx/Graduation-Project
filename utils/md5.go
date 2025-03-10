package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5BySalt(str string, salt string) string {
	strByte := []byte(str)
	saltByte := []byte(salt)
	h := md5.New()
	h.Write(saltByte)
	h.Write(strByte)
	return hex.EncodeToString(h.Sum(nil))
}
