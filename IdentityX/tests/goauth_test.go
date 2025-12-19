package testing

import (
	"GoAuth/internal/utils"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	database "GoAuth/internal/db"
	"GoAuth/internal/router"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/gavv/httpexpect/v2"
	"github.com/spf13/viper"
)

var Port string
var Db *sql.DB

func init() {
	viper.AutomaticEnv()

	err := utils.LoadEd25519Keys(
		viper.GetString("JWT_PRIVATE_KEY"),
		viper.GetString("JWT_PUBLIC_KEY"),
	)

	if err != nil {
		log.Fatal(err)
	}

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

type accountContext struct {
	SuccessEmail    string `json:"email"`
	SuccessPassword string `json:"password"`
	accessToken     string
	refreshToken    string
	sessionID       string
	sessionJIT      string
}

func TestGoAuthService(t *testing.T) {
	runServer()
	defer func(Db *sql.DB) {
		err := Db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(Db)

	rllAcc := &accountContext{
		SuccessEmail:    "success@mail.com",
		SuccessPassword: "Str0ngP4ass!",
	}

	t.Run("RegisterTests", func(t *testing.T) { runRegisterTests(t, rllAcc) })
	t.Run("LoginTests", func(t *testing.T) { runLoginTests(t, rllAcc) })
	t.Run("LogoutTests", func(t *testing.T) { runLogoutTests(t, rllAcc) })

	pingAcc := &accountContext{
		SuccessEmail:    "ping@mail.com",
		SuccessPassword: "Str0ngP4ass!",
	}

	t.Run("PingTests", func(t *testing.T) { runPingTests(t, pingAcc) })

	sessionAccount := &accountContext{
		SuccessEmail:    "session@mail.com",
		SuccessPassword: "Str0ngP4ass!",
	}

	t.Run("SessionTests", func(t *testing.T) { runSessionTests(t, sessionAccount) })

	refreshAccount := &accountContext{
		SuccessEmail:    "refresh@mail.com",
		SuccessPassword: "Str0ngP4ass!",
	}

	t.Run("RefreshTests", func(t *testing.T) { runRefreshTests(t, refreshAccount) })
}

func runRegisterTests(t *testing.T, ctx *accountContext) {
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

	t.Run("RegisterSuccess", registerSuccess(ctx))
	t.Run("AccountAlreadyExists", accountAlreadyExists(ctx))
}

// Register Test Cases
func registerSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("Registered user")

		obj.Value("code").Number().IsEqual(201)
	}
}
func registerNoEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "",
				"password": "N0Email#S4d",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(email) is required")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerNoPassword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "nopass@mail.com",
				"password": "",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(password) is required")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLetters() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass1@mail.com",
				"password": "abc",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(3)
		trace.Value(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Value(1).String().Contains("(password) must contain at least one number")
		trace.Value(2).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersNumber() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass2@mail.com",
				"password": "abc3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(2)
		trace.Value(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Value(1).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersSymbol() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass3@mail.com",
				"password": "abc#",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(2)
		trace.Value(0).String().Contains("(password) must contain at least one uppercase letter")
		trace.Value(1).String().Contains("(password) must contain at least one number")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass4@mail.com",
				"password": "Abc",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(2)
		trace.Value(0).String().Contains("(password) must contain at least one number")
		trace.Value(1).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersSymbolUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass5@mail.com",
				"password": "Abc#",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(password) must contain at least one number")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersNumberUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass6@mail.com",
				"password": "Abc3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(password) must contain at least one symbol or punctuation")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerWeakPasswordLettersNumberSymbolUppercase() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "weakpass7@mail.com",
				"password": "Abc#3",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(password) must be at least 8 characters long")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerInvalidEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "invalid-email.com",
				"password": "Str0ngP4$$",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(email) must be a valid email address: invalid-email.com")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerBigEmail() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@email.com",
				"password": "B1g&mailMan",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(email) must be at most 255 characters long:")

		obj.Value("code").Number().IsEqual(400)
	}
}
func registerBigPassword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "bigpassword@mail.com",
				"password": "1#Aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("validation")
		obj.Value("message").String().IsEqual("Validation failed")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("(password) must be at most 72 characters long")

		obj.Value("code").Number().IsEqual(400)
	}
}
func accountAlreadyExists(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/register").
			WithHeader("Content-Type", "application/json").
			WithJSON(user).
			Expect().
			Status(http.StatusConflict).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("error registering user")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("email already in use")

		obj.Value("code").Number().IsEqual(409)
	}
}

func runLoginTests(t *testing.T, ctx *accountContext) {
	t.Run("LoginWrongPassword", loginWrongPassword(ctx))
	t.Run("LoginWrongEmail", loginWrongEmail(ctx))
	t.Run("LoginWrongEmailAndPassword", loginWrongEmailAndPassword())
	t.Run("LoginSuccess", loginSuccess(ctx))
}

// Login Test Cases
func loginWrongPassword(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    user.SuccessEmail,
				"password": "123",
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("invalid email or password")

		obj.Value("code").Number().IsEqual(401)
	}
}
func loginWrongEmail(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "wrong@email.com",
				"password": user.SuccessPassword,
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("invalid email or password")

		obj.Value("code").Number().IsEqual(401)
	}
}
func loginWrongEmailAndPassword() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    "wrong@email.com",
				"password": "Wr0ngP4$$",
			}).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("invalid email or password")

		obj.Value("code").Number().IsEqual(401)
	}
}
func loginSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		r := e.POST("/auth/login").
			WithHeader("Content-Type", "application/json").
			WithJSON(map[string]interface{}{
				"email":    user.SuccessEmail,
				"password": user.SuccessPassword,
			}).
			Expect().
			Status(http.StatusOK)

		obj := r.JSON().Object()
		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("Logged in")
		obj.Value("code").Number().IsEqual(200)

		access := r.Cookie("access_token")
		if access == nil || access.Raw() == nil {
			t.Fatalf("expected access_token cookie, got nil")
		}

		val := access.Value().Raw()
		if val == "" {
			t.Fatalf("access_token cookie value is empty")
		}
		user.accessToken = val

		refresh := r.Cookie("refresh_token")
		if refresh == nil || refresh.Raw() == nil {
			t.Fatalf("expected refresh_token cookie, got nil")
		}

		val = refresh.Value().Raw()
		if val == "" {
			t.Fatalf("refresh_token cookie value is empty")
		}
		user.refreshToken = val
	}
}

func runRefreshTests(t *testing.T, ctx *accountContext) {
	t.Run("CreateRefreshAccount", registerSuccess(ctx))
	t.Run("LoginRefreshAccount", loginSuccess(ctx))
	t.Run("list1Sessions", listXSessions(ctx, 1))
	oldJIT := ctx.sessionJIT
	access := ctx.accessToken
	refresh := ctx.refreshToken
	t.Run("RefreshTokensSuccess", refreshTokensSuccess(ctx))
	t.Run("list1Sessions", listXSessions(ctx, 1))
	t.Run("TokenJTIsMustNotMatch", func(t *testing.T) {
		if oldJIT == ctx.sessionJIT {
			t.Fatal("refresh token JTIs mustn't match between sessions")
		}
	})
	t.Run("AccessTokensMustNotMatch", func(t *testing.T) {
		if access == ctx.accessToken {
			t.Fatal("access tokens mustn't match between sessions")
		}
	})
	t.Run("RefreshTokensMustNotMatch", func(t *testing.T) {
		if refresh == ctx.refreshToken {
			t.Fatal("refresh Tokens mustn't match between sessions")
		}
	})
}

// Refresh Test Cases
func refreshTokensSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		r := e.POST("/auth/refresh").
			WithHeader("Content-Type", "application/json").
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK)

		obj := r.JSON().Object()
		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("Refreshed tokens")

		access := r.Cookie("access_token")
		if access == nil || access.Raw() == nil {
			t.Fatalf("expected access_token cookie, got nil")
		}

		val := access.Value().Raw()
		if val == "" {
			t.Fatalf("access_token cookie value is empty")
		}
		user.accessToken = val

		refresh := r.Cookie("refresh_token")
		if refresh == nil || refresh.Raw() == nil {
			t.Fatalf("expected refresh_token cookie, got nil")
		}

		val = refresh.Value().Raw()
		if val == "" {
			t.Fatalf("refresh_token cookie value is empty")
		}
		user.refreshToken = val

		obj.Value("code").Number().IsEqual(200)
	}
}

func runLogoutTests(t *testing.T, ctx *accountContext) {
	t.Run("LogoutNoTokens", logoutNoTokens())
	t.Run("LogoutNoRefresh", logoutNoRefresh(ctx))
	t.Run("LogoutSuccess", logoutSuccess(ctx))
	t.Run("LoggedOutAlready", loggedOutAlready(ctx))
}

// Logout Test Cases
func logoutNoTokens() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("missing refresh_token cookie")

		obj.Value("code").Number().IsEqual(401)
	}
}
func logoutNoRefresh(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("missing refresh_token cookie")

		obj.Value("code").Number().IsEqual(401)
	}
}
func logoutSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("Logged out")

		obj.Value("code").Number().IsEqual(200)
	}
}
func loggedOutAlready(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/auth/logout").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("user already logged out")

		trace := obj.Value("trace").Array()
		trace.Length().IsEqual(1)
		trace.Value(0).String().Contains("token already blacklisted")

		obj.Value("code").Number().IsEqual(400)
	}
}

func runPingTests(t *testing.T, ctx *accountContext) {
	t.Run("Ping", ping())
	t.Run("PrivatePingFailure", privatePingFailure())
	t.Run("CreatePingAccount", registerSuccess(ctx))
	t.Run("LoginPingAccount", loginSuccess(ctx))
	t.Run("PrivatePingSuccess", privatePingSuccess(ctx))
}

// Ping Test Cases
func ping() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/public").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("pong")

		obj.Value("code").Number().IsEqual(200)
	}
}
func privatePingFailure() func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/private").
			WithHeader("Content-Type", "application/json").
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("AuthMW")
		obj.Value("message").String().IsEqual("missing access_token cookie")

		obj.Value("code").Number().IsEqual(401)
	}
}
func privatePingSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.POST("/ping/private").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().Contains("pong to you")

		obj.Value("code").Number().IsEqual(200)
	}
}

func runSessionTests(t *testing.T, ctx *accountContext) {
	t.Run("CreateSessionAccount", registerSuccess(ctx))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list1Sessions", listXSessions(ctx, 1))
	t.Run("RevokeSessionByIDFail", revokeSessionByIDFail(ctx))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list2Sessions", listXSessions(ctx, 2))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list3Sessions", listXSessions(ctx, 3))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list4Sessions", listXSessions(ctx, 4))
	t.Run("RevokeSessionByIDSuccess", revokeSessionByIDSuccess(ctx))
	t.Run("list3Sessions", listXSessions(ctx, 3))
	t.Run("RevokeOtherSessions", revokeOtherSessions(ctx))
	t.Run("list1Sessions", listXSessions(ctx, 1))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list2Sessions", listXSessions(ctx, 2))
	t.Run("LoginSessionAccount", loginSuccess(ctx))
	t.Run("list3Sessions", listXSessions(ctx, 3))
	t.Run("RevokeAllSessions", revokeAllSessions(ctx))
	t.Run("ListSessionsFail", listSessionsFail(ctx))
}

// Session Test Cases
func listXSessions(user *accountContext, sessionAmount int) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.GET("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")

		data := obj.Value("data").Array()
		data.Length().IsEqual(sessionAmount)
		user.sessionID = data.Value(sessionAmount - 1).Object().Value("session_id").String().Raw()
		user.sessionJIT = data.Value(sessionAmount - 1).Object().Value("token_id").String().Raw()

		obj.Value("code").Number().IsEqual(200)
	}
}
func listSessionsFail(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.GET("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusUnauthorized).
			JSON().Object()

		obj.Value("module").String().IsEqual("AuthMW")
		obj.Value("message").String().IsEqual("refresh token is invalidated")

		obj.Value("code").Number().IsEqual(401)
	}
}
func revokeSessionByIDFail(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.DELETE("/sessions/"+user.sessionID).
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusBadRequest).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("can't revoke a currently active session, please logout instead")
		obj.Value("code").Number().IsEqual(400)
	}
}
func revokeSessionByIDSuccess(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.DELETE("/sessions/"+user.sessionID).
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("revoked session")
		obj.Value("code").Number().IsEqual(200)
	}
}
func revokeOtherSessions(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.DELETE("/sessions/others").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("revoked sessions")
		obj.Value("code").Number().IsEqual(200)
	}
}
func revokeAllSessions(user *accountContext) func(t *testing.T) {
	return func(t *testing.T) {
		e := createExpect(t)

		obj := e.DELETE("/sessions").
			WithHeader("Content-Type", "application/json").
			WithCookie("access_token", user.accessToken).
			WithCookie("refresh_token", user.refreshToken).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		obj.Value("module").String().IsEqual("go-auth-test")
		obj.Value("message").String().IsEqual("revoked sessions")
		obj.Value("code").Number().IsEqual(200)
	}
}
