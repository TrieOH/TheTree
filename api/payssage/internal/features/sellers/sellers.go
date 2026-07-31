package sellers

import (
	"payssage/internal/features/sellers/handlers"
	"payssage/internal/features/sellers/queries"
	"payssage/internal/features/sellers/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes

type Repo = repos.Repo
type Queries = queries.Queries
type Handlers = handlers.Handlers
