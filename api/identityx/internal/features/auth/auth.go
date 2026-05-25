package auth

import (
	"IdentityX/internal/features/auth/commands"
	"IdentityX/internal/features/auth/handlers"
	"IdentityX/internal/features/auth/queries"
	"IdentityX/internal/features/auth/repos"
)

var NewCommands = commands.NewCommandService
var NewQueries = queries.NewQueryService
var NewHandlers = handlers.NewHandler
var NewRepos = repos.NewRepo
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.CommandService
type Queries = queries.QueryService
type Handlers = handlers.Handler
