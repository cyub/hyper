// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cli

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

// Stub struct
type Stub struct {
	Name string
	Path string
	File string
}

// Cast use for create code file by stub
func (s *Stub) Cast(modulePath string, replace interface{}) error {
	content, err := s.GetStubContent()
	if err != nil {
		return err
	}
	// Parse stub
	t, err := template.New(s.File).Parse(content)
	if err != nil {
		return err
	}

	buffer := bytes.NewBufferString("")
	if err := t.Execute(buffer, replace); err != nil {
		return err
	}

	var destFile string
	if len(s.Path) == 0 || s.Path == "/" {
		destFile = filepath.Join(modulePath, s.Name)
	} else {
		if err := createPathIfNotExist(filepath.Join(modulePath, s.Path)); err != nil {
			return err
		}
		destFile = filepath.Join(modulePath, s.Path, s.Name)
	}

	if err := ioutil.WriteFile(destFile, buffer.Bytes(), os.ModePerm); err != nil {
		return err
	}
	fmt.Printf("%s%s  %s\n", strings.Repeat(" ", 4), s.Name, color.Style{color.Green, color.OpBold}.Render("✓"))
	return nil
}

// GetStubContent use for Get stub file content
func (s *Stub) GetStubContent() (string, error) {
	return app.box.FindString(s.File)
}
