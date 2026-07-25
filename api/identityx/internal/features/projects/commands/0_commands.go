package commands

import (
	"IdentityX/ports"
	"lib/errx"
)

type Commands struct {
	keys     ports.CryptoKeysRepo
	projects ports.ProjectRepo
	actors   ports.ActorRepo
}

func NewCommands(
	keys ports.CryptoKeysRepo,
	projects ports.ProjectRepo,
	actors ports.ActorRepo,
) *Commands {
	return errx.MustProvide(&Commands{
		keys:     keys,
		projects: projects,
		actors:   actors,
	})
}
