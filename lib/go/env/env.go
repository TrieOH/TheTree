package env

import (
	"os"
)

func Get[T any](key string, parse func(string) (T, error), fallback T) T { //nolint:ireturn
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	result, err := parse(v)
	if err != nil {
		return fallback
	}
	return result
}
