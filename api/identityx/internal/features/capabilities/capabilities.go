package capabilities

import (
	"IdentityX/internal/features/capabilities/commands"
	"IdentityX/internal/features/capabilities/handlers"
	"IdentityX/internal/features/capabilities/queries"
	"IdentityX/internal/features/capabilities/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
