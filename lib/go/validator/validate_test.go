package validator

import (
	"errors"
	"testing"
)

// Validate works without SetupValidator: an operation called directly, as
// the integration tests do, used to panic here.
func TestValidateBuildsTheValidatorOnFirstUse(t *testing.T) {
	type input struct {
		Name string `json:"name" validate:"required"`
	}
	err := Validate(input{})
	var verr *Error
	if !errors.As(err, &verr) || len(verr.Fields) != 1 {
		t.Fatalf("Validate(empty) = %v, want one failed field", err)
	}
	if Validate(input{Name: "ada"}) != nil {
		t.Fatal("a valid input failed")
	}
}
