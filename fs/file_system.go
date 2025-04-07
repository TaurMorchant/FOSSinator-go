package fs

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func ParseFile(path string) (*token.FileSet, *ast.File, error) {
	fs := token.NewFileSet()
	result, err := parser.ParseFile(fs, path, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing is failed: %w", err)
	}
	return fs, result, nil
}

func ParseSrc(src string) (*token.FileSet, *ast.File, error) {
	fs := token.NewFileSet()
	result, err := parser.ParseFile(fs, "", src, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing is failed: %w", err)
	}
	return fs, result, nil
}

func WriteFile(fileName, src string) error {
	return os.WriteFile(fileName, []byte(src), 0644)
}

func FindMainFile(dir string) (string, error) {
	var mainFile string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, path, nil, parser.PackageClauseOnly)
		if err != nil {
			return nil
		}
		if node.Name.Name != "main" {
			return nil
		}

		node, err = parser.ParseFile(fs, path, nil, parser.AllErrors)
		if err != nil {
			return nil
		}
		for _, decl := range node.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == "main" {
				mainFile = path
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return mainFile, nil
}

func FindGoModFile(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "go.mod" {
			return filepath.Join(dir, entry.Name()), nil
		}
	}
	return "", errors.New("go.mod not found")
}

func FmtAndWrite(fs *token.FileSet, path string, node *ast.File) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close file %s: %w", path, err))
		}
	}(file)

	//todo vlla: side effect - CRLF converted to LF
	if err = format.Node(file, fs, node); err != nil {
		return err
	}

	//todo vlla: side effect - Fprint replaces spaces by tabs AND CRLF converted to LF
	//if err := printer.Fprint(file, fs, node); err != nil {
	//	return err
	//}

	//todo vlla работает так же, как format.Node, но дольше
	//if err := runGoFmt(file.Name()); err != nil {
	//	return err
	//}

	fmt.Println("Updated:", path)
	return nil
}
