package testing

import (
	"net/http"
	"testing"
)

func registerSuccess(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		r := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().Value("message")

		r.String().Equal("Registered user")
	}
}

func registerNoEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "",
				"password":  "N0Email#S4d",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(email) is required")

		obj.Value("code").Number().Equal(400)
	}
}

func registerNoPassword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "nopass@mail.com",
				"password":  "",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(password) is required")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLetters() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass1@mail.com",
				"password":  "abc",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(3)
		trace.Element(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Element(1).String().Contains("(password) must contain at least one number")
		trace.Element(2).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersNumber() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass2@mail.com",
				"password":  "abc3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(2)
		trace.Element(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Element(1).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersSymbol() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass3@mail.com",
				"password":  "abc#",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(2)
		trace.Element(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Element(1).String().Contains("(password) must contain at least one number")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass4@mail.com",
				"password":  "Abc",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(2)
		trace.Element(0).String().Contains("(password) must contain at least one number")
		trace.Element(1).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersSymbolUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass5@mail.com",
				"password":  "Abc#",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(password) must contain at least one number")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersNumberUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass6@mail.com",
				"password":  "Abc3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().Equal(400)
	}
}

func registerWeakPasswordLettersNumberSymbolUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "weakpass7@mail.com",
				"password":  "Abc#3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(password) must be at least 8 characters long")

		obj.Value("code").Number().Equal(400)
	}
}

func registerInvalidEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "invalid-email.com",
				"password":  "Str0ngP4$$",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().Equal("validation")
		obj.Value("message").String().Equal("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("(email) must be a valid email address: invalid-email.com")

		obj.Value("code").Number().Equal(400)
	}
}

func registerBigEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@email.com",
				"password":  "B1g&mailMan",
			}).
			Expect().
			Status(http.StatusInternalServerError).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("error registering user")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("database-error: pq: value too long for type character varying(255)")

		obj.Value("code").Number().Equal(500)
	}
}

func registerBigPassword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email": "bigpassword@mail.com",
				"password":  "1#Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}).
			Expect().
			Status(http.StatusInternalServerError).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("error hashing user password")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("error: bcrypt: password length exceeds 72 bytes")

		obj.Value("code").Number().Equal(500)
	}
}

func accountAlreadyExists(user *rllCtx) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusInternalServerError).
			JSON().Object()

		obj.Value("module").String().Equal("go-auth-test")
		obj.Value("message").String().Equal("error registering user")

		trace := obj.Value("trace").Array()
		trace.Length().Equal(1)
		trace.Element(0).String().Contains("email is already in use")

		obj.Value("code").Number().Equal(500)
	}
}
