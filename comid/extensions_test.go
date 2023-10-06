package comid

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var FlagTestFlag = Flag(-1)

type TestExtension struct {
	TestFlag *bool
}

func (o *TestExtension) ValidComid(v *Comid) error {
	return errors.New("invalid")
}

func (o *TestExtension) ValidTriples(v *Triples) error {
	return errors.New("invalid")
}

func (o *TestExtension) ValidMval(v *Mval) error {
	return errors.New("invalid")
}

func (o *TestExtension) ValidFlagsMap(v *FlagsMap) error {
	return errors.New("invalid")
}

func (o *TestExtension) SetTrue(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = &True
	}
}
func (o *TestExtension) SetFalse(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = &False
	}
}

func (o *TestExtension) Clear(flag Flag) {
	if flag == FlagTestFlag {
		o.TestFlag = nil
	}
}

func (o *TestExtension) Get(flag Flag) *bool {
	if flag == FlagTestFlag {
		return o.TestFlag
	}

	return nil
}

func (o *TestExtension) AnySet() bool {
	return o.TestFlag != nil
}

func Test_Extensions(t *testing.T) {
	exts := Extensions{}
	exts.Register(&TestExtension{})

	err := exts.ValidComid(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.ValidTriples(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.ValidMval(nil)
	assert.EqualError(t, err, "invalid")

	err = exts.ValidFlagsMap(nil)
	assert.EqualError(t, err, "invalid")

	assert.False(t, exts.AnySet())

	exts.SetTrue(FlagTestFlag)
	assert.True(t, exts.AnySet())
	assert.True(t, *exts.Get(FlagTestFlag))

	exts.SetFalse(FlagTestFlag)
	assert.False(t, *exts.Get(FlagTestFlag))

	exts.Clear(FlagTestFlag)
	assert.Nil(t, exts.Get(FlagTestFlag))
	assert.False(t, exts.AnySet())
}
