// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Wallet Wallet with some coins.
//
// swagger:model wallet
type Wallet []uint32

// Validate validates this wallet
func (m Wallet) Validate(formats strfmt.Registry) error {
	var res []error

	iWalletSize := int64(len(m))

	if err := validate.MaxItems("", "body", iWalletSize, 1000); err != nil {
		return err
	}

	if err := validate.UniqueItems("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this wallet based on context it is used
func (m Wallet) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
