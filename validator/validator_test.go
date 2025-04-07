package validator

import (
	"fossinator/config"
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"testing"
)

func Test_validateDependenciesInternal_doNotContainsProhibited(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
	}()

	//data
	const input = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1
	github.com/spf13/cobra v1.9.1
)
`

	//test
	result := generalValidateDependenciesTest(t, input)

	assert.Equal(t, 0, len(result))
}

func Test_validateDependenciesInternal_containsProhibited(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
	}()

	//data
	const input = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1
	github.com/spf13/cobra v1.9.1
	foo.com/lib1/package1/v3 v3.0.0
	foo.com/lib2/package2/v4 v4.0.0
)
`

	//test
	result := generalValidateDependenciesTest(t, input)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "go.mod contains not permitted dependency: foo.com/lib1/package1/v3", result[0])
	assert.Equal(t, "go.mod contains not permitted dependency: foo.com/lib2/package2/v4", result[1])
}

func Test_validateDependenciesInternal_containsProhibitedButWhitelisted(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	config.CurrentConfig.Go.Validation.LibsWhiteList = []string{
		"foo.com/lib",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
		config.CurrentConfig.Go.Validation.LibsWhiteList = nil
	}()

	//data
	const input = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1
	github.com/spf13/cobra v1.9.1
	foo.com/lib/package/v3 v3.0.0
)
`

	//test
	result := generalValidateDependenciesTest(t, input)

	assert.Equal(t, 0, len(result))
}

//-----------------------------------------------------------------------------

func Test_validateImportsInternal_doNotContainsProhibited(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
	}()

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	result := generalValidateImportsTest(t, input, "filename")

	assert.Equal(t, 0, len(result))
}

func Test_validateImportsInternal_containsProhibited(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
	}()

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo"
	"foo.com/lib1/package1/v3"
	"foo.com/lib2/package2/v3"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	result := generalValidateImportsTest(t, input, "filename")

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "File filename contains not permitted import: foo.com/lib1/package1/v3", result[0])
	assert.Equal(t, "File filename contains not permitted import: foo.com/lib2/package2/v3", result[1])
}

func Test_validateImportsInternal_containsProhibitedButWhitelisted(t *testing.T) {
	//config
	config.CurrentConfig.Go.Validation.ProhibitedWords = []string{
		"foo.com",
	}
	config.CurrentConfig.Go.Validation.LibsWhiteList = []string{
		"foo.com/lib",
	}
	defer func() {
		config.CurrentConfig.Go.Validation.ProhibitedWords = nil
		config.CurrentConfig.Go.Validation.LibsWhiteList = nil
	}()

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo"
	"foo.com/lib1/package1/v3"
	"foo.com/lib2/package2/v3"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	result := generalValidateImportsTest(t, input, "filename")

	assert.Equal(t, 0, len(result))
}

//-------------------------------------------------------------------------------------

func generalValidateDependenciesTest(t *testing.T, input string) []string {
	//test dto
	file, err := modfile.Parse("go.mod", []byte(input), nil)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	//test
	return validateDependenciesInternal(file)
}

func generalValidateImportsTest(t *testing.T, input, path string) []string {
	//test dto
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", input, parser.ParseComments)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	//test
	return validateImportsInternal(path, file)
}
