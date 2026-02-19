package main

import (
	"context"
	"fmt"
	"log"

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

	// 2. Create a Permission with simple string object/action
	fmt.Println("\n--- Creating Permission ---")

	perm, err := client.Permissions.Create(ctx, "document", "read")
	if err != nil {
		handleError("CreatePermission", err)
		return
	}
	fmt.Printf("Created Permission: ID=%s, Object=%s, Action=%s\n", perm.ID, perm.Object, perm.Action)

	// 3. Create a Role
	fmt.Println("\n--- Creating Role ---")
	role, err := client.Roles.Define("DocumentReader").Create(ctx)
	if err != nil {
		handleError("CreateRole", err)
		return
	}
	fmt.Printf("Created Role: ID=%s, Name=%s\n", role.ID, role.Name)

	// 4. Add Permission to Role
	fmt.Println("\n--- Adding Permission to Role ---")
	err = client.Roles.AddPermission(ctx, role.ID, perm.ID)
	if err != nil {
		handleError("AddPermissionToRole", err)
		return
	}
	fmt.Println("Permission added to role successfully")

	// 5. Assign Role to User
	fmt.Println("\n--- Assigning Role to User ---")
	userID := uuid.MustParse("019c3f35-0f2c-7816-83e3-65b578c91adb")
	err = client.Roles.GiveToUser(ctx, userID, role.ID, nil)
	if err != nil {
		handleError("GiveRoleToUser", err)
		return
	}
	fmt.Println("Role assigned to user successfully")

	// 6. Check Authorization
	fmt.Println("\n--- Authorization Check ---")
	allowed, err := client.Authz.Check().
		User(userID).
		Object("document").
		Action("read").
		Allowed(ctx)

	if err != nil {
		handleError("AuthzCheck", err)
		return
	}

	if allowed {
		fmt.Println("✅ Access GRANTED")
	} else {
		fmt.Println("❌ Access DENIED")
	}

	// 7. Check with different action (should fail)
	fmt.Println("\n--- Authorization Check (Different Action) ---")
	allowed, err = client.Authz.Check().
		User(userID).
		Object("document").
		Action("write").
		Allowed(ctx)

	if err != nil {
		handleError("AuthzCheck", err)
		return
	}

	if allowed {
		fmt.Println("✅ Access GRANTED")
	} else {
		fmt.Println("❌ Access DENIED (Expected - user only has 'read' permission)")
	}

	// 8. Create a scoped permission
	fmt.Println("\n--- Creating Scoped Permission ---")

	// Create a scope first
	scope, err := client.Scopes.Create(ctx, "department", nil)
	if err != nil {
		handleError("CreateScope", err)
		return
	}
	fmt.Printf("Created Scope: ID=%s, Name=%s\n", scope.ID, scope.Name)

	// Create permission for that scope
	scopedPerm, err := client.Permissions.Create(ctx, "report", "view")
	if err != nil {
		handleError("CreateScopedPermission", err)
		return
	}
	fmt.Printf("Created Scoped Permission: ID=%s\n", scopedPerm.ID)

	// Give user the permission within the scope
	err = client.Permissions.GiveDirect(ctx, userID, scopedPerm.ID, &scope.ID)
	if err != nil {
		handleError("GiveDirectPermission", err)
		return
	}
	fmt.Println("Permission granted to user within scope")

	// Check authorization within the scope
	fmt.Println("\n--- Scoped Authorization Check ---")
	allowed, err = client.Authz.Check().
		User(userID).
		Scope(scope.ID).
		Object("report").
		Action("view").
		Allowed(ctx)

	if err != nil {
		handleError("ScopedAuthzCheck", err)
		return
	}

	if allowed {
		fmt.Println("✅ Access GRANTED within scope")
	} else {
		fmt.Println("❌ Access DENIED")
	}

	// 9. List effective permissions
	fmt.Println("\n--- Listing Effective Permissions ---")
	effPerms, err := client.Permissions.GetEffective(ctx, userID, nil)
	if err != nil {
		handleError("GetEffectivePermissions", err)
		return
	}
	fmt.Printf("User has %d effective permissions:\n", len(effPerms))
	for _, p := range effPerms {
		fmt.Printf("  - %s:%s (ID: %s)\n", p.Object, p.Action, p.ID)
	}
}

func handleError(op string, err error) {
	if fe, ok := fail.As(err); ok {
		fmt.Printf("Error during %s: [%s] %s\n", op, fe.ID, fe.Message)
	} else {
		fmt.Printf("Error during %s: %v\n", op, err)
	}
}
