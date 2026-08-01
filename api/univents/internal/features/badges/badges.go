package badges

import (
	"univents/internal/features/badges/commands"
	"univents/internal/features/badges/handlers"
	"univents/internal/features/badges/queries"
	"univents/internal/features/badges/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandler
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Queries = queries.Queries
type Commands = commands.Commands
type Handler = handlers.Handler
