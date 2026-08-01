package badges

import (
	"univents/ports"
)

type Operations struct {
	repo ports.BadgeTemplateRepo
}

func NewOperations(
	repo ports.BadgeTemplateRepo,
) *Operations {
	return &Operations{
		repo: repo,
	}
}
