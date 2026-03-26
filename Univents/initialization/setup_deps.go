package initialization

import (
	"context"
	"errors"
	"log"
	"time"
	"univents/internal/plataform/database"
	"univents/internal/plataform/telemetry"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// SimpleLogger is a custom interceptor
type SimpleLogger struct{}

// Intercept is called for responses sent with a context
func (l *SimpleLogger) Intercept(ctx context.Context, rs *resp.Response, statusCode int) {
	l.InterceptSimple(rs, statusCode)
}

// InterceptSimple is called for responses sent without a context
func (l *SimpleLogger) InterceptSimple(rs *resp.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func SetupFUN() {
	module := viper.GetString("MODULE")
	if module == "" {
		module = "univents-module"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})

	err := resp.AddInterceptor(&SimpleLogger{})
	if err != nil {
		log.Fatalf("Failed to add interceptor: %v", err)
	}
}

func SetupPayments(app *UniventsApp) {
	paymentsURL := viper.GetString("TRIEPAYMENTS_URL")
	if paymentsURL == "" {
		log.Fatal("TRIEPAYMENTS_URl not set")
	}
	paymentsAPIKey := viper.GetString("TRIEPAYMENTS_API_KEY")
	if paymentsAPIKey == "" {
		log.Fatal("TRIEPAYMENTS_API_KEY not set")
	}
	paymentsWebhookSecret := viper.GetString("TRIEPAYMENTS_WEBHOOK_SECRET")
	if paymentsWebhookSecret == "" {
		log.Fatal("TRIEPAYMENTS_WEBHOOK_SECRET not set")
	}
	client := paymentsSDK.New(paymentsURL, paymentsAPIKey)
	app.Payments = client
}

func SetupMinio(app *UniventsApp) {
	endpoint := viper.GetString("MINIO_ENDPOINT")
	accessKey := viper.GetString("MINIO_ACCESS_KEY")
	secretKey := viper.GetString("MINIO_SECRET_KEY")
	useSSL := viper.GetBool("MINIO_USE_SSL")

	if endpoint == "" {
		log.Fatal("MINIO_ENDPOINT not set")
	}
	if accessKey == "" {
		log.Fatal("MINIO_ACCESS_KEY not set")
	}
	if secretKey == "" {
		log.Fatal("MINIO_SECRET_KEY not set")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	app.Minio = client
	log.Println("Connected to MinIO")
}

func SetupDB(app *UniventsApp, migrationPath string) {
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

func SetupCron(db *pgxpool.Pool, app *UniventsApp) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create Scheduler: %v", err)
	}

	app.Scheduler = scheduler

	_ = database.NewPGXTxRunner(db)
	productHardDeleteJob(ctx, app)

	go scheduler.Start()
	log.Println("Started the cron Scheduler")
}

func SetupGoAuth(app *UniventsApp) {
	projectID := uuid.MustParse(viper.GetString("GO_AUTH_PROJECT_ID"))
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:      viper.GetString("GOAUTH_URL"),
		APIKey:       viper.GetString("GOAUTH_API_KEY"),
		ProjectID:    projectID,
		Debug:        false,
		SessionCache: NewRedisSessionCache(app.Redis),
	})
	if err != nil {
		log.Fatalf("Error creating goauth client: %s", err.Error())
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, false)
		if err != nil {
			log.Fatalf("error fetching initial JWKS: %s", err.Error())
		}
	}()

	perms, err := client.Permissions.EnsureExists(context.Background(), []goauth.PermissionDefinition{
		EventsCreate,
		ActivitiesAttend,
		ProductsPurchase,
		CheckpointsAccess,
	})
	for _, p := range perms {
		if p.Created {
			log.Println("Created: " + p.Object + ":" + p.Action)
		}
	}

	roles, err := client.Roles.EnsureExists(context.Background(), []goauth.RoleDefinition{
		{
			Name: "Event Owner",
			Permissions: []goauth.PermissionDefinition{
				AttendanceMark,
				EventsPublish,
				EditionsCreate,
				EditionsRead,
				EditionsAnnounce,
				ActivitiesCreate,
				ActivitiesPublish,
				ActivitiesRead,
				ActivitiesManage,
				ProductsCreate,
				ProductsRead,
				ProductsPublish,
				ProductsDelete,
				ProductsRestore,
				ProductsEdit,
				CheckpointsCreate,
				CheckpointsRead,
				TicketsCreate,
				TicketsEdit,
				TicketsRead,
				PaymentsConnect,
				PaymentsDisconnect,
			},
			Meta: map[string]interface{}{
				"color": "#ef4444",
				"icon":  "Shield",
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
	EventsCreate = goauth.PermissionDefinition{
		Object: "events",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	EventsPublish = goauth.PermissionDefinition{
		Object: "events",
		Action: "publish",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Mail",
		},
	}
	EditionsCreate = goauth.PermissionDefinition{
		Object: "editions",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	EditionsAnnounce = goauth.PermissionDefinition{
		Object: "editions",
		Action: "announce",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Mail",
		},
	}
	EditionsRead = goauth.PermissionDefinition{
		Object: "editions",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
	ActivitiesCreate = goauth.PermissionDefinition{
		Object: "activities",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	ActivitiesPublish = goauth.PermissionDefinition{
		Object: "activities",
		Action: "publish",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Mail",
		},
	}
	ActivitiesRead = goauth.PermissionDefinition{
		Object: "activities",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
	ActivitiesAttend = goauth.PermissionDefinition{
		Object: "activities",
		Action: "attend",
		Meta: map[string]interface{}{
			"color": "#203ee6",
			"icon":  "CircleUserRound",
		},
	}
	ActivitiesManage = goauth.PermissionDefinition{
		Object: "activities",
		Action: "manage",
		Meta: map[string]interface{}{
			"color": "#203ee6",
			"icon":  "Shield",
		},
	}
	ProductsCreate = goauth.PermissionDefinition{
		Object: "products",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	ProductsRead = goauth.PermissionDefinition{
		Object: "products",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
	ProductsPurchase = goauth.PermissionDefinition{
		Object: "products",
		Action: "purchase",
		Meta: map[string]interface{}{
			"color": "#203ee6",
			"icon":  "CircleDollarSign",
		},
	}
	ProductsPublish = goauth.PermissionDefinition{
		Object: "products",
		Action: "publish",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Mail",
		},
	}
	ProductsDelete = goauth.PermissionDefinition{
		Object: "products",
		Action: "delete",
		Meta: map[string]interface{}{
			"color": "#ed1a1a",
			"icon":  "Trash2",
		},
	}
	ProductsRestore = goauth.PermissionDefinition{
		Object: "products",
		Action: "restore",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "RotateCcw",
		},
	}
	ProductsEdit = goauth.PermissionDefinition{
		Object: "products",
		Action: "edit",
		Meta: map[string]interface{}{
			"color": "#6366f1",
			"icon":  "PenLine",
		},
	}
	CheckpointsCreate = goauth.PermissionDefinition{
		Object: "checkpoints",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	CheckpointsRead = goauth.PermissionDefinition{
		Object: "checkpoints",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
	CheckpointsAccess = goauth.PermissionDefinition{
		Object: "checkpoints",
		Action: "access",
		Meta: map[string]interface{}{
			"color": "#203ee6",
			"icon":  "CircleUserRound",
		},
	}
	TicketsCreate = goauth.PermissionDefinition{
		Object: "tickets",
		Action: "create",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "Zap",
		},
	}
	TicketsRead = goauth.PermissionDefinition{
		Object: "tickets",
		Action: "read",
		Meta: map[string]interface{}{
			"color": "#f59e0b",
			"icon":  "Eye",
		},
	}
	TicketsEdit = goauth.PermissionDefinition{
		Object: "tickets",
		Action: "edit",
		Meta: map[string]interface{}{
			"color": "#6366f1",
			"icon":  "PenLine",
		},
	}
	AttendanceMark = goauth.PermissionDefinition{
		Object: "attendance",
		Action: "mark",
		Meta: map[string]interface{}{
			"color": "#028cdb",
			"icon":  "UserCheck",
		},
	}
	PaymentsConnect = goauth.PermissionDefinition{
		Object: "payments",
		Action: "connect",
		Meta: map[string]interface{}{
			"color": "#10b981",
			"icon":  "PlugZap",
		},
	}
	PaymentsDisconnect = goauth.PermissionDefinition{
		Object: "payments",
		Action: "disconnect",
		Meta: map[string]interface{}{
			"color": "#ed1a1a",
			"icon":  "GlobeOff",
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
