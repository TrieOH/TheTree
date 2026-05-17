package errx

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

func Exit(err error, msg string) {
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error(msg, "err", err)
		os.Exit(1)
	}
}

func MustEnv[T any](key string, parse func(string) (T, error)) T {
	v := os.Getenv(key)
	result, err := parse(v)
	if err != nil {
		Exit(err, fmt.Sprintf("MustEnv failed to parse %s=%q", key, v))
	}
	return result
}

// MustProvide validates that all fields of the given struct are non-zero/non-nil.
// Panics with a list of all missing fields if any are zero.
// Returns the pointer typed, so it can wrap the return directly.
//
// Usage:
//
//	func NewCommandService(deps Deps) *CommandService {
//	    return fun.MustProvide(&CommandService{
//	        users:    deps.Users,
//	        accounts: deps.Accounts,
//	    })
//	}
func MustProvide[T any](v *T) *T {
	if v == nil {
		panic("fun.MustProvide: nil pointer provided")
	}

	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()

	var missing []string
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanInterface() {
			continue // unexported — skip
		}
		if isZero(field.Interface()) {
			missing = append(missing, rt.Field(i).Name)
		}
	}

	if len(missing) > 0 {
		panic("fun.MustProvide: missing fields in " + rt.Name() + ": " + strings.Join(missing, ", "))
	}

	return v
}

func isZero(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice,
		reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return v.IsZero()
}
