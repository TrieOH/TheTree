package profile_schemas

import (
	"IdentityX/internal/features/profile_schemas/commands"
	"IdentityX/internal/features/profile_schemas/queries"
	"IdentityX/internal/features/profile_schemas/repos"
)

var NewRepo = repos.NewSchemaRepo
var NewQueries = queries.New
var NewCommands = commands.New

type Queries = queries.Queries
type Commands = commands.Commands
