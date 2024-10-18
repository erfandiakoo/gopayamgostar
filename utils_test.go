package gopayamgostar_test

import (
	"testing"

	"github.com/erfandiakoo/gopayamgostar/v1"
	"github.com/stretchr/testify/assert"
)

func TestStringP(t *testing.T) {
	p := gopayamgostar.StringP("test value")
	assert.Equal(t, "test value", *p)
}
func TestPString(t *testing.T) {
	p := "test value"
	v := gopayamgostar.PString(&p)
	assert.Equal(t, p, v)
}

func TestPStringNil(t *testing.T) {
	v := gopayamgostar.PString(nil)
	assert.Equal(t, "", v)
}

func TestBoolP(t *testing.T) {
	p1 := gopayamgostar.BoolP(false)
	assert.False(t, *p1)
	p2 := gopayamgostar.BoolP(false)
	assert.False(t, *p1)
	assert.False(t, p1 == p2)

	p := gopayamgostar.BoolP(true)
	assert.True(t, *p)
}

func TestPBool(t *testing.T) {
	p := true
	v := gopayamgostar.PBool(&p)
	assert.True(t, v)

	p = false
	v = gopayamgostar.PBool(&p)
	assert.False(t, v)

	v = gopayamgostar.PBool(nil)
	assert.False(t, v)
}

func TestIntP(t *testing.T) {
	p := gopayamgostar.IntP(42)
	assert.Equal(t, 42, *p)
}

func TestInt32P(t *testing.T) {
	v := int32(42)
	p := gopayamgostar.Int32P(v)
	assert.Equal(t, v, *p)
}

func TestInt64P(t *testing.T) {
	v := int64(42)
	p := gopayamgostar.Int64P(v)
	assert.Equal(t, v, *p)
}

func TestPInt(t *testing.T) {
	p := 42
	v := gopayamgostar.PInt(&p)
	assert.Equal(t, p, v)
	assert.IsType(t, p, v)

	v = gopayamgostar.PInt(nil)
	assert.Equal(t, int(0), v)
	assert.IsType(t, int(0), v)
}

func TestPInt32(t *testing.T) {
	var p int32 = 42
	v := gopayamgostar.PInt32(&p)
	assert.Equal(t, p, v)
	assert.IsType(t, p, v)

	v = gopayamgostar.PInt32(nil)
	assert.Equal(t, int32(0), v)
	assert.IsType(t, int32(0), v)
}

func TestPInt64(t *testing.T) {
	var p int64 = 42
	v := gopayamgostar.PInt64(&p)
	assert.Equal(t, p, v)
	assert.IsType(t, p, v)

	v = gopayamgostar.PInt64(nil)
	assert.Equal(t, int64(0), v)
	assert.IsType(t, int64(0), v)
}

func TestFloat32P(t *testing.T) {
	v := float32(42.42)
	p := gopayamgostar.Float32P(v)
	assert.Equal(t, v, *p)
}
func TestFloat64P(t *testing.T) {
	v := 42.42
	p := gopayamgostar.Float64P(v)
	assert.Equal(t, v, *p)
}

func TestPFloat32(t *testing.T) {
	var p float32 = 42.42
	v := gopayamgostar.PFloat32(&p)
	assert.Equal(t, p, v)
	assert.IsType(t, p, v)

	v = gopayamgostar.PFloat32(nil)
	assert.Equal(t, float32(0), v)
	assert.IsType(t, float32(0), v)
}
func TestPFloat64(t *testing.T) {
	p := 42.42
	v := gopayamgostar.PFloat64(&p)
	assert.Equal(t, p, v)
	assert.IsType(t, p, v)

	v = gopayamgostar.PFloat64(nil)
	assert.Equal(t, float64(0), v)
	assert.IsType(t, float64(0), v)
}
func TestNilOrEmptyArray(t *testing.T) {
	a := gopayamgostar.NilOrEmptyArray(&[]string{"c", "d"})
	b := gopayamgostar.NilOrEmptyArray(&[]string{"", "b"})
	c := gopayamgostar.NilOrEmptyArray(&[]string{})
	assert.False(t, a)
	assert.True(t, b)
	assert.True(t, c)

}
