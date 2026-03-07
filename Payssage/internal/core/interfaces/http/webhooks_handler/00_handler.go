package workspaces_handler

import (
	"TriePayments/internal/core/application/webhooks/commands"
	"TriePayments/internal/core/application/webhooks/queries"
)

type Handler struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewWebhooksHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}
