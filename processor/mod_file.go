package processor

import (
	"fmt"
	"fossinator/config"
	"fossinator/fs"
	"golang.org/x/mod/modfile"
	"os"
)

// side effects - new requirements are added to the end of the last "require" block
func UpdateGoMod(dir string) error {
	fmt.Printf("----- Update go.mod [START] -----\n")
	defer fmt.Printf("----- Update go.mod [END] -----\n\n")

	filename, err := fs.FindGoModFile(dir)
	if err != nil {
		return err
	}

	src, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	mf, err := modfile.Parse("go.mod", src, nil)
	if err != nil {
		return err
	}

	update := processModFile(mf)

	if update {
		fmt.Println("Updated: go.mod")
		newContent, err := mf.Format()
		if err != nil {
			return err
		}

		return os.WriteFile(filename, newContent, 0644)
	}

	return nil
}

//--------------------------------------------------------------------------------

func processModFile(file *modfile.File) bool {
	update := replaceGoVersion(file)
	update = replaceToolchain(file) || update

	for _, r := range file.Require {
		if ok := replaceDependencies(file, r); ok {
			update = true
			continue
		}
		if ok := removeDependencies(file, r); ok {
			update = true
			continue
		}
	}
	return update
}

func replaceDependencies(mf *modfile.File, r *modfile.Require) bool {
	for _, replacement := range config.CurrentConfig.Go.LibsToReplace {
		if r.Mod.Path == replacement.OldName {
			_ = mf.DropRequire(replacement.OldName)
			_ = mf.AddRequire(replacement.NewName, replacement.NewVersion)
			return true
		}
	}
	return false
}

func removeDependencies(mf *modfile.File, r *modfile.Require) bool {
	for _, lib := range config.CurrentConfig.Go.LibsToRemove {
		if r.Mod.Path == lib.Name {
			_ = mf.DropRequire(lib.Name)
			return true
		}
	}
	return false
}

func replaceGoVersion(mf *modfile.File) bool {
	goVersion := config.CurrentConfig.Go.Version
	if len(goVersion) == 0 {
		return false
	}
	if mf.Go != nil && mf.Go.Version != goVersion {
		_ = mf.AddGoStmt(goVersion)
		return true
	}
	return false
}

func replaceToolchain(mf *modfile.File) bool {
	toolchain := config.CurrentConfig.Go.Toolchain
	if len(toolchain) == 0 {
		return false
	}
	if mf.Toolchain == nil || mf.Toolchain.Name != toolchain {
		if err := mf.AddToolchainStmt(toolchain); err != nil {
			println(err.Error())
		}
		return true
	}
	return false
}
