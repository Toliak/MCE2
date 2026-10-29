// Package optional provides an explicit representation of values that may be absent.
//
// Use Optional[T] instead of nil pointers for optional values in internal APIs.
// Pointers are expected to be non-nil by default (see configapp/NULLABILITY.md).
package optional

// Optional represents a value that may or may not be present.
type Optional[T any] struct {
	value   T
	present bool
}

// Some returns an Optional containing value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, present: true}
}

// None returns an empty Optional.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// FromPtr converts a pointer into an Optional. A nil pointer yields None.
func FromPtr[T any](ptr *T) Optional[T] {
	if ptr == nil {
		return None[T]()
	}
	return Some(*ptr)
}

// IsPresent reports whether the Optional contains a value.
func (o Optional[T]) IsPresent() bool {
	return o.present
}

// IsEmpty reports whether the Optional is empty.
func (o Optional[T]) IsEmpty() bool {
	return !o.present
}

// Get returns the contained value and a boolean indicating presence.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.present
}

// Or returns the contained value if present, otherwise def.
func (o Optional[T]) Or(def T) T {
	if o.present {
		return o.value
	}
	return def
}

// Value returns the contained value. It panics if the Optional is empty.
func (o Optional[T]) Value() T {
	if !o.present {
		panic("optional: Value() called on an empty Optional")
	}
	return o.value
}

// Map applies f to the contained value and returns an Optional of the result.
func Map[T, U any](o Optional[T], f func(T) U) Optional[U] {
	if !o.present {
		return None[U]()
	}
	return Some(f(o.value))
}
