/**
* @program: exception
*
* @description:
*
* @author: lemo
*
* @create: 2021-05-22 00:49
**/
package exception

import (
	"testing"

	"github.com/lemoyxk/caller"
	"github.com/stretchr/testify/assert"
)

type T struct {
	Name string
}

func (t *T) A() string {
	return t.Name
}

type B interface {
	A() string
}

func TestTry(t *testing.T) {
	file, line := caller.Deep(1)
	assert.Equal(t, file, "ex_test.go", file)
	assert.Equal(t, line, 32, line)

	var err = New("hello error")
	assert.Equal(t, err.Error(), "hello error", err.Error())
	assert.Equal(t, err.Line(), 36, err.Line())
	assert.Equal(t, err.File(), "ex_test.go", err.File())

	err = Try(func() {
		panic(1)
		println("never work")
	}).Catch(func(err Error) {
		assert.Equal(t, err.Error(), "1", err.Error())
		assert.Equal(t, err.Line(), 42, err.Line())
		assert.Equal(t, err.File(), "ex_test.go", err.File())
	}).Finally(func(err Error) {
		assert.Equal(t, err.Error(), "1", err.Error())
		assert.Equal(t, err.Line(), 42, err.Line())
		assert.Equal(t, err.File(), "ex_test.go", err.File())
	})

	assert.Equal(t, err.Error(), "1", err.Error())
}

func BenchmarkIsNil(b *testing.B) {
	var c []int
	for i := 0; i < b.N; i++ {
		IsNil(c)
	}
}
