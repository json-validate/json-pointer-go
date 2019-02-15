package jsonpointer

import "errors"

// ErrInvalidPtr indicates that a string inputted to New does not correspond to
// a representation of any JSON Pointer.
var ErrInvalidPtr = errors.New("invalid JSON Pointer")

// ErrEvalPtr indicates that a JSON Pointer referred to a nonexistent property
// of data, or that the data is not a Golang representation of JSON data.
var ErrEvalPtr = errors.New("error evaluating JSON Pointer")
