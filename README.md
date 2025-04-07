# FOSSinator

FOSSinator is a tool for automating the replacement of libraries and imports in go repositories.

# How to use
- build repo
```
go build ./fossinator.go
```
- prepare config.yaml in directory with fossinator.exe
- run tool with target repo in args (if `-dir` arg is empty - run in current folder)
```
./fossinator.exe -dir <path to your go project>
```
- optional flags:
  - `-fmt` - perform code formatting
  - `-tidy` - perform 'go mod tidy'

# Features
- Replace lib names + lib versions in go.mod
- Remove libs from go.mod
- Replace lib names in all imports (without change package names, case: lib renamed)
- Replace lib names in all imports (with change package names and alias, case: lib package moved to another lib)
- Exec 'go fmt' and 'go mod tidy' after all manipulations 

# Config structure
Config has name config.yaml and should be placed in /config directory. It embedded into exe file during build.
Config fields:
- `go.version` - defines the version of golang in the mod file to replace
- `go.toolchain` - defines the toolchain version in the mod file to replace
- `go.libs-to-replace` - defines list of libs to replace. FOSSinator will replace them both in go.mod file and in imports. Suitable for the case when the lib has not changed structurally, but its version or name has changed.
  - `old-name` - old name of lib (without package name)
  - `new-name` - name to replace with
  - `new-version` - version to replace with
- `go.libs-to-remove` - defines list of libs to remove from go.mod
  - `name` - name of lib to remove
- `go.imports-to-replace` - defines list of packages to replace in import statements. Suitable for the case when the package has moved from one lib to another
  - `old-name` - old name of import (with package name)
  - `new-name` - name to replace with
- `go.service-loading` - defines general configuration of service loading mechanism. FOSSinator will find file with main function and insert imports and init() method with SL configuration in it
  - `imports` - list of imports to insert in file with main function
  - `instructions` - list of go instructions to insert in init() method in file with main function