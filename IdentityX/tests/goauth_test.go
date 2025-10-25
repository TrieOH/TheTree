package testing

import (
  "database/sql"
	"log"
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

type accountContext struct {
	SuccessEmail string `json:"email"`
	SuccessPasword string `json:"password"`
	accessToken string `json:"-"`
	refreshToken string `json:"-"`
  sessionID string `json:"-"`
  sessionJIT string `json:"-"`
}

func TestGoAuthService(t *testing.T) {
	runServer()
	defer Db.Close()

	rllAcc := &accountContext{
		SuccessEmail: "success@mail.com",
		SuccessPasword: "Str0ngP4ass!",
	}

	t.Run("RegisterTests", func(t *testing.T) {
		runRegisterTests(t, rllAcc)
	})

	t.Run("LoginTests", func(t *testing.T) {
    runLoginTests(t, rllAcc)
	})

	t.Run("LogoutTests", func(t *testing.T) {
    runLogoutTests(t, rllAcc)
	})

	pingAcc := &accountContext{
		SuccessEmail: "ping@mail.com",
		SuccessPasword: "Str0ngP4ass!",
	}

	t.Run("PingTests", func(t *testing.T) {
    runPingTests(t, pingAcc)
	})

	sessionAccount := &accountContext{
		SuccessEmail: "session@mail.com",
		SuccessPasword: "Str0ngP4ass!",
	}

	t.Run("SessionTests", func(t *testing.T) {
    runSessionTests(t, sessionAccount)
	})

	refreshAccount := &accountContext{
		SuccessEmail: "refresh@mail.com",
		SuccessPasword: "Str0ngP4ass!",
	}

	t.Run("RefreshTests", func(t *testing.T) {
    runRefreshTests(t, refreshAccount)
	})
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

func runLoginTests(t *testing.T, ctx *accountContext) {
	t.Run("LoginWrongPassword", loginWrongPassword(ctx))
	t.Run("LoginWrongEmail", loginWrongEmail(ctx))
	t.Run("LoginWrongEmailAndPasword", LoginWrongEmailAndPasword())

	t.Run("LoginSuccess", loginSuccess(ctx))
}

func runLogoutTests(t *testing.T, ctx *accountContext) {
	t.Run("LogoutNoTokens", logoutNoTokens())
	t.Run("LogoutNoRefresh", logoutNoRefresh(ctx))
	t.Run("LogoutSuccess", logoutSuccess(ctx))
	t.Run("LoggedOutAlready", loggedOutAlready(ctx))
}

func runPingTests(t *testing.T, ctx *accountContext) {
  t.Run("Ping", ping())
  t.Run("PrivatePingFailure", privatePingFailure())
  t.Run("CreatePingAccount", registerSuccess(ctx))
  t.Run("LoginPingAccount", loginSuccess(ctx))
  t.Run("PrivatePingSuccess", privatePingSuccess(ctx))
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
    if (oldJIT == ctx.sessionJIT) {
      t.Fatal("refresh token JTIs mustn't match between sessions")
    }
	})
	t.Run("AccessTokensMustNotMatch", func(t *testing.T) {
    if (access == ctx.accessToken) {
      t.Fatal("access tokens mustn't match between sessions")
    }
	})
	t.Run("RefreshTokensMustNotMatch", func(t *testing.T) {
    if (refresh == ctx.refreshToken) {
      t.Fatal("refresh Tokens mustn't match between sessions")
    }
	})

}
