package profiles

import (
	"IdentityX/internal/features/profiles/queries"
	"IdentityX/internal/features/profiles/repos"
)

var NewRepo = repos.NewProfileRepo
var NewQueries = queries.New

type Queries = queries.Queries
