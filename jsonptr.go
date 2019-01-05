package jsonpointer

import (
	"fmt"
	"strconv"
	"strings"
)

// Ptr represents a JSON Pointer in parsed form.
type Ptr struct {
	// The "reference tokens" (see RFC6901, Section 4) of the pointer. Special
	// sequences such as "~0" and "~1" are already parsed into "~" and "/",
	// respectively.
	Tokens []string
}

// New parses a JSON Pointer represented as a string value.
//
// This function will handle JSON Pointer escape sequences, converting "~0" to
// "~" and "~1" to "/".
//
// If s is not empty and does not begin with "/", an error will be returned.
func New(s string) (Ptr, error) {
	// From the ABNF syntax of JSON Pointer, the only valid initial character for
	// a JSON Pointer is "/". Empty strings are acceptable.
	//
	// https://tools.ietf.org/html/rfc6901#section-3
	//
	// Other than this limitation, all strings are valid JSON Pointers.
	if s == "" {
		return Ptr{Tokens: []string{}}, nil
	}

	if !strings.HasPrefix(s, "/") {
		return Ptr{}, &Error{parseError: s}
	}

	tokens := strings.Split(s, "/")[1:]
	for i, token := range tokens {
		// This sequence of replacements follows the instructions from:
		//
		// https://tools.ietf.org/html/rfc6901#section-4
		token = strings.Replace(token, "~1", "/", -1)
		token = strings.Replace(token, "~0", "~", -1)
		tokens[i] = token
	}

	return Ptr{Tokens: tokens}, nil
}

// String is an implementation of Stringer for Ptr.
//
// This functions acts as the inverse of New.
func (p Ptr) String() string {
	// Special case: empty sequence of tokens is represented as empty string.
	if len(p.Tokens) == 0 {
		return ""
	}

	parts := make([]string, len(p.Tokens))
	for i, token := range p.Tokens {
		token = strings.Replace(token, "~", "~0", -1)
		token = strings.Replace(token, "/", "~1", -1)
		parts[i] = token
	}

	return fmt.Sprintf("/%s", strings.Join(parts, "/"))
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
