package tegnbuilder

import "fmt"

type TegnParameterOption func(self *TegnParameter)
type TegnValidator func (self *TegnParameter, newValue string) bool

// Represents a parameter specification for a Tegn.
type TegnParameter struct {
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

func NewTegnParameter(name string, paramType TegnParameterType, opts ...TegnParameterOption) TegnParameter{
	p := TegnParameter{
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
func WithAvailability(available TegnAvailability) TegnParameterOption {
	return func(p *TegnParameter) {
		p.available = available
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
	return p.paramType
}

func (p *TegnParameter) GetAvailable() TegnAvailability {
	return p.available
}

func (p *TegnParameter) GetValidator() TegnValidator {
	return p.validator
}

func (p *TegnParameter) Validate(newValue string) bool {
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
	return fmt.Sprintf("%d", v)
}

func TegnParameterToBool(s string) bool {
	return s == "y"
}

func TegnParameterToInt(s string, fallback int) int {
	var result *int
	_, err := fmt.Sscanf(s, "%d", result)
	if err != nil || result == nil {
		return fallback
	}

	return *result
}
