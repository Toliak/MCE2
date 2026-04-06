package tegnbuilder

import (
	"fmt"
	"strconv"
)

type TegnParameterOption func(self *TegnParameter)
type TegnValidator func (self *TegnParameter, newValue string) error

// Represents a parameter specification for a Tegn.
type TegnParameter struct {
	id string
	name string
	description string

	// Current parameter DefaultValue.
	defaultValue string

	// Parameter's data type.
	paramType TegnParameterType

	// Indicates the availability and contains the reason, if not available.
	available TegnAvailability

	// Validates the value
	validator TegnValidator
}

func NewTegnParameter(id string, name string, paramType TegnParameterType, opts ...TegnParameterOption) TegnParameter{
	p := TegnParameter{
		id: id,
		name:        name,
		paramType: paramType,
	}
	for _, opt := range opts {
		opt(&p)
	}

	return p
}

// WithDefaultValue sets the private value field.
func WithDefaultValue(val string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.defaultValue = val
	}
}

// WithAvailability sets the availability status.
func WithAvailability(status bool, reason string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.available = TegnAvailability{
			Available: status,
			Reason: reason,
		}
	}
}

// WithAvailability sets the availability status.
func WithAvailabilityTrue() TegnParameterOption {
	return func(p *TegnParameter) {
		p.available = NewTegnAvailable()
	}
}

// WithAvailability sets the availability status.
func WithAvailabilityFalse(reason string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.available = NewTegnNotAvailable(reason)
	}
}

// WithDescription overrides the description (useful if you want to set it optionally).
func WithDescription(desc string) TegnParameterOption {
	return func(p *TegnParameter) {
		p.description = desc
	}
}

// WithDescription overrides the description (useful if you want to set it optionally).
func WithValidator(validator TegnValidator) TegnParameterOption {
	return func(p *TegnParameter) {
		p.validator = validator
	}
}

func ReadOnlyValidator(param *TegnParameter, newValue string) error {
	return fmt.Errorf("The field '%s' is read-only", param.GetName())
}

// WithDescription overrides the description (useful if you want to set it optionally).
func WithReadOnlyValidator() TegnParameterOption {
	return func(p *TegnParameter) {
		p.validator = ReadOnlyValidator
	}
}

func (p *TegnParameter) GetID() string {
	return p.id
}

func (p *TegnParameter) GetName() string {
	return p.name
}

func (p *TegnParameter) GetDescription() string {
	return p.description
}

func (p *TegnParameter) GetDefaultValue() string {
	return p.defaultValue
}

func (p *TegnParameter) GetParamType() TegnParameterType {
	// TODO: about "param type" and "validate". Consolidate that somehow.
	// So, two ways:
	// 1. Validate checks the parameter based on the type and custom validator
	// 2. Remove param type and leave the custom validator. But how to store templates?
	// TODO: resolved??
	return p.paramType
}

func (p *TegnParameter) GetAvailable() TegnAvailability {
	return p.available
}

func (p *TegnParameter) GetValidator() TegnValidator {
	return p.validator
}

func (p *TegnParameter) validateBasedOnParameterType(newValue string) error {
	switch p.paramType {
	case TegnParameterTypeInt:
		_, err := strconv.Atoi(newValue)
		if err != nil {
			return fmt.Errorf("Unable to convert '%s' to integer", newValue)
		}
	case TegnParameterTypeBool:
		if newValue != "y" && newValue != "n" {
			return fmt.Errorf("Unable to convert '%s' to boolean. Expected 'y' or 'n'", newValue)
		}
	case TegnParameterTypeString:
		// usually always valid
	}

	return nil
}

func (p *TegnParameter) Validate(newValue string) error {
	if err := p.validateBasedOnParameterType(newValue); err != nil {
		return err
	}

	if p.validator == nil {
		return nil
	}
	return p.validator(p, newValue)
}

// Data type of a TegnParameter.
type TegnParameterType int

const (
	TegnParameterTypeInt TegnParameterType = iota
	TegnParameterTypeBool
	TegnParameterTypeString
)

func (t TegnParameterType) String() string {
	switch t {
	case TegnParameterTypeInt:
		return "int"
	case TegnParameterTypeBool:
		return "bool"
	case TegnParameterTypeString:
		return "string"
	default:
		return "invalid"
	}
}

func TegnParameterFromBool(v bool) string {
	if v {
		return "y"
	}

	return "n"
}

func TegnParameterFromInt(v int) string {
	return strconv.Itoa(v)
}

func TegnParameterToBool(s string) bool {
	return s == "y"
}

func TegnParameterToInt(s string, fallback int) int {
	var result int
	_, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}

	return result
}
