package jsonpointer

import (
	"strconv"
)

// Ptr represents a JSON Pointer in parsed form.
type Ptr struct {
	// The "reference tokens" (see RFC6901, Section 4) of the pointer. Special
	// sequences such as "~0" and "~1" are already parsed into "~" and "/",
	// respectively.
	Tokens []string
}

// Eval evaluates a Ptr against a document, returning a (Golang) pointer into
// that document.
//
// Errors, if returned, will be instances of Error from this package.
func (p Ptr) Eval(doc interface{}) (*interface{}, error) {
	i := 0

	for i < len(p.Tokens) {
		switch v := doc.(type) {
		case nil, bool, float64, string:
			return nil, &Error{derefPrimitive: true}
		case []interface{}:
			n, err := strconv.ParseInt(p.Tokens[i], 10, 0)
			if err != nil {
				return nil, &Error{numParseError: p.Tokens[i]}
			}

			if n < 0 || int(n) >= len(v) {
				return nil, &Error{indexOutOfBounds: int(n)}
			}

			doc = v[n]
		case map[string]interface{}:
			var ok bool
			doc, ok = v[p.Tokens[i]]

			if !ok {
				return nil, &Error{noSuchProperty: p.Tokens[i]}
			}
		}

		i++
	}

	return &doc, nil
}
