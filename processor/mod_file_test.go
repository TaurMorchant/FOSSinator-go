package processor

import (
	"fossinator/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
	"testing"
)

func Test_processModFile_changeLib(t *testing.T) {
	//config
	config.CurrentConfig.Go.LibsToReplace = []config.LibToReplace{
		{
			OldName:    "company1.com/import1/v2",
			NewName:    "company2.com/import2/v3",
			NewVersion: "v3.0.0",
		}}
	defer func() {
		config.CurrentConfig.Go.LibsToReplace = nil
	}()

	//data
	const input = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1
	company1.com/import1/v2 v2.0.0
)
`

	const expected = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1

	company2.com/import2/v3 v3.0.0
)
`

	//test
	generalProcessModFileTest(t, input, expected, true)
}

func Test_processModFile_removeLib(t *testing.T) {
	//config
	config.CurrentConfig.Go.LibsToRemove = []config.LibToRemove{
		{
			Name: "company1.com/import1/v2",
		}}
	defer func() {
		config.CurrentConfig.Go.LibsToRemove = nil
	}()

	//data
	const input = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1
	company1.com/import1/v2 v2.0.0
)
`

	const expected = `module fossinator

go 1.23.0

require (
	gopkg.in/yaml.v3 v3.0.1

)
`

	//test
	generalProcessModFileTest(t, input, expected, true)
}

func Test_processModFile_changeGoVersion(t *testing.T) {
	//config
	config.CurrentConfig.Go.Version = "1.23.0"
	config.CurrentConfig.Go.Toolchain = "go1.23.4"
	defer func() {
		config.CurrentConfig.Go.Version = ""
		config.CurrentConfig.Go.Toolchain = ""
	}()

	//data
	const input = `module fossinator

go 1.22.0

toolchain go1.22.0

require (
	gopkg.in/yaml.v3 v3.0.1
)
`

	const expected = `module fossinator

go 1.23.0

toolchain go1.23.4

require (
	gopkg.in/yaml.v3 v3.0.1
)
`

	//test
	generalProcessModFileTest(t, input, expected, true)
}

func Test_processModFile_changeGoVersion_toolchainIsAbsent(t *testing.T) {
	//config
	config.CurrentConfig.Go.Version = "1.23.0"
	config.CurrentConfig.Go.Toolchain = "go1.23.4"
	defer func() {
		config.CurrentConfig.Go.Version = ""
		config.CurrentConfig.Go.Toolchain = ""
	}()

	//data
	const input = `module fossinator

go 1.22.0

require (
	gopkg.in/yaml.v3 v3.0.1
)
`

	const expected = `module fossinator

go 1.23.0

toolchain go1.23.4

require (
	gopkg.in/yaml.v3 v3.0.1
)
`

	//test
	generalProcessModFileTest(t, input, expected, true)
}

//--------------------------------------------------------------------------

func generalProcessModFileTest(t *testing.T, input, expected string, shouldBeUpdated bool) {
	//test dto
	file, err := modfile.Parse("go.mod", []byte(input), nil)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	//test
	updated := processModFile(file)

	//result preparation
	actualBytes, err := file.Format()
	actual := string(actualBytes)

	//assertions
	assert.Equal(t, shouldBeUpdated, updated)
	assert.Equal(t, expected, actual)
}
