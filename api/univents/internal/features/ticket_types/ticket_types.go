package ticket_types

import (
	"univents/internal/features/ticket_types/commands"
	"univents/internal/features/ticket_types/handlers"
	"univents/internal/features/ticket_types/queries"
	"univents/internal/features/ticket_types/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
