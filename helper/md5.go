package helper

import (
	"crypto/md5"
	"fmt"
)

// Md5 return the md5 value of string
func Md5(str string) string {
	data := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", data)
}
