// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/cyub/hyper/cli"
)

func main() {
	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
