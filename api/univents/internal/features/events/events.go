package events

import (
	"univents/internal/features/events/commands"
	"univents/internal/features/events/handlers"
	"univents/internal/features/events/queries"
	"univents/internal/features/events/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
