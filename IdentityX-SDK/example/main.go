package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/google/uuid"
	"github.com/MintzyG/fail/v3"
)

func main() {
	ctx := context.Background()

	// 1. Initialize the Client
	projectID := uuid.MustParse("your-project-uuid-here")
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:   "http://localhost:8080",
		APIKey:    "gk_project_secret_here",
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 2. Setup Roles and Permissions
	fmt.Println("\n--- RBAC Setup ---")
	
	// Create a permission directly
	perm, err := client.Permissions.Define().
		Object(goauth.Object("doc", "*")).
		Action(goauth.Action("read")).
		Create(ctx)
	if err != nil {
		handleError("CreatePermission", err)
		return
	}
	fmt.Printf("Created Permission: %s:%s\n", perm.Object, perm.Action)

	// Create a role directly
	role, err := client.Roles.Define("Viewer").
		Description("Can read all documents").
		Create(ctx)
	if err != nil {
		handleError("CreateRole", err)
		return
	}
	fmt.Printf("Created Role: %s\n", role.Name)

	// Attach permission to role
	err = client.Roles.AddPermission(ctx, role.ID, perm.ID)
	if err != nil {
		handleError("AddPermissionToRole", err)
		return
	}
	fmt.Println("Attached permission to role")

	// 3. Authz Check
	fmt.Println("\n--- Authorization Check ---")
	
	userID := uuid.New() 

	allowed, err := client.Authz.Check().
		User(userID).
		Object(goauth.Object("doc", "456")).
		Action(goauth.Action("read")).
		Allowed(ctx)
	
	if err != nil {
		handleError("AuthzCheck", err)
		return
	}

	if allowed {
		fmt.Println("✅ User is ALLOWED")
	} else {
		fmt.Println("❌ User is FORBIDDEN")
	}
}

func handleError(op string, err error) {
	if fe, ok := fail.As(err); ok {
		fmt.Printf("Error during %s: [%s] %s\n", op, fe.ID, fe.Message)
	} else {
		fmt.Printf("Error during %s: %v\n", op, err)
	}
}