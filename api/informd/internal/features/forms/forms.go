package forms

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/handler"
	"Informd/internal/features/forms/queries"
	"Informd/internal/features/forms/repos"
)

var NewFormRepo = repos.NewFormRepo
var NewStepRepo = repos.NewStepRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handler.NewHandlers
var RegisterRoutes = handler.RegisterRoutes

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handlers = handler.Handlers
