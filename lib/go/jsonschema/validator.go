// Package jsonschema provides generic JSON Schema validation for any JSON instance.
package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validate validates an instance JSON document against a JSON Schema document.
// Returns nil if valid, or a structured error listing each validation failure.
// An empty or null schema skips validation (passthrough).
func Validate(schema, instance json.RawMessage) error {
	if len(schema) == 0 || string(schema) == "{}" || string(schema) == "null" {
		return nil
	}

	schemaDoc, err := jsonschema.UnmarshalJSON(strings.NewReader(string(schema)))
	if err != nil {
		return fmt.Errorf("invalid schema json: %w", err)
	}

	c := jsonschema.NewCompiler()
	err = c.AddResource("schema.json", schemaDoc)
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	sch, err := c.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("invalid schema: %w", err)
	}

	inst, err := jsonschema.UnmarshalJSON(strings.NewReader(string(instance)))
	if err != nil {
		return fmt.Errorf("invalid instance json: %w", err)
	}

	err = sch.Validate(inst)
	if err == nil {
		return nil
	}

	return formatValidationError(err)
}

func formatValidationError(err error) error {
	if validationErr, ok := errors.AsType[*jsonschema.ValidationError](err); ok {
		var msg strings.Builder
		msg.WriteString("validation failed:")
		for _, detail := range validationErr.DetailedOutput().Errors {
			msg.WriteString("\n  - ")
			msg.WriteString(detail.InstanceLocation)
			msg.WriteString(": ")
			msg.WriteString(outputErrorString(detail.Error))
		}
		return fmt.Errorf("%s", msg.String())
	}
	return fmt.Errorf("validation failed:\n  - %s", err.Error())
}

func outputErrorString(oe *jsonschema.OutputError) string {
	if oe == nil {
		return "invalid value"
	}
	return oe.String()
}
