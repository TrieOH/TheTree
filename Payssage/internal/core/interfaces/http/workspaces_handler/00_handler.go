package workspaces_handler

import (
	"TriePayments/internal/core/application/workspaces/commands"
	"TriePayments/internal/core/application/workspaces/queries"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewWorkspacesHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}
