package checkouts

import (
	"context"
	"time"

	"univents/internal/openapi"
)

// GetEditionAttendeeCount is the public `GET /editions/{edition_id}/attendees/count`
// route: how many confirmed (paid) registrations the edition has — the
// storefront's "N people already registered" number. No identity required,
// like the edition browsing and store stock reads. Unknown editions are
// NOT_FOUND (the edition check happens in the service).
func (h *Handlers) GetEditionAttendeeCount(ctx context.Context, req openapi.GetEditionAttendeeCountRequestObject) (openapi.GetEditionAttendeeCountResponseObject, error) {
	count, err := h.ops.AttendeeCount(ctx, req.EditionId)
	if err != nil {
		return nil, err
	}
	return openapi.GetEditionAttendeeCount200JSONResponse{
		Code: 200, Data: &openapi.AttendeeCount{Count: count}, Timestamp: time.Now(), Module: module,
	}, nil
}
