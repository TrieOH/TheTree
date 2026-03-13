package csvwriter

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/google/uuid"
)

const timeFormat = "2006-01-02 15:04:05"

// Writer is a generic CSV streamer driven by `csv` struct tags.
// It reflects over T once at construction time and caches the field index,
// so per-row cost is just field access + string conversion.
type Writer[T any] struct {
	fields []fieldMeta
	keySet map[string]struct{}
}

type fieldMeta struct {
	key   string
	index int
}

// New builds a Writer by reflecting over T's `csv` tags.
// Panics on duplicate keys or if T has no csv-tagged fields — programming errors.
func New[T any]() *Writer[T] {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make([]fieldMeta, 0, t.NumField())
	keySet := make(map[string]struct{}, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		key := f.Tag.Get("csv")
		if key == "" || key == "-" {
			continue
		}
		if _, exists := keySet[key]; exists {
			panic(fmt.Sprintf("csvwriter: duplicate csv tag %q on type %s", key, t.Name()))
		}
		fields = append(fields, fieldMeta{key: key, index: i})
		keySet[key] = struct{}{}
	}

	if len(fields) == 0 {
		panic(fmt.Sprintf("csvwriter: type %s has no csv-tagged fields", t.Name()))
	}

	return &Writer[T]{fields: fields, keySet: keySet}
}

// Stream writes a header row then one row per entry in rows.
// selectedKeys controls which columns appear and in what order.
// Returns an error for unknown keys or write failures.
func (w *Writer[T]) Stream(dst io.Writer, selectedKeys []string, rows []T) error {
	if len(selectedKeys) == 0 {
		return fmt.Errorf("csvwriter: at least one column must be selected")
	}

	// Build an ordered index of selected fields, validating all keys up front
	// before writing anything so we never emit a partial header.
	selected := make([]fieldMeta, len(selectedKeys))
	for i, key := range selectedKeys {
		fm, ok := w.fieldByKey(key)
		if !ok {
			return fmt.Errorf("csvwriter: unknown column key %q", key)
		}
		selected[i] = fm
	}

	cw := csv.NewWriter(dst)

	headers := make([]string, len(selected))
	for i, fm := range selected {
		headers[i] = fm.key
	}
	if err := cw.Write(headers); err != nil {
		return fmt.Errorf("csvwriter: writing header: %w", err)
	}

	record := make([]string, len(selected))
	for _, row := range rows {
		rv := reflect.ValueOf(row)
		for i, fm := range selected {
			record[i] = stringify(rv.Field(fm.index))
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("csvwriter: writing row: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// ValidKeys returns all column keys available for this writer.
func (w *Writer[T]) ValidKeys() map[string]struct{} {
	out := make(map[string]struct{}, len(w.keySet))
	for k := range w.keySet {
		out[k] = struct{}{}
	}
	return out
}

func (w *Writer[T]) fieldByKey(key string) (fieldMeta, bool) {
	for _, fm := range w.fields {
		if fm.key == key {
			return fm, true
		}
	}
	return fieldMeta{}, false
}

// stringify converts a reflected field value to its CSV string representation.
// Handles the types that appear in export rows: string, uuid.UUID, time.Time,
// enums (underlying string), and pointer variants of all of the above.
func stringify(v reflect.Value) string {
	// Unwrap pointer — nil pointer becomes empty string
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	iface := v.Interface()

	switch val := iface.(type) {
	case string:
		return val
	case uuid.UUID:
		return val.String()
	case time.Time:
		if val.IsZero() {
			return ""
		}
		return val.Format(timeFormat)
	case bool:
		if val {
			return "true"
		}
		return "false"
	}

	// Catch enums and any other type with an underlying string kind
	if v.Kind() == reflect.String {
		return v.String()
	}

	return fmt.Sprintf("%v", iface)
}
