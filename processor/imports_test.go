package processor

import (
	"bytes"
	"fossinator/config"
	"github.com/stretchr/testify/assert"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func Test_processFile_withLibsToReplace_shouldChangeFullyEqualName(t *testing.T) {
	//config
	config.CurrentConfig.Go.LibsToReplace = []config.LibToReplace{
		{
			OldName: "company1/import1/v1",
			NewName: "company2/import2/v2",
		}}

	//data
	const input = `package main

import (
	"company1/import1/v1"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	"company2/import2/v2"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, true)
}

func Test_processFile_withLibsToReplace_shouldChangeNotFullyEqualName(t *testing.T) {
	//config
	config.CurrentConfig.Go.LibsToReplace = []config.LibToReplace{
		{
			OldName: "company1/import1/v1",
			NewName: "company2/import2/v2",
		}}

	//data
	const input = `package main

import (
	"company1/import1/v1/foo"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	"company2/import2/v2/foo"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, true)
}

func Test_processFile_withImportsToReplace_shouldChangeFullyEqualName_withAlias(t *testing.T) {
	//config
	config.CurrentConfig.Go.ImportsToReplace = []config.ImportToReplace{
		{
			OldName: "company1/import1/package_foo",
			NewName: "company2/import2/package_bar",
		}}

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	alias1 "company2/import2/package_bar"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, true)
}

func Test_processFile_withImportsToReplace_shouldChangeFullyEqualName_withoutAlias(t *testing.T) {
	//config
	config.CurrentConfig.Go.ImportsToReplace = []config.ImportToReplace{
		{
			OldName: "company1/import1/package_foo",
			NewName: "company2/import2/package_bar",
		}}

	const input = `package main

import (
	"company1/import1/package_foo"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	package_foo "company2/import2/package_bar"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, true)
}

func Test_processFile_withImportsToReplace_shouldNotChangeNotFullyEqualName(t *testing.T) {
	//config
	config.CurrentConfig.Go.ImportsToReplace = []config.ImportToReplace{
		{
			OldName: "company1/import1/package_foo",
			NewName: "company2/import2/package_bar",
		}}

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo/baz"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	alias1 "company1/import1/package_foo/baz"
	"fmt"
	"testing"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, false)
}

func Test_processFile_combinedTest(t *testing.T) {
	//config
	config.CurrentConfig.Go.ImportsToReplace = []config.ImportToReplace{
		{
			OldName: "company1/import1/package_foo",
			NewName: "company2/import2/package_bar",
		},
		{
			OldName: "company3/import3/package_qwe",
			NewName: "company4/import4/package_asd",
		},
	}

	config.CurrentConfig.Go.LibsToReplace = []config.LibToReplace{
		{
			OldName: "company1/import1",
			NewName: "company5/import5",
		},
		{
			OldName: "company6/import6",
			NewName: "company7/import7",
		},
	}

	//data
	const input = `package main

import (
	alias1 "company1/import1/package_foo"
	alias2 "company1/import1/package_baz"
	"company3/import3/package_qwe"
	"company6/import6/qwerty"
)

func main() {
	fmt.Println("hello")
}
`

	const expected = `package main

import (
	alias1 "company2/import2/package_bar"
	package_qwe "company4/import4/package_asd"
	alias2 "company5/import5/package_baz"
	"company7/import7/qwerty"
)

func main() {
	fmt.Println("hello")
}
`

	//test
	generalProcessFileTest(t, input, expected, true)
}

//--------------------------------------------------------------------------

func generalProcessFileTest(t *testing.T, input, expected string, shouldBeUpdated bool) {
	//test dto
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", input, parser.ParseComments)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	//test
	updated := processFile(file)

	//result preparation
	var buf bytes.Buffer
	err = format.Node(&buf, fset, file)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}
	actual := buf.String()

	//assertions
	assert.Equal(t, shouldBeUpdated, updated)
	assert.Equal(t, expected, actual)
}
