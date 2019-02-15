# json-pointer [![documentation][badge]][godoc]

This library is a Golang implementation of [RFC6901][rfc6901] "JavaScript Object
Notation (JSON) Pointer".

[badge]: https://godoc.org/github.com/ucarion/json-pointer?status.svg
[godoc]: https://godoc.org/github.com/ucarion/json-pointer

## Usage

This library does not perform JSON encoding/decoding. Instead, you must provide
data in the default Golang format for JSON data, namely:

```txt
bool, for JSON booleans
float64, for JSON numbers
string, for JSON strings
[]interface{}, for JSON arrays
map[string]interface{}, for JSON objects
nil for JSON null
```

_(from the docs for [`encoding/json#Unmarshal`][encoding/json#unmarshal])_

To create a JSON Pointer, use `New`. To evaluate a pointer, use `Eval`:

```golang
ptr, err := jsonpointer.New("/foo/1/bar")

data := map[string]interface{}{
  "foo": []interface{}{
    nil,
    map[string]interface{}{
      "bar": "hello, world",
    },
  },
}

val, err := ptr.Eval(&data)
fmt.Println(val) // outputs "hello, world"
```

[rfc6901]: https://tools.ietf.org/html/rfc6901
[encoding/json#unmarshal]: https://golang.org/pkg/encoding/json/#Unmarshal
