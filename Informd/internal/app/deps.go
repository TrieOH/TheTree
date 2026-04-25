package app

import (
	"Informd/internal/platform/database"
	"Informd/internal/shared/errx"
	"context"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/MintzyG/FastUtilitiesNet"
	"github.com/MintzyG/FastUtilitiesNet/bind"
	"github.com/TrieOH/IdentityX-SDK-Go"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupValidator() *validator.Validate {
	var v = validator.New()
	if err := v.RegisterValidation("uuid7", func(fl validator.FieldLevel) bool {
		v := fl.Field().String()

		u, err := uuid.Parse(v)
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

func SetupFUN(module string) {
	fun.SetConfig(fun.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	v := SetupValidator()
	bind.SetValidator(v)
}

func SetupDB(migrationPath, dsn string) *pgxpool.Pool {
	db, err := database.WaitForDB(30*time.Second, dsn)
	if err != nil {
		errx.Must(err, "Failed to connect DB")
	}
	if err = database.RunMigrations(db, migrationPath); err != nil {
		errx.Must(err, "Failed migrations")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errx.Must(errx.ValidateConstraintRegistry(ctx, db), "unregistered constraints found")
	return db
}

func SetupRedis(timeout time.Duration, addr, password string, db int) *redis.Client {
	rdb, err := database.WaitForRedis(timeout, addr, password, db)
	if err != nil {
		errx.Must(err, "Failed to connect Redis")
	}
	return rdb
}

func SetupCron(db *pgxpool.Pool) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		errx.Must(err, "Failed to create Scheduler")
	}

	//_ = database.NewPGXTxRunner(db)
	// Call Jobs in jobs.go

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
	_, err = client.WriteSchema(ctx, &pb.WriteSchemaRequest{
		Schema: string(schema),
	})
	if err != nil {
		errx.Must(err, "failed to write SpiceDB schema")
	}

	log.Println("SpiceDB schema ensured")
	return client
}

func SetupIdentityX(cfg Config) *idx.Client {
	client, err := idx.NewClient(idx.Config{
		BaseURL:   cfg.IdxURL,
		APIKey:    cfg.IdxAPIKey,
		ProjectID: cfg.IdxProjectID,
		Debug:     true,
	})
	if err != nil {
		errx.Must(err, "error creating identity_x client")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, true)
		if err != nil {
			errx.Must(err, "error fetching initial JWKS")
		}
	}()
	return client
}
