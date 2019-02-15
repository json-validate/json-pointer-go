package jsonpointer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var evalTests = []struct {
	tokens []string
	in     interface{}
	out    interface{}
	err    error
}{
	{
		[]string{},
		nil,
		nil,
		nil,
	},
	{
		[]string{},
		true,
		true,
		nil,
	},
	{
		[]string{},
		3.14,
		3.14,
		nil,
	},
	{
		[]string{},
		"a",
		"a",
		nil,
	},
	{
		[]string{},
		"a",
		"a",
		nil,
	},
	{
		[]string{},
		[]interface{}{true, 3.14, "a"},
		[]interface{}{true, 3.14, "a"},
		nil,
	},
	{
		[]string{},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		nil,
	},
	{
		[]string{"a"},
		nil,
		nil,
		&Error{derefPrimitive: "a"},
	},
	{
		[]string{"a"},
		true,
		nil,
		&Error{derefPrimitive: "a"},
	},
	{
		[]string{"a"},
		3.14,
		nil,
		&Error{derefPrimitive: "a"},
	},
	{
		[]string{"a"},
		"a",
		nil,
		&Error{derefPrimitive: "a"},
	},
	{
		[]string{"0"},
		[]interface{}{true, 3.14, "a"},
		true,
		nil,
	},
	{
		[]string{"1"},
		[]interface{}{true, 3.14, "a"},
		3.14,
		nil,
	},
	{
		[]string{"2"},
		[]interface{}{true, 3.14, "a"},
		"a",
		nil,
	},
	{
		[]string{"3"},
		[]interface{}{true, 3.14, "a"},
		nil,
		&Error{indexOutOfBounds: 3},
	},
	{
		[]string{"a"},
		[]interface{}{true, 3.14, "a"},
		nil,
		&Error{numParseError: "a"},
	},
	{
		[]string{"foo"},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		true,
		nil,
	},
	{
		[]string{"bar"},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		3.14,
		nil,
	},
	{
		[]string{"baz"},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		"a",
		nil,
	},
	{
		[]string{"quux"},
		map[string]interface{}{"foo": true, "bar": 3.14, "baz": "a"},
		nil,
		&Error{noSuchProperty: "quux"},
	},
	{
		[]string{"foo", "1", "bar"},
		map[string]interface{}{
			"foo": []interface{}{
				nil,
				map[string]interface{}{
					"bar": "hello, world",
				},
			},
		},
		"hello, world",
		nil,
	},
	{
		[]string{""},
		struct{}{},
		nil,
		&Error{notJSON: &wrappedEmptyStruct},
	},
}

// wrappedEmptyStruct is created as a separate variable because one cannot take
// the address of this value when described as a literal.
var wrappedEmptyStruct = interface{}(struct{}{})

func TestEval(t *testing.T) {
	for i, tt := range evalTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ptr := Ptr{Tokens: tt.tokens}
			val, err := ptr.Eval(tt.in)
			assert.Equal(t, tt.err, err)

			if tt.out != nil {
				assert.Equal(t, &tt.out, val)
			}
		})
	}
}

// These test cases are lifted from RFC6901, Section 5:
//
// https://tools.ietf.org/html/rfc6901#section-5
var newAndParseTests = []struct {
	in  string
	out []string
	err error
}{
	{"", []string{}, nil},
	{"/foo", []string{"foo"}, nil},
	{"/foo/0", []string{"foo", "0"}, nil},
	{"/", []string{""}, nil},
	{"/a~1b", []string{"a/b"}, nil},
	{"/c%d", []string{"c%d"}, nil},
	{"/e^f", []string{"e^f"}, nil},
	{"/g|h", []string{"g|h"}, nil},
	{"/i\\j", []string{"i\\j"}, nil},
	{"/k\"l", []string{"k\"l"}, nil},
	{"/ ", []string{" "}, nil},
	{"/m~0n", []string{"m~n"}, nil},
	{"/o~0~1p/q~1~0r", []string{"o~/p", "q/~r"}, nil},
	{" ", nil, &Error{parseError: " "}},
}

func TestNewAndParse(t *testing.T) {
	for _, tt := range newAndParseTests {
		t.Run(tt.in, func(t *testing.T) {
			// Test that New parses the expected tokens, and that String inverts the
			// process. Verify, when valid, that JSON marshalling/unmarshalling works
			// similarly.
			ptr, err := New(tt.in)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.out, ptr.Tokens)

			// Only attempt to convert back to string, or serialize if no parse error
			// was expected.
			if tt.err == nil {
				assert.Equal(t, tt.in, ptr.String())

				// Test that parsing from JSON yields the same result as New.
				inJSON, err := json.Marshal(tt.in)
				assert.Nil(t, err)

				var ptrFromJSON Ptr
				err = json.Unmarshal(inJSON, &ptrFromJSON)
				assert.Nil(t, err)
				assert.Equal(t, tt.out, ptrFromJSON.Tokens)

				// Test that serializing to JSON yields the same result as String.
				outJSON, err := json.Marshal(ptr)
				assert.Nil(t, err)
				assert.Equal(t, inJSON, outJSON)
			}
		})
	}
}
