# FOSSinator

FOSSinator is a tool for automating the replacement of libraries and imports in go repositories.

# How to use
- build repo
```
go build ./main.go
```
- prepare config.yaml in directory with main.exe
- run tool with target repo in args
```
./main.exe <path to your go project>
```

# Features
- Replace lib names + lib versions in go.mod
- Remove libs from go.mod
- Replace lib names in all imports (without change package names, case: lib renamed)
- Replace lib names in all imports (with change package names and alias, case: lib package moved to another lib)
- Exec 'go fmt' and 'go mod tidy' after all manipulations 