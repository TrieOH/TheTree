package certifications

import (
	"univents/internal/features/certifications/commands"
	"univents/internal/features/certifications/handlers"
	"univents/internal/features/certifications/queries"
	"univents/internal/features/certifications/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
