package signatures

import (
	"univents/internal/features/signatures/commands"
	"univents/internal/features/signatures/handlers"
	"univents/internal/features/signatures/queries"
	"univents/internal/features/signatures/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
