package api_keys

import (
	"IdentityX/internal/features/api_keys/commands"
	"IdentityX/internal/features/api_keys/handlers"
	"IdentityX/internal/features/api_keys/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Handlers = handlers.Handlers
