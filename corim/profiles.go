// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jraman567/corim/comid"
	"github.com/jraman567/corim/extensions"
	"github.com/veraison/eat"
	"github.com/veraison/go-cose"
)

// SignedCorimMapExtensionPoints is a list of extension.Point's valid for a
// SignedCorim.
var SignedCorimMapExtensionPoints = []extensions.Point{
	ExtSigner,
	ExtUnsignedCorim,
	ExtEntity,
}

// UnsignedCorimMapExtensionPoints is a list of extension.Point's valid for a
// UnsignedCorim.
var UnsignedCorimMapExtensionPoints = []extensions.Point{
	ExtUnsignedCorim,
	ExtEntity,
}

// ComidMapExtensionPoints is a list of extension.Point's valid for a comid.Comid.
var ComidMapExtensionPoints = []extensions.Point{
	comid.ExtComid,
	comid.ExtEntity,
	comid.ExtTriples,
	comid.ExtReferenceValue,
	comid.ExtReferenceValueFlags,
	comid.ExtEndorsedValue,
	comid.ExtEndorsedValueFlags,
}

// AllExtensionPoints is a list of all valid extension.Point's
var AllExtensionPoints = make(map[extensions.Point]bool) // populated inside init() below

// UnmarshalSignedCorimFromCBOR unmarshals a SignedCorim from provided
// CBOR data. If there are extensions associated with the profile specified by
// the data, they will be registered with the UnsignedCorim before it is
// unmarshaled.
func UnmarshalSignedCorimFromCBOR(buf []byte) (*SignedCorim, error) {
	message := cose.NewSign1Message()

	if err := message.UnmarshalCBOR(buf); err != nil {
		return nil, fmt.Errorf("failed CBOR decoding for COSE-Sign1 signed CoRIM: %w", err)
	}

	profiled := struct {
		Profile *eat.Profile `cbor:"3,keyasint,omitempty"`
	}{}

	if err := dm.Unmarshal(message.Payload, &profiled); err != nil {
		return nil, err
	}

	ret := GetSignedCorim(profiled.Profile)
	if err := ret.FromCOSE(buf); err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalUnsignedCorimFromCBOR unmarshals an UnsignedCorim from provided
// CBOR data. If there are extensions associated with the profile specified by
// the data, they will be registered with the UnsignedCorim before it is
// unmarshaled.
func UnmarshalUnsignedCorimFromCBOR(buf []byte) (*UnsignedCorim, error) {
	profiled := struct {
		Profile *eat.Profile `cbor:"3,keyasint,omitempty"`
	}{}

	if err := dm.Unmarshal(buf, &profiled); err != nil {
		return nil, err
	}

	ret := GetUnsignedCorim(profiled.Profile)
	if err := ret.FromCBOR(buf); err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalUnsignedCorimFromJSON unmarshals an UnsignedCorim from provided
// JSON data. If there are extensions associated with the profile specified by
// the data, they will be registered with the UnsignedCorim before it is
// unmarshaled.
func UnmarshalUnsignedCorimFromJSON(buf []byte) (*UnsignedCorim, error) {
	profiled := struct {
		Profile *eat.Profile `json:"profile,omitempty"`
	}{}

	if err := json.Unmarshal(buf, &profiled); err != nil {
		return nil, err
	}

	ret := GetUnsignedCorim(profiled.Profile)
	if err := ret.FromJSON(buf); err != nil {
		return nil, err
	}

	return ret, nil
}

// UnmarshalComidFromCBOR unmarshals a comid.Comid from provided CBOR data. If
// there are extensions associated with the profile specified by the data, they
// will be registered with the comid.Comid before it is unmarshaled.
func UnmarshalComidFromCBOR(buf []byte, profileID *eat.Profile) (*comid.Comid, error) {
	var ret *comid.Comid

	profile, ok := GetProfile(profileID)
	if ok {
		ret = profile.GetComid()
	} else {
		ret = comid.NewComid()
	}

	if err := ret.FromCBOR(buf); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetSingedCorim returns a pointer to a new SingedCorim instance. If there
// are extensions associated with the provided profileID, they will be
// registered with the instance.
func GetSignedCorim(profileID *eat.Profile) *SignedCorim {
	var ret *SignedCorim

	if profileID == nil {
		ret = NewSignedCorim()
	} else {
		profile, ok := GetProfile(profileID)
		if !ok {
			// unknown profile -- treat here like an unprofiled
			// CoRIM. While the CoRIM spec states that unknown
			// profiles should be rejected, we're not actually
			// validating the profile here, just trying to identify
			// any extensions we may need to load. Profile
			// validation is left up to the calling code, as a
			// profile only needs to be registered here if it
			// defines extensions. Profiles that do not add any
			// additional fields may not be registered.
			ret = NewSignedCorim()
		} else {
			ret = profile.GetSignedCorim()
		}
	}

	return ret
}

// GetUnsignedCorim returns a pointer to a new UnsignedCorim instance. If there
// are extensions associated with the provided profileID, they will be
// registered with the instance.
func GetUnsignedCorim(profileID *eat.Profile) *UnsignedCorim {
	var ret *UnsignedCorim

	if profileID == nil {
		ret = NewUnsignedCorim()
	} else {
		profile, ok := GetProfile(profileID)
		if !ok {
			// unknown profile -- treat here like an unprofiled
			// CoRIM. While the CoRIM spec states that unknown
			// profiles should be rejected, we're not actually
			// validating the profile here, just trying to identify
			// any extensions we may need to load. Profile
			// validation is left up to the calling code, as a
			// profile only needs to be registered here if it
			// defines extensions. Profiles that do not add any
			// additional fields may not be registered.
			ret = NewUnsignedCorim()
		} else {
			ret = profile.GetUnsignedCorim()
		}
	}

	return ret
}

// Profile associates an EAT profile ID with a set of extensions. It allows
// obtaining new CoRIM and CoMID structures that had associated extensions
// registered.
type Profile struct {
	ID            *eat.Profile
	MapExtensions extensions.Map
}

// GetComid returns a pointer to a new comid.Comid that had the Profile's
// extensions (if any) registered.
func (o *Profile) GetComid() *comid.Comid {
	ret := comid.NewComid()
	o.registerExtensions(ret, ComidMapExtensionPoints)
	return ret
}

// GetUnsignedCorim returns a pointer to a new UnsignedCorim that had the
// Profile's extensions (if any) registered.
func (o *Profile) GetUnsignedCorim() *UnsignedCorim {
	ret := NewUnsignedCorim()
	ret.Profile = o.ID
	o.registerExtensions(ret, UnsignedCorimMapExtensionPoints)
	return ret
}

// GetSignedCorim returns a pointer to a new SignedCorim that had the
// Profile's extensions (if any) registered.
func (o *Profile) GetSignedCorim() *SignedCorim {
	ret := NewSignedCorim()
	ret.UnsignedCorim.Profile = o.ID
	o.registerExtensions(ret, SignedCorimMapExtensionPoints)
	return ret
}

func (o *Profile) registerExtensions(e iextensible, points []extensions.Point) {
	exts := extensions.NewMap()
	for _, p := range points {
		if v, ok := o.MapExtensions[p]; ok {
			exts[p] = v
		}
	}

	if err := e.RegisterExtensions(exts); err != nil {
		// exts is a subset of o.MapExtensions which have been
		// validated when the profile was registered, so we should never
		// get here.
		panic(err)
	}
}

// RegisterProfile registers a set of extensions with the specified profile. If
// the profile has already been registered, or if the extensions are invalid,
// an error is returned.
func RegisterProfile(id *eat.Profile, exts extensions.Map) error {
	strID, err := id.Get()
	if err != nil {
		return err
	}

	if _, ok := profilesRegister[strID]; ok {
		return fmt.Errorf("profile with id %q already registered", strID)
	}

	for p, v := range exts {
		if _, ok := AllExtensionPoints[p]; !ok {
			return fmt.Errorf("%w: %q", extensions.ErrUnexpectedPoint, p)
		}

		if reflect.TypeOf(v).Kind() != reflect.Pointer {
			return fmt.Errorf("attempting to register a non-pointer IMapValue for %q", p)
		}
	}

	profilesRegister[strID] = Profile{ID: id, MapExtensions: exts}

	return nil
}

// UnregisterProfile ensures there are no extensions registered for the
// specified profile ID. Returns true if extensions were previously registered
// and have been removed, and false otherwise.
func UnregisterProfile(id *eat.Profile) bool {
	if id == nil {
		return false
	}

	strID, err := id.Get()
	if err != nil {
		return false
	}

	if _, ok := profilesRegister[strID]; ok {
		delete(profilesRegister, strID)
		return true
	}

	return false
}

// GetProfile returns the Profile associated with the specified ID, or an empty
// profile if no Profile has been registered for the id. The second return
// value indicates whether a profile for the ID has been found.
func GetProfile(id *eat.Profile) (Profile, bool) {
	if id == nil {
		return Profile{}, false
	}

	strID, err := id.Get()
	if err != nil {
		return Profile{}, false
	}

	prof, ok := profilesRegister[strID]
	return prof, ok
}

type iextensible interface {
	RegisterExtensions(exts extensions.Map) error
}

var profilesRegister = make(map[string]Profile)

func init() {
	for _, p := range SignedCorimMapExtensionPoints {
		AllExtensionPoints[p] = true
	}

	for _, p := range UnsignedCorimMapExtensionPoints {
		AllExtensionPoints[p] = true
	}

	for _, p := range ComidMapExtensionPoints {
		AllExtensionPoints[p] = true
	}
}
