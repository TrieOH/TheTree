package authn

import (
	"IdentityX/internal/features/authn/commands"
	"IdentityX/internal/features/authn/handlers"
)

var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Commands = commands.Commands
type Handlers = handlers.Handlers
