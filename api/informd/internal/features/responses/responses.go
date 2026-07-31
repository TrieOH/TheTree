package responses

import (
	"Informd/internal/features/responses/commands"
	"Informd/internal/features/responses/handlers"
	"Informd/internal/features/responses/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Commands = commands.Commands
type Handlers = handlers.Handlers
