package validator

import (
	"fmt"
	"fossinator/config"
	"fossinator/fs"
	"go/ast"
	"golang.org/x/mod/modfile"
	fs2 "io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Validate(dir string) []string {
	var result []string
	result = append(result, validateDependencies(dir)...)
	result = append(result, validateImports(dir)...)
	return result
}

func validateDependencies(dir string) []string {
	filename, err := fs.FindGoModFile(dir)
	if err != nil {
		return []string{err.Error()}
	}

	src, err := os.ReadFile(filename)
	if err != nil {
		return []string{err.Error()}
	}

	mf, err := modfile.Parse("go.mod", src, nil)
	if err != nil {
		return []string{err.Error()}
	}

	return validateDependenciesInternal(mf)
}

func validateDependenciesInternal(mf *modfile.File) []string {
	var result []string
	for _, req := range mf.Require {
		if isProhibited(req.Mod.Path) && !inWhitelistList(req.Mod.Path) {
			validationMessage := fmt.Sprintf("go.mod contains not permitted dependency: %v", req.Mod.Path)
			result = append(result, validationMessage)
		}
	}
	return result
}

func validateImports(dir string) []string {
	var result []string
	err := filepath.Walk(dir, func(path string, _ fs2.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		_, file, err := fs.ParseFile(path)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse file: %v", err)
			result = append(result, msg)
			return nil
		}

		result = append(result, validateImportsInternal(path, file)...)

		return nil
	})
	if err != nil {
		msg := fmt.Sprintf("Error while iterating through files: %v", err)
		result = append(result, msg)
	}
	return result
}

func validateImportsInternal(path string, file *ast.File) []string {
	var result []string
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		if isProhibited(importPath) && !inWhitelistList(importPath) {
			validationMessage := fmt.Sprintf("File %v contains not permitted import: %v", path, importPath)
			result = append(result, validationMessage)
		}
	}
	return result
}

func isProhibited(dep string) bool {
	for _, prohibitedWord := range config.CurrentConfig.Go.Validation.ProhibitedWords {
		if strings.Contains(dep, prohibitedWord) {
			return true
		}
	}
	return false
}

func inWhitelistList(dep string) bool {
	for _, whiteListLib := range config.CurrentConfig.Go.Validation.LibsWhiteList {
		if strings.HasPrefix(dep, whiteListLib) {
			return true
		}
	}
	return false
}
