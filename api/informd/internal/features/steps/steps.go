package steps

import (
	"Informd/internal/features/steps/commands"
	"Informd/internal/features/steps/handlers"
	"Informd/internal/features/steps/queries"
	"Informd/internal/features/steps/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Command
type Queries = queries.Queries
type Handlers = handlers.Handlers
