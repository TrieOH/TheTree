package checkouts

import (
	"context"

	idx "sdk/identityx"
	"univents/internal/features/checkouts/service"
	"univents/internal/openapi"
	"univents/models"
)

type Handlers struct {
	ops *service.Operations
}

func NewHandlers(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

// CreateEditionCheckout is the money path (split 7): reserve the cart and
// create the Payssage intent in one request. The response mirrors the
// shared Purchase read shape (resume == snapshot == checkout uniform) plus
// the one-time ws_token so the front can open the socket immediately.
func (h *Handlers) CreateEditionCheckout(ctx context.Context, req openapi.CreateEditionCheckoutRequestObject) (openapi.CreateEditionCheckoutResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	res, err := h.ops.Checkout(ctx, req.EditionId, ident.Sub.ID, toCheckoutInput(req.Body))
	if err != nil {
		return nil, err
	}

	return openapi.CreateEditionCheckout201JSONResponse(openapi.CheckoutResult{
		PurchaseId:       res.PurchaseID,
		EditionId:        res.EditionID,
		Status:           openapi.PurchaseStatus(res.Status),
		StatusReason:     res.StatusReason,
		TotalCents:       res.TotalCents,
		Currency:         res.Currency,
		PaymentMethod:    res.PaymentMethod,
		PayssageIntentId: res.PayssageIntentID,
		QrCode:           res.QRCode,
		QrCodeBase64:     res.QRCodeBase64,
		ExpiresAt:        res.ExpiresAt,
		CreatedAt:        &res.CreatedAt,
		Items:            toItems(res.Items),
		WsToken:          res.WsToken,
		WsTokenExpiresAt: &res.WsTokenExpiresAt,
	}), nil
}

// toCheckoutInput maps the generated request body to the service's domain
// input. Validation (ticket attendee, per-type quantity, duplicates,
// payment fields for paid orders) runs in the service — the spec only
// enforces the static request shape.
func toCheckoutInput(body *openapi.CreateCheckoutRequest) service.CheckoutInput {
	paymentMethod := ""
	if body.PaymentMethod != nil {
		paymentMethod = string(*body.PaymentMethod)
	}
	payer := service.Payer{}
	if body.Payer != nil {
		payer = service.Payer{
			Email:                body.Payer.Email,
			IdentificationType:   body.Payer.IdentificationType,
			IdentificationNumber: body.Payer.IdentificationNumber,
		}
	}
	in := service.CheckoutInput{
		PaymentMethod:   paymentMethod,
		CardToken:       body.CardToken,
		PaymentMethodID: body.PaymentMethodId,
		Installments:    body.Installments,
		IssuerID:        body.IssuerId,
		Payer:           payer,
		Items:           make([]service.CheckoutLine, 0, len(body.Items)),
	}
	for _, item := range body.Items {
		line := service.CheckoutLine{
			ItemType: models.PurchaseItemType(item.ItemType),
			ItemID:   item.ItemId,
			Quantity: item.Quantity,
		}
		if item.Attendee != nil {
			line.Attendee = &service.Attendee{
				UserID: item.Attendee.UserId,
				Email:  item.Attendee.Email,
				Name:   item.Attendee.Name,
			}
		}
		in.Items = append(in.Items, line)
	}
	return in
}

// GetCheckout is the resume read: the buyer's source of truth for one
// purchase (owner-scoped — a non-owner gets 404 via the service). The
// response mirrors the split-6 WS snapshot shape so the front treats
// resume == snapshot == checkout uniformly.
func (h *Handlers) GetCheckout(ctx context.Context, req openapi.GetCheckoutRequestObject) (openapi.GetCheckoutResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.ops.Get(ctx, req.PurchaseId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}
	return openapi.GetCheckout200JSONResponse(openapi.Checkout{
		PurchaseId:       res.Purchase.ID,
		EditionId:        res.Purchase.EditionID,
		Status:           openapi.PurchaseStatus(res.Purchase.Status),
		StatusReason:     res.Purchase.StatusReason,
		TotalCents:       res.Purchase.TotalCents,
		Currency:         res.Purchase.Currency,
		PaymentMethod:    res.Purchase.PaymentMethod,
		PayssageIntentId: res.Purchase.PayssageIntentID,
		QrCode:           res.Purchase.QRCode,
		QrCodeBase64:     res.Purchase.QRCodeBase64,
		ExpiresAt:        res.Purchase.ExpiresAt,
		CreatedAt:        &res.Purchase.CreatedAt,
		Items:            toItems(res.Items),
		IntentStatus:     res.IntentStatus,
	}), nil
}

// toItems maps the ledger rows to the wire shape (item identity + server-
// computed price; materialization links stay internal).
func toItems(items []models.PurchaseItem) []openapi.PurchaseItem {
	out := make([]openapi.PurchaseItem, 0, len(items))
	for _, item := range items {
		out = append(out, openapi.PurchaseItem{
			ItemType:       openapi.PurchaseItemType(item.ItemType),
			ItemId:         item.ItemID,
			Quantity:       item.Quantity,
			UnitPriceCents: item.UnitPriceCents,
		})
	}
	return out
}

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
	return openapi.GetEditionAttendeeCount200JSONResponse(openapi.AttendeeCount{Count: count}), nil
}

// GetEditionMyTicket is the "what do I hold" read: the caller's active
// ticket for the edition (pending or confirmed registration + its ticket
// type), or `data: null` when they hold none — the front uses it to show
// upgrade options on the more expensive ticket types. Read-only; the
// upgrade itself is a follow-up.
func (h *Handlers) GetEditionMyTicket(ctx context.Context, req openapi.GetEditionMyTicketRequestObject) (openapi.GetEditionMyTicketResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	res, err := h.ops.MyTicket(ctx, req.EditionId, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	var data *openapi.MyTicket
	if res != nil {
		data = &openapi.MyTicket{
			RegistrationId: res.RegistrationID,
			Status:         openapi.MyTicketStatus(res.Status),
			TicketType:     res.TicketType,
		}
	}
	return openapi.GetEditionMyTicket200JSONResponse(*data), nil
}

// ListEditionPurchases is the organizer orders read (owner/admin): every
// purchase of an edition with items + attendees.
func (h *Handlers) ListEditionPurchases(ctx context.Context, req openapi.ListEditionPurchasesRequestObject) (openapi.ListEditionPurchasesResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := h.ops.ListEditionPurchases(ctx, ident.Sub.ID, req.EditionId)
	if err != nil {
		return nil, err
	}
	data := make([]openapi.EditionPurchase, 0, len(rows))
	for _, row := range rows {
		attendees := make([]openapi.PurchaseAttendee, 0, len(row.Attendees))
		for _, a := range row.Attendees {
			attendees = append(attendees, openapi.PurchaseAttendee{Name: a.Name, Email: a.Email})
		}
		data = append(data, openapi.EditionPurchase{
			PurchaseId:       row.Purchase.ID,
			EditionId:        row.Purchase.EditionID,
			Status:           openapi.PurchaseStatus(row.Purchase.Status),
			StatusReason:     row.Purchase.StatusReason,
			TotalCents:       row.Purchase.TotalCents,
			Currency:         row.Purchase.Currency,
			PaymentMethod:    row.Purchase.PaymentMethod,
			PayerEmail:       row.Purchase.PayerEmail,
			PayssageIntentId: row.Purchase.PayssageIntentID,
			CreatedAt:        &row.Purchase.CreatedAt,
			Items:            toItems(row.Items),
			Attendees:        attendees,
		})
	}
	return openapi.ListEditionPurchases200JSONResponse{Items: data}, nil
}

// RefundPurchase initiates a full refund of an approved purchase
// (owner/admin). The purchase stays approved until the webhook flips it.
func (h *Handlers) RefundPurchase(ctx context.Context, req openapi.RefundPurchaseRequestObject) (openapi.RefundPurchaseResponseObject, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	purchase, err := h.ops.RefundPurchase(ctx, ident.Sub.ID, req.PurchaseId)
	if err != nil {
		return nil, err
	}
	return openapi.RefundPurchase200JSONResponse(openapi.Purchase{
		PurchaseId:       purchase.ID,
		EditionId:        purchase.EditionID,
		Status:           openapi.PurchaseStatus(purchase.Status),
		StatusReason:     purchase.StatusReason,
		TotalCents:       purchase.TotalCents,
		Currency:         purchase.Currency,
		PaymentMethod:    purchase.PaymentMethod,
		PayssageIntentId: purchase.PayssageIntentID,
		PayerEmail:       purchase.PayerEmail,
		QrCode:           purchase.QRCode,
		QrCodeBase64:     purchase.QRCodeBase64,
		ExpiresAt:        purchase.ExpiresAt,
		CreatedAt:        &purchase.CreatedAt,
		Items:            []openapi.PurchaseItem{},
	}), nil
}
