package oauth

import (
	"payssage/internal/features/oauth/commands"
	"payssage/internal/features/oauth/handlers"
	"payssage/internal/features/oauth/repos"
)

var NewRepos = repos.NewRepo
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Handlers = handlers.Handlers
