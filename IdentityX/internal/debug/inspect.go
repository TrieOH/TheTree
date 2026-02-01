package debug

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func InspectError(err error) {
	fmt.Printf("=== ERROR CHAIN ===\n")
	current := err
	for i := 0; current != nil && i < 10; i++ {
		fmt.Printf("[%d] Type: %T\n", i, current)
		fmt.Printf("     Error(): %v\n", current.Error())

		// Check if it implements Unwrap
		if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
			next := unwrapper.Unwrap()
			if next == current {
				fmt.Println("     (self-referencing unwrap, stopping)")
				break
			}
			current = next
		} else {
			fmt.Println("     (no Unwrap method)")
			break
		}
	}
	fmt.Println("===================")
}

func DeepInspectError(err error) string {
	fmt.Printf("=== DEEP INSPECT ===\n")
	if err == nil {
		return "nil"
	}

	var result string

	// Basic info
	result += fmt.Sprintf("=== Error Inspection ===\n")
	result += fmt.Sprintf("Error(): %s\n", err.Error())
	result += fmt.Sprintf("Type: %T\n", err)
	result += fmt.Sprintf("Reflect: %v\n", reflect.TypeOf(err))

	// Check if it's a struct and show fields
	v := reflect.ValueOf(err)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		result += fmt.Sprintf("Fields:\n")
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanInterface() {
				result += fmt.Sprintf("  %s: %v (type: %T)\n",
					t.Field(i).Name,
					field.Interface(),
					field.Interface())
			}
		}
	}

	// Check pgx specific types
	result += fmt.Sprintf("\n=== pgx Checks ===\n")

	if errors.Is(err, pgx.ErrNoRows) {
		result += "✓ Is pgx.ErrNoRows\n"
	}
	if errors.Is(err, pgx.ErrTooManyRows) {
		result += "✓ Is pgx.ErrTooManyRows\n"
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		result += fmt.Sprintf("✓ Is *pgconn.PgError:\n")
		result += fmt.Sprintf("    Code: %s\n", pgErr.Code)
		result += fmt.Sprintf("    Message: %s\n", pgErr.Message)
	}

	// Unwrap chain
	result += fmt.Sprintf("\n=== Unwrap Chain ===\n")
	current := err
	depth := 0
	for {
		result += fmt.Sprintf("%d: %T\n", depth, current)
		unwrapped := errors.Unwrap(current)
		if unwrapped == nil {
			break
		}
		current = unwrapped
		depth++
	}

	fmt.Println("===================")
	return result
}

func ClassifyPgxError(err error) {
	fmt.Printf("=== CLASSIFY PGX ERROR ===\n")
	// pgx.ErrNoRows - query returned no rows
	{
		var counter int
		if err == pgx.ErrNoRows {
			fmt.Println("USING EQUALITY Type: pgx.ErrNoRows")
			counter++
		}

		if errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("USING IS Type: pgx.ErrNoRows")
			counter++
		}
		if counter == 2 {
			return
		}
	}

	// pgx.ErrTooManyRows - query returned multiple rows when expecting one
	if err == pgx.ErrTooManyRows {
		fmt.Println("Type: pgx.ErrTooManyRows")
		return
	}

	// PostgreSQL error with SQLSTATE code
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Printf("Type: *pgconn.PgError\n")
		fmt.Printf("  SQLSTATE: %s\n", pgErr.Code)
		fmt.Printf("  Message:  %s\n", pgErr.Message)
		fmt.Printf("  Detail:   %s\n", pgErr.Detail)
		fmt.Printf("  Table:    %s\n", pgErr.TableName)
		fmt.Printf("  Column:   %s\n", pgErr.ColumnName)
		return
	}

	// Connection errors
	var connErr *pgconn.ConnectError
	if errors.As(err, &connErr) {
		fmt.Println("Type: *pgconn.ConnectError")
		return
	}

	// Context errors
	if errors.Is(err, context.Canceled) {
		fmt.Println("Type: context.Canceled")
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("Type: context.DeadlineExceeded")
		return
	}

	// Fallback
	fmt.Printf("Unknown type: %T\n", err)
	fmt.Println("===================")
}
