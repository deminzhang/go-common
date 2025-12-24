package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// Deprecated: use Md5 instead
func Md5Sum(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// lowercase
func Md5(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// uppercase
func MD5(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%X", hash)
}
