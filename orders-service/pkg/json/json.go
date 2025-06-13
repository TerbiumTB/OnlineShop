package json

import (
	"encoding/json"
	"io"
)

func ToJSON(i any, w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")

	return e.Encode(i)
}

func FromJSON(i any, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
