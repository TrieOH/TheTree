package forms

import (
	"Informd/internal/features/forms/commands"
	"Informd/internal/features/forms/handler"
	"Informd/internal/features/forms/queries"
	"Informd/internal/features/forms/repo"
)

var NewFormRepo = repo.NewFormRepo
var NewStepRepo = repo.NewStepRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandler = handler.NewHandler

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handler = handler.Handler
