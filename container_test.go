package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type as struct {
	id int
}

type bs struct {
	A *as
}

func TestContainer_Make(t *testing.T) {
	con := New()
	con.Register(func() *as {
		return &as{1}
	})
	con.Register(func(a *as) *bs {
		return &bs{A: a}
	})

	assert.NotPanics(t, func() {
		con.Make(func(a *as, b *bs) {
			assert.Equal(t, a.id, 1)
			assert.Equal(t, a.id, b.A.id)
		})
	})
}

func TestContainer_Register(t *testing.T) {
	con := New()
	assert.NotPanics(t, func() {
		con.Register(func() *as {
			return &as{1}
		})
	})
}

func TestContainer_String(t *testing.T) {
	assert.NotPanics(t, func() {
		con := New()
		con.Register(func() *as {
			return &as{1}
		})
		con.Register(func(a *as) *bs {
			return &bs{A: a}
		})
		fmt.Println(con.String())
	})
}

type I interface {
	Get() string
}

type i struct{}

func (*i) Get() string {
	return "i"
}

type testHandler struct {
	i I
}

func TestContainer_Interface(t *testing.T) {
	con := New()
	assert.NotPanics(t, func() {
		con.Register(func() I {
			return &i{}
		})
		con.Make(func(i I) {
			h := &testHandler{
				i,
			}
			assert.Equal(t, "i", h.i.Get())
		})

		var i2 I
		con.Make(&i2)
		assert.Equal(t, "i", i2.Get())
	})
}

func do(di *Container, receiver interface{}) {
	di.Make(receiver)
}

func TestContainer_Make_Interface(t *testing.T) {
	di := New()
	assert.NotPanics(t, func() {
		di.Register(func() I {
			return &i{}
		})
		do(di, func(i I) {
			h := &testHandler{
				i,
			}
			assert.Equal(t, "i", h.i.Get())
		})
	})
}
