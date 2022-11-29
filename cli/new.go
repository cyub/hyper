// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

var newCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "create an application",
	Action:  newApp,
}

func newApp(c *cli.Context) error {
	// check cli args number
	if c.NArg() != 1 {
		fmt.Println("Usage: hyper new your_project_name")
		os.Exit(1)
	}
	name := strings.ToLower(c.Args().Get(0))
	path := filepath.Join(app.pwd, name)
	_, err := os.Stat(path)
	if err == nil {
		color.Red.Println("app path had exists")
		os.Exit(1)
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		color.Red.Printf("create app path failure %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("start create %s's modules\n", name)
	for _, module := range modules {
		color.Bold.Printf("%s%s\n", strings.Repeat(" ", 2), module.Name)
		if err := module.Create(filepath.Join(app.pwd, name), name); err != nil {
			return err
		}
	}
	color.Style{color.Green, color.OpBold}.Printf("%s create success, enjoy it!\n", name)
	return nil
}
