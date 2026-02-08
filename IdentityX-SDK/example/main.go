package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	// 1. Initialize the Client
	projectID := uuid.MustParse("019c3f28-5431-7005-b596-9344ea3c05b3")
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:   "http://localhost:8080",
		APIKey:    "gk_019c3f28-5431-7005-b596-9344ea3c05b3_Q6vbKw1LZkzLrBEGNcuDyBlV49D0XimjLJlX0DRjmVQ",
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 2. Define a Permission with Conditions
	fmt.Println("\n--- Condition Builder Example ---")

	c := goauth.NewCondition()
	// Logic: (status == 'active') AND (requester == owner OR role == 'admin')
	cond := c.And(
		c.Path("resource.status").Eq("active"),
		c.Or(
			c.Path("resource.requester_id").RefEq("resource.owner_id"),
			c.Path("subject.role").Eq("admin"),
		),
	)

	perm, err := client.Permissions.Define().
		Object(goauth.Object("document", "secure_*")).
		Action("read").
		Conditions(cond). // Automatically calls .Build()
		Create(ctx)
	if err != nil {
		handleError("DefinePermission", err)
		return
	}
	fmt.Printf("Created Permission with ID: %s\n", perm.ID)

	// 3. Authorization Check with Resource Data
	fmt.Println("\n--- Authz Check with Resource ---")

	userID := uuid.MustParse("019c3f35-0f2c-7816-83e3-65b578c91adb")

	// Match: status is active and requester is the owner
	allowed, err := client.Authz.Check().
		User(userID).
		Object("document:secure_report_001").
		Action("read").
		WithResource(map[string]any{
			"status":       "active",
			"owner_id":     userID.String(),
			"requester_id": userID.String(),
		}).
		Allowed(ctx)

	if err != nil {
		handleError("AuthzCheck", err)
		return
	}

	if allowed {
		fmt.Println("✅ Access GRANTED (Condition Matched)")
	} else {
		fmt.Println("❌ Access DENIED")
	}

	// 4. Temporal Grace Example
	fmt.Println("\n--- Temporal Grace Example ---")

	graceCond := c.Field("resource.start_time").GraceBefore("15m")

	gracePerm, err := client.Permissions.Define().
		Object("event:checkin").
		Action("execute").
		Conditions(graceCond).
		Create(ctx)
	if err != nil {
		handleError("DefineGracePermission", err)
		return
	}
	fmt.Printf("Created Grace Permission with ID: %s\n", gracePerm.ID)

	// Check if allowed 5 minutes before event
	eventTime := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	allowed, _ = client.Authz.Check().
		User(userID).
		Object("event:checkin").
		Action("execute").
		WithResource(map[string]any{"start_time": eventTime}).
		Allowed(ctx)

	if allowed {
		fmt.Println("✅ Check-in ALLOWED (Within 15m grace before start)")
	} else {
		fmt.Println("❌ Check-in DENIED")
	}
}

func handleError(op string, err error) {
	if fe, ok := fail.As(err); ok {
		fmt.Printf("Error during %s: [%s] %s\n", op, fe.ID, fe.Message)
	} else {
		fmt.Printf("Error during %s: %v\n", op, err)
	}
}
