package organizations

import (
	"IdentityX/internal/features/organizations/commands"
	"IdentityX/internal/features/organizations/handlers"
	"IdentityX/internal/features/organizations/queries"
	"IdentityX/internal/features/organizations/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Queries = queries.Queries
type Handlers = handlers.Handlers
