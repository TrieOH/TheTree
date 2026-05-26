package steps

import (
	"Informd/internal/features/steps/commands"
	"Informd/internal/features/steps/handlers"
	"Informd/internal/features/steps/repos"
)

var NewRepo = repos.NewRepo
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Command
type Handlers = handlers.Handlers
