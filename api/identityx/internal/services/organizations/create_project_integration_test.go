package organizations

import (
	"context"
	"strings"
	"testing"
	"time"

	"IdentityX/internal/authz"
	"IdentityX/internal/keys"
	"IdentityX/internal/repos"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"lib/crypto"
	"lib/database"
	"lib/testdb"

	"github.com/google/uuid"
)

// TestCreateProjectProvisionsKeysRealDB pins the org-project gap fix
// against a real Postgres (testcontainers): an org-created project must
// ship with signing and encryption keys in the same transaction, instead
// of being token-broken until the next boot's catch-up.
func TestCreateProjectProvisionsKeysRealDB(t *testing.T) {
	pool := testdb.Postgres(t, "../../../db/migrations")
	database.SetDefaultRunner(database.NewPGXTxRunner(pool))

	ctx := context.Background()
	q := sqlc.New(pool)
	r := repos.New(q)

	authzSvc := authz.New(r.Organizations, r.Projects, r.PlatformRoles)
	keysMgr := keys.NewManager(r.CryptoKeys, r.Projects, keys.Config{
		KeyLifetime:    168 * time.Hour,
		RefreshTTL:     7 * time.Hour,
		RotateInterval: 1 * time.Hour,
	}, keys.WithKeyGen(func(models.CryptoKeyType) (*crypto.KeyPair, error) {
		return &crypto.KeyPair{Public: "pub", EncryptedPrivate: "enc", Algorithm: "test"}, nil
	}))
	ops := NewOperations(r.Projects, r.Actors, r.Organizations, keysMgr, authzSvc)

	// org owner (platform actor) with a member row, so CheckOrg admits them
	email := "owner@trieoh.com"
	actor, err := r.Actors.Register(ctx, models.Actor{
		AuthMethod: models.PasswordAuthMethod,
		Email:      &email,
		Type:       models.HumanActorType,
	})
	if err != nil {
		t.Fatalf("register owner: %v", err)
	}
	org, err := r.Organizations.Create(ctx, models.Organization{
		OwnerID: actor.ID,
		Name:    "Test Org",
		Slug:    "testorg",
	})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	err = r.Organizations.AddMember(ctx, models.OrganizationMember{
		OrganizationID: org.ID,
		ActorID:        actor.ID,
		Role:           models.OrganizationRoleOwner,
	})
	if err != nil {
		t.Fatalf("add owner member: %v", err)
	}

	identCtx := models.WithIdentity(ctx, &models.Identity{
		Sub: models.Subject{ID: actor.ID, Type: models.HumanActorType},
	})

	slug := randLettersSlug()
	created, err := ops.CreateProject(identCtx, models.CreateOrgProjectInput{
		OrganizationID: org.ID,
		Name:           "Proj",
		BrandSlug:      slug,
	})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	for _, typ := range []models.CryptoKeyType{models.SigningCryptoKeyType, models.EncryptionCryptoKeyType} {
		key, err := r.CryptoKeys.GetActive(ctx, typ, &created.ID)
		if err != nil {
			t.Fatalf("project %s has no active %s key after org creation: %v", created.ID, typ, err)
		}
		if key.ExpiresAt == nil {
			t.Fatalf("%s key has no expiry — the Key-lifecycle module must stamp a lifetime", typ)
		}
	}
}

// randLettersSlug builds a brand slug satisfying the schema constraint
// (^[a-z]+$, length 3–32): pure lowercase letters.
func randLettersSlug() string {
	var sb strings.Builder
	for _, c := range uuid.New() {
		sb.WriteByte('a' + c%26)
	}
	return "proj" + sb.String()[:8]
}
