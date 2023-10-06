// Copyright 2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package corim

import (
	"github.com/veraison/corim/extensions"
)

type IEntityValidator interface {
	ValidEntity(*Entity) error
}

type ICorimValidator interface {
	ValidCorim(*UnsignedCorim) error
}

type ISignerValidator interface {
	ValidSigner(*Signer) error
}

type Extensions struct {
	extensions.Extensions
}

func (o *Extensions) ValidEntity(entity *Entity) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(IEntityValidator)
	if ok {
		if err := ev.ValidEntity(entity); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) ValidCorim(c *UnsignedCorim) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(ICorimValidator)
	if ok {
		if err := ev.ValidCorim(c); err != nil {
			return err
		}
	}

	return nil
}

func (o *Extensions) ValidSigner(signer *Signer) error {
	if !o.HaveExtensions() {
		return nil
	}

	ev, ok := o.IExtensionsValue.(ISignerValidator)
	if ok {
		if err := ev.ValidSigner(signer); err != nil {
			return err
		}
	}

	return nil
}
