package jsonpointer

import "fmt"

// Error represents all types that may be returned from jsonpointer.
type Error struct {
	derefPrimitive   string
	numParseError    string
	indexOutOfBounds int
	noSuchProperty   string
	parseError       string
}

func (e *Error) Error() string {
	if e.IsDerefPrimitive() {
		return fmt.Sprintf("cannot get property of primitive type: %#v", e.derefPrimitive)
	} else if e.NumParseError() {
		return fmt.Sprintf("cannot parse as number: %#v", e.numParseError)
	} else if e.IndexOutOfBounds() {
		return fmt.Sprintf("index out of bounds: %d", e.indexOutOfBounds)
	} else if e.NoSuchProperty() {
		return fmt.Sprintf("no such property: %#v", e.noSuchProperty)
	} else if e.ParseError() {
		return fmt.Sprintf("error parsing JSON Pointer: %#v", e.parseError)
	} else {
		return "unknown error"
	}
}

// IsDerefPrimitive indicates that the error is due to attempting to dereference
// a property of a primitive type (namely null, a boolean, a number, or a
// string).
func (e *Error) IsDerefPrimitive() bool {
	return e.derefPrimitive != ""
}

// NumParseError indicates that the error is due to attempting to dereference a
// property of an array, but the property wasn't a valid base-10 number.
func (e *Error) NumParseError() bool {
	return e.numParseError != ""
}

// IndexOutOfBounds indicates that the error is due to attempting to dereference
// a element beyond the end of an array.
func (e *Error) IndexOutOfBounds() bool {
	return e.indexOutOfBounds != 0
}

// NoSuchProperty indicates that the error is due to attempting to dereference a
// property of an object that does not have that property.
func (e *Error) NoSuchProperty() bool {
	return e.noSuchProperty != ""
}

// ParseError indicates that the error is due to parsing a string that does not
// correctly represent a JSON Pointer.
func (e *Error) ParseError() bool {
	return e.parseError != ""
}
