package webhook_deliveries

import (
	"payssage/internal/features/webhook_deliveries/handlers"
	"payssage/internal/features/webhook_deliveries/queries"
	"payssage/internal/features/webhook_deliveries/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Queries = queries.Queries
type Handlers = handlers.Handlers
