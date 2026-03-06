package initialization

import (
	"context"
	"log"
	"time"
	"univents/internal/plataform/database"

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

func SetupGoAuth(app *UniventsApp) {
	projectID := uuid.MustParse(viper.GetString("GO_AUTH_PROJECT_ID"))
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:   viper.GetString("GOAUTH_URL"),
		APIKey:    viper.GetString("GOAUTH_API_KEY"),
		ProjectID: projectID,
		Debug:     false,
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
				EventsPublish,
				EditionsCreate,
				EditionsAnnounce,
				ActivitiesCreate,
				ActivitiesPublish,
				ActivitiesRead,
				ProductsCreate,
				ProductsRead,
				CheckpointsCreate,
				CheckpointsRead,
				TicketsCreate,
				TicketsEdit,
				TicketsRead,
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
)
