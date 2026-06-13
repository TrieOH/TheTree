package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"univents/internal/platform/database"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	idx "sdk/identityx"
	"sdk/payssage"
)

type SimpleLogger struct{}

func (l *SimpleLogger) Intercept(_ context.Context, rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func (l *SimpleLogger) InterceptSimple(rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func SetupValidator() *validator.Validate {
	var v = validator.New()
	if err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		vv := fl.Field().String()
		u, err := uuid.Parse(vv)
		if err != nil {
			return false
		}
		return u.Version() == 7
	}); err != nil {
		errx.Must(err, "failed to register uuid7 validator")
	}

	// Use JSON field names for better API responses
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return fld.Name
		}
		return name
	})

	return v
}

func SetupFUN() {
	module := "UniventsAPI"
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	err := fun.AddInterceptor(&SimpleLogger{})
	if err != nil {
		errx.Must(err, "failed to add interceptor")
	}

	v := SetupValidator()
	bind.SetValidator(v)
	fun.SetPathParamFunc(func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
	})
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	return rdb
}

func SetupIdentityX(cfg Config) *idx.Client {
	projectID := cfg.IdxProjectID
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: projectID,
		Debug:     true,
	})
	if err != nil {
		errx.Must(err, "error creating goauth client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, false)
		if err != nil {
			errx.Must(err, "error fetching initial JWKS")
		}
	}()
	return client
}

func SetupPayssage(cfg Config) *payssage.Client {
	paymentsURL := cfg.PayssageURL
	paymentsAPIKey := cfg.PayssageAPIKey
	client := payssage.New(paymentsURL, paymentsAPIKey)
	return client
}

func SetupObjectStorage(cfg Config) *minio.Client {
	client, err := minio.New(cfg.ObjStorageURL, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.ObjStorageAccessKey, cfg.ObjStorageSecretKey, ""),
		Secure: cfg.ObjStorageUseSSL,
	})
	if err != nil {
		errx.Must(err, "failed to create ObjectStorage client")
	}
	log.Println("Connected to ObjectStorage")
	return client
}

func SetupDB(migrationPath string) *pgxpool.Pool {
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	if err = database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
	return db
}

func SetupCron(app *Univents) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create Scheduler: %v", err)
	}
	productHardDeleteJob(scheduler, app.minio, app.db)
	go scheduler.Start()
	log.Println("Started the cron Scheduler")
	return scheduler
}

func SetupSpiceDB(cfg Config) *authzed.Client {
	client, err := authzed.NewClient(
		cfg.SpiceDBAddr,
		grpcutil.WithInsecureBearerToken(cfg.SpiceDBToken),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		errx.Must(err, "failed to connect to SpiceDB")
	}
	schema, err := os.ReadFile("./schema.zed")
	if err != nil {
		errx.Must(err, "failed to read SpiceDB schema")
	}
	ctx := context.Background()
	_, err = client.WriteSchema(ctx, &pb.WriteSchemaRequest{Schema: string(schema)})
	if err != nil {
		errx.Must(err, "failed to write SpiceDB schema")
	}
	log.Println("SpiceDB schema ensured")
	return client
}
