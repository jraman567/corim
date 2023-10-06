// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package comid

import (
	"github.com/veraison/corim/extensions"
)

type IComidValidator interface {
	ValidComid(*Comid) error
}

type ITriplesValidator interface {
	ValidTriples(*Triples) error
}

type IMvalValidator interface {
	ValidMval(*Mval) error
}

type IFlagsMapValidator interface {
	ValidFlagsMap(*FlagsMap) error
}

type IFlagSetter interface {
	AnySet() bool
	SetTrue(Flag)
	SetFalse(Flag)
	Clear(Flag)
	Get(Flag) *bool
}

type Extensions struct {
	extensions.Extensions
}

func (o *Extensions) ValidComid(comid *Comid) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IComidValidator)
	if ok {
		if err := ev.ValidComid(comid); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) ValidTriples(triples *Triples) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(ITriplesValidator)
	if ok {
		if err := ev.ValidTriples(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) ValidMval(triples *Mval) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IMvalValidator)
	if ok {
		if err := ev.ValidMval(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) ValidFlagsMap(triples *FlagsMap) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IFlagsMapValidator)
	if ok {
		if err := ev.ValidFlagsMap(triples); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) SetTrue(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IExtensionsValue.(IFlagSetter)
	if ok {
		ev.SetTrue(flag)
	}
}

func (o *Extensions) SetFalse(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IExtensionsValue.(IFlagSetter)
	if ok {
		ev.SetFalse(flag)
	}
}

func (o *Extensions) Clear(flag Flag) {
	if !o.HaveExtensions() {
		return
	}

	ev, ok := o.IExtensionsValue.(IFlagSetter)
	if ok {
		ev.Clear(flag)
	}
}

func (o *Extensions) Get(flag Flag) *bool {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IFlagSetter)
	if ok {
		return ev.Get(flag)
	}

	return nil
}

func (o *Extensions) AnySet() bool {
	if !o.HaveExtensions() {
		return false
	}

	ev, ok := o.IExtensionsValue.(IFlagSetter)
	if ok {
		return ev.AnySet()
	}

	return false
}
