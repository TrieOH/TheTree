package certifications

import (
	"univents/internal/features/certifications/commands"
	"univents/internal/features/certifications/handlers"
	"univents/internal/features/certifications/jobs"
	"univents/internal/features/certifications/queries"
	"univents/internal/features/certifications/repos"
)

var NewRepos = repos.NewRepo
var NewQueries = queries.NewQueries
var NewCommands = commands.NewCommands
var NewHandlers = handlers.NewHandlers
var RegisterRoutes = handlers.RegisterRoutes
var NewGrantCertsWorker = jobs.NewGrantCertsWorker
var NewGrantCertsForOccurrenceWorker = jobs.NewGrantCertsForOccurrenceWorker

type Repo = repos.Repo
type Queries = queries.Queries
type Commands = commands.Commands
type Handlers = handlers.Handlers
