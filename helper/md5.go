// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
