package initialization

import (
	"TrieForms/internal/plataform/database"
	"context"
	"errors"
	"log"
	"os"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	v1 "github.com/authzed/authzed-go/v1"
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
	module := viper.GetString("MODULE")
	if module == "" {
		module = "trie-forms-module"
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

func SetupDB(app *TrieForms, migrationPath string) {
	var err error
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	app.DB = db

	if err := database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	return rdb
}

func SetupCron(db *pgxpool.Pool, app *TrieForms) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create Scheduler: %v", err)
	}

	app.Scheduler = scheduler

	_ = database.NewPGXTxRunner(db)

	go scheduler.Start()
	log.Println("Started the cron Scheduler")
}

func SetupSpiceDB(app *TrieForms) {
	client, err := v1.NewClient(
		viper.GetString("SPICEDB_ADDR"),
		grpcutil.WithInsecureBearerToken(viper.GetString("SPICEDB_TOKEN")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to SpiceDB: %v", err)
	}
	app.AzClient = client

	schema, err := os.ReadFile("./internal/plataform/authz/schema.zed")
	if err != nil {
		log.Fatalf("failed to read SpiceDB schema: %v", err)
	}

	ctx := context.Background()
	_, err = app.AzClient.WriteSchema(ctx, &pb.WriteSchemaRequest{
		Schema: string(schema),
	})
	if err != nil {
		log.Fatalf("failed to write SpiceDB schema: %v", err)
	}

	log.Println("SpiceDB schema ensured")
}

func SetupGoAuth(app *TrieForms) {
	projectID := uuid.MustParse(viper.GetString("GO_AUTH_PROJECT_ID"))
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:      viper.GetString("GOAUTH_URL"),
		APIKey:       viper.GetString("GOAUTH_API_KEY"),
		ProjectID:    projectID,
		Debug:        true,
		SessionCache: NewRedisSessionCache(app.Redis),
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

	app.GaClient = client
}

type RedisSessionCache struct {
	client *redis.Client
}

func NewRedisSessionCache(rdb *redis.Client) *RedisSessionCache {
	return &RedisSessionCache{
		client: rdb,
	}
}

func (r *RedisSessionCache) GetSession(ctx context.Context, id string) ([]byte, error) {
	val, err := r.client.Get(ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *RedisSessionCache) SetSession(ctx context.Context, id string, data []byte, ttl time.Duration) error {
	return r.client.Set(ctx, id, data, ttl).Err()
}

func (r *RedisSessionCache) DeleteSession(ctx context.Context, id string) error {
	return r.client.Del(ctx, id).Err()
}
