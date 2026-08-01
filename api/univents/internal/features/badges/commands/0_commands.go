package commands

import "univents/ports"

type Commands struct {
	repo ports.BadgeTemplateRepo
}

func NewCommands(repo ports.BadgeTemplateRepo) *Commands {
	return &Commands{
		repo: repo,
	}
}
