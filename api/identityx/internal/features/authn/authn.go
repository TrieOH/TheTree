package authn

import (
	"IdentityX/internal/features/authn/commands"
	"IdentityX/internal/features/authn/handlers"
	"IdentityX/internal/features/authn/queries"
	"IdentityX/internal/features/authn/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Queries = queries.Queries
type Handlers = handlers.Handlers
