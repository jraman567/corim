package comid

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

type IInstanceValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// Instance stores an instance identity. The supported formats are UUID and UEID.
type Instance struct {
	Value IInstanceValue
}

// NewInstanceUEID instantiates a new instance with the supplied UEID identity
func NewInstanceUEID(val any) (*Instance, error) {
	if val == nil {
		return &Instance{&TaggedUEID{}}, nil
	}

	ret, err := NewTaggedUEID(val)
	if err != nil {
		return nil, err
	}
	return &Instance{ret}, nil
}

func MustNewInstanceUEID(val any) *Instance {
	ret, err := NewInstanceUEID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewInstanceUUID instantiates a new instance with the supplied UUID identity
func NewInstanceUUID(val any) (*Instance, error) {
	if val == nil {
		return &Instance{&TaggedUUID{}}, nil
	}

	ret, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &Instance{ret}, nil
}

func MustNewInstanceUUID(val any) *Instance {
	ret, err := NewInstanceUUID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// Valid checks for the validity of given instance
func (o Instance) Valid() error {
	if o.String() == "" {
		return fmt.Errorf("invalid instance id")
	}
	return nil
}

// String returns a printable string of the Instance value.  UUIDs use the
// canonical 8-4-4-4-12 format, UEIDs are hex encoded.
func (o Instance) String() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.String()
}

func (o Instance) Type() string {
	return o.Value.Type()
}

func (o Instance) Bytes() []byte {
	return o.Value.Bytes()
}

// MarshalCBOR serializes the target instance to CBOR
func (o Instance) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

func (o *Instance) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target Instance
// The instance object must have the following shape:
//
//	{
//	  "type": "<INSTANCE_TYPE>",
//	  "value": "<INSTANCE_STRING_VALUE>"
//	}
//
// where <INSTANCE_TYPE> must be one of the known IInstanceValue implementation
// type names (available in the base implementation: "ueid" and "uuid"), and
// <INSTANCE_STRING_VALUE> is the instance value encoded as a string. The exact
// encoding is <INSTANCE_TYPE> depenent. For the base implmentation types it is
//
//	ueid: base64-encoded bytes, e.g. "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
//	uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"
func (o *Instance) UnmarshalJSON(data []byte) error {
	var value encoding.TypeAndValue

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value.Type == "" {
		return errors.New("key type not set")
	}

	factory, ok := instanceValueRegister[value.Type]
	if !ok {
		return fmt.Errorf("unknown class id type: %q", value.Type)
	}

	v, err := factory(value.Value)
	if err != nil {
		return err
	}

	o.Value = v.Value

	return o.Valid()
}

func (o Instance) MarshalJSON() ([]byte, error) {
	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: o.Value.String(),
	}

	return json.Marshal(value)
}

type IInstanceFactory func(any) (*Instance, error)

var instanceValueRegister = map[string]IInstanceFactory{
	UEIDType: NewInstanceUEID,
	UUIDType: NewInstanceUUID,
}

func RegisterInstanceType(tag uint64, factory IInstanceFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := instanceValueRegister[typ]; exists {
		return fmt.Errorf("class ID type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	instanceValueRegister[typ] = factory

	return nil
}
