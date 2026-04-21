package app

import (
	"TrieForms/internal/platform/database"
	"context"
	"log"
	"os"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/IdentityX-SDK-Go"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupFUN() {
	var module string
	if module = viper.GetString("MODULE"); module == "" {
		module = "FormService"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})
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

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	return rdb
}

func SetupCron(db *pgxpool.Pool) gocron.Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create Scheduler: %v", err)
	}

	//_ = database.NewPGXTxRunner(db)
	// Call Jobs in jobs.go

	go scheduler.Start()
	log.Println("Started the cron Scheduler")
	return scheduler
}

func SetupSpiceDB() *authzed.Client {
	client, err := authzed.NewClient(
		viper.GetString("SPICEDB_ADDR"),
		grpcutil.WithInsecureBearerToken(viper.GetString("SPICEDB_TOKEN")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to SpiceDB: %v", err)
	}

	schema, err := os.ReadFile("./schema.zed")
	if err != nil {
		log.Fatalf("failed to read SpiceDB schema: %v", err)
	}

	ctx := context.Background()
	_, err = client.WriteSchema(ctx, &pb.WriteSchemaRequest{
		Schema: string(schema),
	})
	if err != nil {
		log.Fatalf("failed to write SpiceDB schema: %v", err)
	}

	log.Println("SpiceDB schema ensured")
	return client
}

func SetupGoAuth() *idx.Client {
	projectID := uuid.MustParse(viper.GetString("GO_AUTH_PROJECT_ID"))
	client, err := idx.NewClient(idx.Config{
		BaseURL:   viper.GetString("GOAUTH_URL"),
		APIKey:    viper.GetString("GOAUTH_API_KEY"),
		ProjectID: projectID,
		Debug:     true,
	})
	if err != nil {
		log.Fatalf("Error creating goauth client: %s", err.Error())
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, true)
		if err != nil {
			log.Fatalf("error fetching initial JWKS: %s", err.Error())
		}
	}()
	return client
}
