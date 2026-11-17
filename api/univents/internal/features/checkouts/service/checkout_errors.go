package service

import (
	"errors"
	"fmt"
	"lib/database"
	"lib/fielderr"
	"univents/errs"
	"uuid"
)

// outOfStockError is the 409 listing the items that cannot be fully
// satisfied (availability semantics from split 3).
func outOfStockError(itemIDs []uuid.UUID) error {
	fields := make([]fielderr.Field, 0, len(itemIDs))
	for _, id := range itemIDs {
		fields = append(fields, fielderr.Field{Field: "items", Message: "not enough stock for item " + id.String()})
	}
	return fielderr.Wrap(errs.ErrItemsUnavailable, fields...)
}

func itemNotInEditionError(line CheckoutLine) error {
	return fielderr.On(errs.ErrItemNotInEdition, fmt.Sprintf("items[].%s", line.ItemType), line.ItemID.String())
}

// mapItemLoadError turns a not-found from the FOR UPDATE reads into the 400
// item-list contract (unknown item ids are validation, not 404 — the cart
// the buyer submitted is wrong, the endpoint exists).
func (o *Operations) mapItemLoadError(err error, line CheckoutLine) error {
	if errors.Is(err, database.ErrNotFound) {
		return fielderr.On(errs.ErrUnknownItem, fmt.Sprintf("items[].%s", line.ItemType), line.ItemID.String())
	}
	return err
}
