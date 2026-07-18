package webhook_endpoints

import (
	"payssage/internal/features/webhook_endpoints/commands"
	"payssage/internal/features/webhook_endpoints/handlers"
	"payssage/internal/features/webhook_endpoints/queries"
	"payssage/internal/features/webhook_endpoints/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
