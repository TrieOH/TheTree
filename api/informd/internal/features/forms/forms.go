package forms

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/handlers"
	"Informd/internal/features/forms/queries"
	"Informd/internal/features/forms/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handlers = handlers.Handlers
