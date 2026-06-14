package env

import (
	"fmt"
	"lib/errx"
	"os"
)

func Get[T any](key string, parse func(string) (T, error), fallback T) T {
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

func MustGet[T any](key string, parse func(string) (T, error)) T {
	v := os.Getenv(key)
	result, err := parse(v)
	if err != nil {
		errx.Exit(err, fmt.Sprintf("MustEnv failed to parse %s=%q", key, v))
	}
	return result
}
