package initialization

import (
	"TriePayments/internal/plataform/database"
	"context"
	"errors"
	"log"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func SetupFUN() {
	module := viper.GetString("MODULE")
	if module == "" {
		module = "trie-payments-module"
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

func SetupDB(app *TriePayments, migrationPath string) {
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

func SetupCron(db *pgxpool.Pool, app *TriePayments) {
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

func SetupGoAuth(app *TriePayments) {
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

	/*perms, err := client.Permissions.EnsureExists(context.Background(), []goauth.PermissionDefinition{})
	for _, p := range perms {
		if p.Created {
			log.Println("Created: " + p.Object + ":" + p.Action)
		}
	}*/

	roles, err := client.Roles.EnsureExists(context.Background(), []goauth.RoleDefinition{
		{
			Name: "Client",
			Permissions: []goauth.PermissionDefinition{
				ApiKeysCreate,
				ApiKeysRevoke,
				WorkspacesCreate,
				WebhooksCreate,
				WebhooksDelete,
				WebhooksRead,
			},
			Meta: map[string]interface{}{
				"color": "#10b981",
				"icon":  "UserStar",
			},
		},
	})

	for _, r := range roles {
		if r.Created {
			log.Println("Created: " + r.Name)
		}
	}

	app.GaClient = client
}

var (
	ApiKeysCreate = goauth.PermissionDefinition{
		Object: "api_keys",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	ApiKeysRevoke = goauth.PermissionDefinition{
		Object: "api_keys",
		Action: "revoke",
		Meta: map[string]interface{}{
			"color": "#ef4444",
			"icon":  "Trash2",
		},
	}
	WorkspacesCreate = goauth.PermissionDefinition{
		Object: "workspaces",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	WebhooksCreate = goauth.PermissionDefinition{
		Object: "webhooks",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	WebhooksDelete = goauth.PermissionDefinition{
		Object: "webhooks",
		Action: "delete",
		Meta: map[string]interface{}{
			"color": "#ef4444",
			"icon":  "Trash2",
		},
	}
	WebhooksRead = goauth.PermissionDefinition{
		Object: "webhooks",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
)

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
