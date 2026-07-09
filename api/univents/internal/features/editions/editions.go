package editions

import (
	"univents/internal/features/editions/commands"
	"univents/internal/features/editions/handlers"
	"univents/internal/features/editions/queries"
	"univents/internal/features/editions/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
