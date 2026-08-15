package authn

import (
	"IdentityX/models"

	"github.com/google/uuid"
)

func testActor() models.Actor {
	email := "actor@trieoh.com"
	return models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
}
