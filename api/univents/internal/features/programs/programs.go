package programs

import (
	"univents/internal/features/programs/commands"
	"univents/internal/features/programs/handlers"
	"univents/internal/features/programs/queries"
	"univents/internal/features/programs/repos"
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
