package products

import (
	"univents/internal/features/products/commands"
	"univents/internal/features/products/handlers"
	"univents/internal/features/products/queries"
	"univents/internal/features/products/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
