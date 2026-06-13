package projects

import (
	"IdentityX/internal/features/projects/commands"
	"IdentityX/internal/features/projects/handlers"
	"IdentityX/internal/features/projects/queries"
	"IdentityX/internal/features/projects/repos"
)

var NewRepos = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Queries = queries.Queries
type Handlers = handlers.Handlers
