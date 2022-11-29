// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cli

import (
	"os"
	"path/filepath"
)

// Module struct
type Module struct {
	Name  string
	Path  string
	Stubs []Stub
}

// Create use for create app module
func (m *Module) Create(rootPath string, packageName string) error {
	modulePath, err := m.CreatePath(rootPath)
	if err != nil {
		return err
	}

	replace := map[string]string{ // the part's of stub to be repalce
		"Package": packageName,
	}
	for _, stub := range m.Stubs {
		if err := stub.Cast(modulePath, replace); err != nil {
			return err
		}
	}
	return nil
}

// CreatePath use for create path
func (m *Module) CreatePath(rootPath string) (string, error) {
	path := filepath.Join(rootPath, m.Path)
	return path, createPathIfNotExist(filepath.Join(path))
}

func createPathIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

var modules = []Module{
	{
		Name: "Main",
		Path: "/",
		Stubs: []Stub{
			{
				Name: "main.go",
				Path: "/",
				File: "main.go.stub",
			},
		},
	},
	{
		Name: "Api",
		Path: "/api",
		Stubs: []Stub{
			{
				Name: "handle.go",
				Path: "/welcome",
				File: "api.welcome.handle.go.stub",
			},
		},
	},
	{
		Name: "Model",
		Path: "/model",
		Stubs: []Stub{
			{
				Name: "user.go",
				Path: "/",
				File: "model.user.go.stub",
			},
		},
	},
	{
		Name: "Config",
		Path: "/",
		Stubs: []Stub{
			{
				Name: "config.yml",
				Path: "/",
				File: "config.yml.stub",
			},
			{
				Name: "config.yml.example",
				Path: "/",
				File: "config.yml.stub",
			},
		},
	},
	{
		Name: "Router",
		Path: "/router",
		Stubs: []Stub{
			{
				Name: "router.go",
				Path: "/",
				File: "router.go.stub",
			},
		},
	},
	{
		Name: "Misc",
		Path: "/",
		Stubs: []Stub{
			{
				Name: "Makefile",
				Path: "/",
				File: "makefile.stub",
			},
			{
				Name: "go.mod",
				Path: "/",
				File: "go.mod.stub",
			},
			{
				Name: ".air.conf",
				Path: "/",
				File: "air.conf.stub",
			},
			{
				Name: ".gitignore",
				Path: "/",
				File: "gitignore.stub",
			},
		},
	},
}
