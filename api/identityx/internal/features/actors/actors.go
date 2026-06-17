package actors

import (
	"IdentityX/internal/features/actors/handlers"
	"IdentityX/internal/features/actors/queries"
	"IdentityX/internal/features/actors/repos"
)

var NewRepo = repos.NewRepo
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Handlers = handlers.Handlers
