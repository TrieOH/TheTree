package utils

import "encoding/json"

// MapTo unmarshals src into a value of type T and returns a pointer to it.
func MapTo[T any](dst *T, src json.RawMessage) error {
	return json.Unmarshal(src, dst)
}
