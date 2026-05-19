package keys

import (
	"Informd/internal/features/keys/commands"
	"Informd/internal/features/keys/handlers"
	"Informd/internal/features/keys/queries"
	"Informd/internal/features/keys/repos"
)

var NewRepos = repos.NewRepos
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handlers = handlers.Handler
