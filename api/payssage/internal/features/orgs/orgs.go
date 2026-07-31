package orgs

import (
	"payssage/internal/features/orgs/commands"
	"payssage/internal/features/orgs/handlers"
	"payssage/internal/features/orgs/queries"
	"payssage/internal/features/orgs/repos"
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
