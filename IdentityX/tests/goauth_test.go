package integration_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	database "GoAuth/internal/db"
	"GoAuth/internal/router"

	resp "github.com/MintzyG/GoResponse/response"
	"github.com/gavv/httpexpect/v2"
	"github.com/spf13/viper"
)

var Port string
var Db *sql.DB

func init() {
	viper.AutomaticEnv()

	Port = viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "go-auth-test",
	})

	var err error
	Db, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(Db, "../migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
}

var serverUrl string

func runServer() {
	r := router.CreateTestRouter(Db)
	server := httptest.NewServer(r)
	serverUrl = server.URL
}

func createExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  serverUrl,
		Reporter: httpexpect.NewAssertReporter(t),
	})
}

func TestGoAuthService(t *testing.T) {
	runServer()
	defer Db.Close()

	t.Run("RegisterTests", func(t *testing.T) {
		runRegisterTests(t)
	})
}

// var ctx createUserContext
type createUserContext struct {
	id string
}

func runRegisterTests(t *testing.T) {
	t.Run("RegisterNoEmail", registerNoEmail())
	t.Run("RegisterInvalidEmail", registerInvalidEmail())
	t.Run("RegisterBigEmail", registerBigEmail())
	t.Run("RegisterNoPassword", registerNoPassword())
	t.Run("RegisterBigPassword", registerBigPassword())
	t.Run("RegisterWeakPasswordLetters", registerWeakPasswordLetters())
	t.Run("RegisterWeakPasswordLettersNumber", registerWeakPasswordLettersNumber())
	t.Run("RegisterWeakPasswordLettersSymbol", registerWeakPasswordLettersSymbol())
	t.Run("RegisterWeakPasswordLettersUppercase", registerWeakPasswordLettersUppercase())
	t.Run("RegisterWeakPasswordLettersSymbolUppercase", registerWeakPasswordLettersSymbolUppercase())
	t.Run("RegisterWeakPasswordLettersNumberUppercase", registerWeakPasswordLettersNumberUppercase())
	t.Run("RegisterWeakPasswordLettersNumberSymbolUppercase", registerWeakPasswordLettersNumberSymbolUppercase())

	success := map[string]interface{}{
		"email": "success@mail.com",
		"password":  "Str0ngP4ass!",
	}

	t.Run("RegisterSuccess", registerSuccess(&success))
	t.Run("AccountAlreadyExists", accountAlreadyExists(&success))
}

// func getUserAndVerify(ctx *createUserContext, expected *map[string]interface{}) func(t *testing.T) {
// 	return func(t *testing.T) {
// 		e := createExpect(t)
//
// 		userObj := e.GET("/users/{id}", ctx.id).
// 			WithHeader("Content-Type", "application/json").
// 			Expect().
// 			Status(http.StatusOK).
// 			JSON().Object().Value("data").Object()
//
// 		userObj.Value("id").NotNull()
// 		userObj.Value("first_name").NotNull()
//
// 		if expected != nil {
// 			for key, val := range *expected {
// 				userObj.Value(key).Equal(val)
// 			}
// 		}
// 	}
// }

func registerSuccess(user *map[string]interface{}) func(t *testing.T) {
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

func accountAlreadyExists(user *map[string]interface{}) func(t *testing.T) {
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
