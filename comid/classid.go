// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/veraison/corim/encoding"
	"github.com/veraison/corim/extensions"
)

const (
	IntType    = "int"
	ImplIDType = "psa.impl-id"
)

type IClassIDValue interface {
	extensions.ITypeChoiceValue

	Bytes() []byte
}

// ClassID represents a $class-id-type-choice, which can be one of TaggedUUID,
// TaggedOID, or TaggedImplID (PSA-specific extension)
type ClassID struct {
	Value IClassIDValue
}

func (o ClassID) Valid() error {
	if o.Value == nil {
		return errors.New("nil value")
	}

	return o.Value.Valid()
}

// Type returns the type of the ClassID
func (o ClassID) Type() string {
	if o.Value == nil {
		return ""
	}

	return o.Value.Type()
}

// Bytes returns a []byte containing the raw bytes of the class id value
func (o ClassID) Bytes() []byte {
	if o.Value == nil {
		return []byte{}
	}
	return o.Value.Bytes()
}

// IsSet returns true iff the underlying class id value has been set (is not nil)
func (o ClassID) IsSet() bool {
	return o.Value != nil
}

// MarshalCBOR serializes the target ClassID to CBOR
func (o ClassID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.Value)
}

// UnmarshalCBOR deserializes the supplied CBOR buffer into the target ClassID.
// It is undefined behavior to try and inspect the target ClassID in case this
// method returns an error.
func (o *ClassID) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.Value)
}

// UnmarshalJSON deserializes the supplied JSON object into the target ClassID
// The class id object must have the following shape:
//
//	{
//	  "type": "<CLASS_ID_TYPE>",
//	  "value": "<CLASS_ID_STRING_VALUE>"
//	}
//
// where <CLASS_ID_TYPE> must be one of the known IClassIDValue implementation
// type names (available in the base implementation: "uuid", "oid",
// "psa.impl-id"), and <CLASS_ID_STRING_VALUE> is the class id value encoded as
// a string. The exact encoding is <CLASS_ID_TYPE> depenent. For the base
// implmentation types it is
//
//	oid: dot-seprated integers, e.g. "1.2.3.4"
//	psa.impl-id: base64-encoded bytes, e.g. "YWNtZS1pbXBsZW1lbnRhdGlvbi1pZC0wMDAwMDAwMDE="
//	uuid: standard UUID string representation, e.g. "550e8400-e29b-41d4-a716-446655440000"
func (o *ClassID) UnmarshalJSON(data []byte) error {
	var value encoding.TypeAndValue

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if value.Type == "" {
		return errors.New("class id type not set")
	}

	factory, ok := classIDValueRegister[value.Type]
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

// MarshalJSON serializes the target ClassID to JSON
func (o ClassID) MarshalJSON() ([]byte, error) {
	value := encoding.TypeAndValue{
		Type:  o.Value.Type(),
		Value: o.Value.String(),
	}

	return json.Marshal(value)
}

// String returns a printable string of the ClassID value. UUIDs use the
// canonical 8-4-4-4-12 format, PSA Implementation IDs are base64 encoded.
// OIDs are output in dotted-decimal notation.
func (o ClassID) String() string {
	return o.Value.String()
}

type ImplID [32]byte

func (o ImplID) String() string {
	return base64.StdEncoding.EncodeToString(o[:])
}

func (o ImplID) Valid() error {
	return nil
}

type TaggedImplID ImplID

func NewImplIDClassID(val any) (*ClassID, error) {
	var ret TaggedImplID

	if val == nil {
		return &ClassID{&TaggedImplID{}}, nil
	}

	switch t := val.(type) {
	case []byte:
		if nb := len(t); nb != 32 {
			return nil, fmt.Errorf("bad ImplID: got %d bytes, want 32", nb)
		}

		copy(ret[:], t)
	case string:
		v, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return nil, fmt.Errorf("bad ImplID: %w", err)
		}

		if nb := len(v); nb != 32 {
			return nil, fmt.Errorf("bad ImplID: decoded %d bytes, want 32", nb)
		}

		copy(ret[:], v)
	case TaggedImplID:
		copy(ret[:], t[:])
	case *TaggedImplID:
		copy(ret[:], (*t)[:])
	case ImplID:
		copy(ret[:], t[:])
	case *ImplID:
		copy(ret[:], (*t)[:])
	default:
		return nil, fmt.Errorf("unexpected type for ImplID: %T", t)
	}

	return &ClassID{&ret}, nil
}

func MustNewImplIDClassID(val any) *ClassID {
	ret, err := NewImplIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func (o TaggedImplID) Valid() error {
	return ImplID(o).Valid()
}

func (o TaggedImplID) String() string {
	return ImplID(o).String()
}

func (o TaggedImplID) Type() string {
	return ImplIDType
}

func (o TaggedImplID) Bytes() []byte {
	return o[:]
}

func NewOIDClassID(val any) (*ClassID, error) {
	ret, err := NewTaggedOID(val)
	if err != nil {
		return nil, err
	}

	return &ClassID{ret}, nil
}

func MustNewOIDClassID(val any) *ClassID {
	ret, err := NewOIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

func NewUUIDClassID(val any) (*ClassID, error) {
	if val == nil {
		return &ClassID{&TaggedUUID{}}, nil
	}

	ret, err := NewTaggedUUID(val)
	if err != nil {
		return nil, err
	}

	return &ClassID{ret}, nil
}

func MustNewUUIDClassID(val any) *ClassID {
	ret, err := NewUUIDClassID(val)
	if err != nil {
		panic(err)
	}

	return ret
}

type TaggedInt int

func NewIntClassID(val any) (*ClassID, error) {
	if val == nil {
		zeroVal := TaggedInt(0)
		return &ClassID{&zeroVal}, nil
	}

	var ret TaggedInt

	switch t := val.(type) {
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("bad int: %w", err)
		}
		ret = TaggedInt(i)
	case []byte:
		if len(t) != 8 {
			return nil, fmt.Errorf("bad int: want 8 bytes, got %d bytes", len(t))
		}
		ret = TaggedInt(binary.BigEndian.Uint64(t))
	case int:
		ret = TaggedInt(t)
	case *int:
		ret = TaggedInt(*t)
	case int64:
		ret = TaggedInt(t)
	case *int64:
		ret = TaggedInt(*t)
	case uint64:
		ret = TaggedInt(t)
	case *uint64:
		ret = TaggedInt(*t)
	default:
		return nil, fmt.Errorf("unexpected type for int: %T", t)
	}

	if err := ret.Valid(); err != nil {
		return nil, err
	}

	return &ClassID{&ret}, nil
}

func (o TaggedInt) String() string {
	return fmt.Sprint(int(o))
}

func (o TaggedInt) Valid() error {
	return nil
}

func (o TaggedInt) Type() string {
	return "int"
}

func (o TaggedInt) Bytes() []byte {
	var ret [8]byte
	binary.BigEndian.PutUint64(ret[:], uint64(o))
	return ret[:]
}

type IClassIDFactory func(any) (*ClassID, error)

var classIDValueRegister = map[string]IClassIDFactory{
	OIDType:  NewOIDClassID,
	UUIDType: NewUUIDClassID,
	IntType:  NewIntClassID,

	ImplIDType: NewImplIDClassID,
}

func RegisterClassIDType(tag uint64, factory IClassIDFactory) error {
	nilVal, err := factory(nil)
	if err != nil {
		return err
	}

	typ := nilVal.Type()
	if _, exists := classIDValueRegister[typ]; exists {
		return fmt.Errorf("class ID type with name %q already exists", typ)
	}

	if err := registerCOMIDTag(tag, nilVal.Value); err != nil {
		return err
	}

	classIDValueRegister[typ] = factory

	return nil
}
