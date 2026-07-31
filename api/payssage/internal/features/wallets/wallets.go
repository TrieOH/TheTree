package wallets

import (
	"payssage/internal/features/wallets/commands"
	"payssage/internal/features/wallets/handlers"
	"payssage/internal/features/wallets/queries"
	"payssage/internal/features/wallets/repos"
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
