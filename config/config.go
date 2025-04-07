package config

import (
	_ "embed"
	"gopkg.in/yaml.v3"
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
		Validation struct {
			LibsWhiteList   []string `yaml:"libs-whitelist"`
			ProhibitedWords []string `yaml:"prohibited-words"`
		} `yaml:"validation"`
	} `yaml:"go"`
}

var CurrentConfig Config

//go:embed config.yaml
var configSrc []byte

func Load() error {
	if err := yaml.Unmarshal(configSrc, &CurrentConfig); err != nil {
		return err
	}

	return nil
}
