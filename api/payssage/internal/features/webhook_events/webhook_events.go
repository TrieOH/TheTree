package webhook_events

import (
	"payssage/internal/features/webhook_events/handlers"
	"payssage/internal/features/webhook_events/queries"
	"payssage/internal/features/webhook_events/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Handlers = handlers.Handlers
