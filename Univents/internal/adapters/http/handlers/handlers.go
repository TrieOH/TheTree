package handlers

import (
	"univents/internal/application"
)

type HandlerBundle struct {
	UniventsHandler *UniventsHandler
}

func New(app *application.Application) *HandlerBundle {
	return &HandlerBundle{
		UniventsHandler: NewUniventsHandler(),
	}
}
