package namespaces

import (
	"Informd/internal/features/namespaces/commands"
	"Informd/internal/features/namespaces/handlers"
	"Informd/internal/features/namespaces/queries"
	"Informd/internal/features/namespaces/repos"
)

var NewRepo = repos.NewRepo
var NewHandler = handlers.NewHandler
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handlers = handlers.Handler
