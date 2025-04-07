package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type LibToReplace struct {
	OldName    string `yaml:"old-name"`
	NewName    string `yaml:"new-name"`
	NewVersion string `yaml:"new-version"`
}

type ImportToReplace struct {
	OldName string `yaml:"old-name"`
	NewName string `yaml:"new-name"`
}

type LibToRemove struct {
	Name string `yaml:"name"`
}

type Config struct {
	Go struct {
		Version          string            `yaml:"version"`
		Toolchain        string            `yaml:"toolchain"`
		LibsToReplace    []LibToReplace    `yaml:"libs-to-replace"`
		ImportsToReplace []ImportToReplace `yaml:"imports-to-replace"`
		LibsToRemove     []LibToRemove     `yaml:"libs-to-remove"`
		ServiceLoading   struct {
			Imports      []string `yaml:"imports"`
			Instructions []string `yaml:"instructions"`
		} `yaml:"service-loading"`
	} `yaml:"go"`
}

var CurrentConfig Config

func Load(name string) {
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(data, &CurrentConfig); err != nil {
		panic(err)
	}
}
