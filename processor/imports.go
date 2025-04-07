package processor

import (
	"fmt"
	"fossinator/config"
	"fossinator/fs"
	"go/ast"
	fs2 "io/fs"
	"path"
	"path/filepath"
	"strings"
)

func UpdateImports(dir string) error {
	fmt.Printf("----- Update imports [START] -----\n")
	defer fmt.Printf("----- Update imports [END] -----\n\n")
	return filepath.Walk(dir, func(path string, _ fs2.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		fileSet, node, err := fs.ParseFile(path)
		if err != nil {
			return err
		}

		updated := processFile(node)

		if updated {
			return fs.FmtAndWrite(fileSet, path, node)
		}
		return nil
	})
}

//-------------------------------------------------------------------------------------

func processFile(file *ast.File) (updated bool) {
	for _, imp := range file.Imports {
		updated = replaceFullPackage(imp) || updated
		updated = replacePackagePrefix(imp) || updated
	}

	return updated
}

func replacePackagePrefix(imp *ast.ImportSpec) bool {
	importPath := strings.Trim(imp.Path.Value, `"`)
	for _, replacement := range config.CurrentConfig.Go.LibsToReplace {
		if strings.HasPrefix(importPath, replacement.OldName) {
			imp.Path.Value = `"` + replacement.NewName + importPath[len(replacement.OldName):] + `"`
			return true
		}
	}
	return false
}

func replaceFullPackage(imp *ast.ImportSpec) bool {
	importPath := strings.Trim(imp.Path.Value, `"`)
	localName := imp.Name
	for _, replacement := range config.CurrentConfig.Go.ImportsToReplace {
		if importPath == replacement.OldName {
			imp.Path.Value = `"` + replacement.NewName + `"`

			if localName == nil {
				oldPackageName := getPackageName(replacement.OldName)
				newPackageName := getPackageName(replacement.NewName)

				if oldPackageName != newPackageName {
					imp.Name = &ast.Ident{Name: oldPackageName}
				}
			}

			return true
		}
	}
	return false
}

func getPackageName(s string) string {
	return path.Base(s)
}
