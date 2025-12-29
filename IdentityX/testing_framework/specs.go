package testing

// ============================================================================
// TEST SPECS - Declarative test definitions
// ============================================================================

type ValidationSpec struct {
	Name   string
	Email  string
	Pass   string
	Errors []string
}

type PasswordSpec struct {
	Name     string
	Password string
	Errors   []string
}

// Test data
var (
	ValidPassword = "Str0ngP4$$!"

	ValidationTests = []ValidationSpec{
		{
			Name:   "NoEmail",
			Email:  "",
			Pass:   ValidPassword,
			Errors: []string{"(email) is required"},
		},
		{
			Name:   "InvalidEmail",
			Email:  "not-an-email",
			Pass:   ValidPassword,
			Errors: []string{"valid email address"},
		},
		{
			Name:   "NoPassword",
			Email:  "test@mail.com",
			Pass:   "",
			Errors: []string{"(password) is required"},
		},
	}

	WeakPasswordTests = []PasswordSpec{
		{"OnlyLetters", "abc", []string{"uppercase", "number", "symbol"}},
		{"LettersNumber", "abc3", []string{"uppercase", "symbol"}},
		{"LettersSymbol", "abc#", []string{"uppercase", "number"}},
		{"LettersUppercase", "Abc", []string{"number", "symbol"}},
		{"NoNumber", "Abc#", []string{"number"}},
		{"NoSymbol", "Abc3", []string{"symbol"}},
		{"TooShort", "Abc#3", []string{"at least 8 characters"}},
	}
)
