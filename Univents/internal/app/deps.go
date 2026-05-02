package app

import (
	"context"
	"log"
	"os"
	"time"
	"univents/internal/platform/database"
	"univents/internal/platform/telemetry"

	"git.trieoh.com/TrieOH/IdentityX-SDK-Go"
	"git.trieoh.com/TrieOH/Payssage-SDK-Go"
	"github.com/MintzyG/fun"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SimpleLogger is a custom interceptor
type SimpleLogger struct{}

// Intercept is called for responses sent with a context
func (l *SimpleLogger) Intercept(_ context.Context, rs *fun.Response, statusCode int) {
	l.InterceptSimple(rs, statusCode)
}

// InterceptSimple is called for responses sent without a context
func (l *SimpleLogger) InterceptSimple(rs *fun.Response, statusCode int) {
	if statusCode == 500 {
		telemetry.Log().Info("InternalServerError Response", zap.Any("response", rs))
	}
}

func SetupFUN() {
	var module string
	if module = viper.GetString("MODULE"); module == "" {
		module = "UniventsAPI"
	}

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
		log.Fatalf("Failed to add interceptor: %v", err)
	}
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	return rdb
}

func SetupIdentityX() *idx.Client {
	projectID := uuid.MustParse(viper.GetString("IDENTITY_X_PROJECT_ID"))
	client, err := idx.NewClient(idx.Config{
		BaseURL:   viper.GetString("IDENTITY_X_URL"),
		APIKey:    viper.GetString("IDENTITY_X_API_KEY"),
		ProjectID: projectID,
		Debug:     true,
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
	return client
}

func SetupPayssage() *payssage.Client {
	paymentsURL := viper.GetString("PAYSSAGE_URL")
	if paymentsURL == "" {
		log.Fatal("PAYSSAGE_URl not set")
	}
	paymentsAPIKey := viper.GetString("PAYSSAGE_API_KEY")
	if paymentsAPIKey == "" {
		log.Fatal("PAYSSAGE_API_KEY not set")
	}
	paymentsWebhookSecret := viper.GetString("PAYSSAGE_WEBHOOK_SECRET")
	if paymentsWebhookSecret == "" {
		log.Fatal("PAYSSAGE_WEBHOOK_SECRET not set")
	}
	client := payssage.New(paymentsURL, paymentsAPIKey)
	return client
}

func SetupObjectStorage() *minio.Client {
	endpoint := viper.GetString("OBJECT_STORAGE_ENDPOINT")
	accessKey := viper.GetString("OBJECT_STORAGE_ACCESS_KEY")
	secretKey := viper.GetString("OBJECT_STORAGE_SECRET_KEY")
	useSSL := viper.GetBool("OBJECT_STORAGE_USE_SSL")

	if endpoint == "" {
		log.Fatal("OBJECT_STORAGE_ENDPOINT not set")
	}
	if accessKey == "" {
		log.Fatal("OBJECT_STORAGE_ACCESS_KEY not set")
	}
	if secretKey == "" {
		log.Fatal("OBJECT_STORAGE_SECRET_KEY not set")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}
	log.Println("Connected to MinIO")
	return client
}

func SetupDB(migrationPath string) *pgxpool.Pool {
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	if err := database.RunMigrations(db, migrationPath); err != nil {
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
