package activities

import (
	"univents/internal/features/activities/commands"
	"univents/internal/features/activities/handlers"
	"univents/internal/features/activities/queries"
	"univents/internal/features/activities/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
