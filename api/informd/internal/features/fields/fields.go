package fields

import (
	"Informd/internal/features/fields/commands"
	"Informd/internal/features/fields/handlers"
	"Informd/internal/features/fields/queries"
	"Informd/internal/features/fields/repos"
)

var NewRepos = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Command
type Queries = queries.Queries
type Handlers = handlers.Handlers
