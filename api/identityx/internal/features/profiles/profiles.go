package profiles

import (
	"IdentityX/internal/features/profiles/commands"
	"IdentityX/internal/features/profiles/handlers"
	"IdentityX/internal/features/profiles/queries"
	"IdentityX/internal/features/profiles/repos"
)

var NewRepo = repos.NewProfileRepo
var NewQueries = queries.New
var NewCommands = commands.New
var NewHandlers = handlers.New
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
