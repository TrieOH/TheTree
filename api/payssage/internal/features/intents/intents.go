package intents

import (
	"payssage/internal/features/intents/commands"
	"payssage/internal/features/intents/handlers"
	"payssage/internal/features/intents/queries"
	"payssage/internal/features/intents/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
