package utils

import (
	"encoding/json"
	"io"
)

func Decode(r io.Reader, res interface{}) error {
	return json.NewDecoder(r).Decode(&res)
}

func Encode(w io.Writer, res interface{}) error {
	return json.NewEncoder(w).Encode(&res)
}
func EncodeMultipleStructs(w io.Writer, res ...interface{}) error {
	return json.NewEncoder(w).Encode(&res)
}
