package actors

import (
	"IdentityX/internal/features/actors/commands"
	"IdentityX/internal/features/actors/handlers"
	"IdentityX/internal/features/actors/queries"
	"IdentityX/internal/features/actors/repos"
)

var NewRepo = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
