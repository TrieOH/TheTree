package authn

import (
	"IdentityX/internal/features/authn/commands"
	"IdentityX/internal/features/authn/handlers"
	"IdentityX/internal/features/authn/repos"
)

var NewCommands = commands.NewCommands
var NewRepo = repos.NewRepo
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Handlers = handlers.Handlers
