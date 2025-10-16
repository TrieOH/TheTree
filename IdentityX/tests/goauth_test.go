package testing

import (
	"database/sql"
	"log"
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

type rllCtx struct {
	SuccessEmail string `json:"email"`
	SuccessPasword string `json:"password"`
	accessToken string `json:"-"`
	refreshToken string `json:"-"`
}

func runRegisterTests(t *testing.T) {
	ctx := &rllCtx{
		SuccessEmail: "success@mail.com",
		SuccessPasword: "Str0ngP4ass!",
	}

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
