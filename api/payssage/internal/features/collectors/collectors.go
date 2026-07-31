package collectors

import (
	"payssage/internal/features/collectors/handlers"
	"payssage/internal/features/collectors/queries"
	"payssage/internal/features/collectors/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Queries = queries.Queries
type Handlers = handlers.Handlers
