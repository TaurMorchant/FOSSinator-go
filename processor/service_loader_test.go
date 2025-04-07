package processor

import (
	"fossinator/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_insertInitPart_configIsEmpty(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

func init() {
	statement1()
	statement2()
}
`

	expected := `
package transformer

import (
	"testing"
)

func init() {
	statement1()
	statement2()
}
`
	generalTestInsertInitPart(t, src, expected, nil)
}

func Test_insertInitPart_initFuncExists(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

func init() {
	statement1()
	statement2()
}
`

	expected := `
package transformer

import (
	"testing"
)

func init() {
	statement1()
	statement2()
	` + PreComment + `
	one
	two
	three
}
`
	generalTestInsertInitPart(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertInitPart_initFuncNotExists_funcIsFirst(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

const pi = 3.14

func first() {
}

func second() {
}
`

	expected := `
package transformer

import (
	"testing"
)

const pi = 3.14

func init() {
	` + PreComment + `
	one
	two
	three
}

func first() {
}

func second() {
}
`
	generalTestInsertInitPart(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertInitPart_initFuncNotExists_methodIsFirst(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

type MyType struct {}

func (m *MyType) first() {
}

func second() {
}
`

	expected := `
package transformer

import (
	"testing"
)

type MyType struct {}

func init() {
	` + PreComment + `
	one
	two
	three
}

func (m *MyType) first() {
}

func second() {
}
`
	generalTestInsertInitPart(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertInitPart_initFuncNotExists_firstFuncWithComment(t *testing.T) {
	config.CurrentConfig.Go.ServiceLoading.Instructions = []string{"one", "two", "three"}
	src := `
package transformer

import (
	"testing"
)

//some comment
//some additional comment
func first() {
}

func second() {
}
`

	expected := `
package transformer

import (
	"testing"
)

func init() {
	` + PreComment + `
	one
	two
	three
}

//some comment
//some additional comment
func first() {
}

func second() {
}
`
	generalTestInsertInitPart(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertInitPart_initFuncNotExists_noFuncExists(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

const pi = 3.14

`

	expected := `
package transformer

import (
	"testing"
)

const pi = 3.14

func init() {
	` + PreComment + `
	one
	two
	three
}

`
	generalTestInsertInitPart(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertImports_configIsEmpty(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

func init() {
}
`

	expected := `
package transformer

import (
	"testing"
)

func init() {
}
`
	generalTestInsertImports(t, src, expected, nil)
}

func Test_insertImports_importBlockNotExists_funcIsFirst(t *testing.T) {
	config.CurrentConfig.Go.ServiceLoading.Imports = []string{"one", "two", "three"}
	src := `
package transformer

func first() {
}

func second() {
}
`

	expected := `
package transformer

import (
	"one"
	"two"
	"three"
)

func first() {
}

func second() {
}
`
	generalTestInsertImports(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertImports_importBlockNotExists_methodIsFirst(t *testing.T) {
	config.CurrentConfig.Go.ServiceLoading.Imports = []string{"one", "two", "three"}
	src := `
package transformer

type MyType struct {}

func (m *MyType) method() {
}
`

	expected := `
package transformer

type MyType struct {}

import (
	"one"
	"two"
	"three"
)

func (m *MyType) method() {
}
`
	generalTestInsertImports(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertImports_importBlockNotExists_firstFuncWithComment(t *testing.T) {
	config.CurrentConfig.Go.ServiceLoading.Imports = []string{"one", "two", "three"}
	src := `
package transformer

//some comment
//some additional comment
func first() {
}
`

	expected := `
package transformer

import (
	"one"
	"two"
	"three"
)

//some comment
//some additional comment
func first() {
}
`
	generalTestInsertImports(t, src, expected, []string{"one", "two", "three"})
}

func Test_insertImports_importBlockExists(t *testing.T) {
	src := `
package transformer

import (
	"testing"
)

func first() {
}
`

	expected := `
package transformer

import (
	"testing"
	"one"
	"two"
	"three"
)

func first() {
}
`
	generalTestInsertImports(t, src, expected, []string{"one", "two", "three"})
}

//----------------------------------------------------------------

func generalTestInsertImports(t *testing.T, src, expected string, list []string) {
	actual, err := insertImports(src, list)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func generalTestInsertInitPart(t *testing.T, src, expected string, list []string) {
	actual, err := insertInitPart(src, list)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
