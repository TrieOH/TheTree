package webhooks

import (
	"payssage/internal/features/webhooks/commands"
	"payssage/internal/features/webhooks/handlers"
)

var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Handlers = handlers.Handlers
