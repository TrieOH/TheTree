package purchases

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/models"

	"github.com/google/uuid"
)

// Availability returns the stock position of every purchasable item in the
// edition (available = base - reserved; nil base = unlimited). It fans out
// to the three per-type queries so callers can run it inside the checkout
// tx under the item row locks (split 7); item ids never cross tables.
func (repo *Repo) Availability(ctx context.Context, editionID uuid.UUID) ([]models.ItemAvailability, error) {
	ctx, span := telemetry.StartSpan(ctx, "PurchasesRepo.Availability")
	defer span.End()

	var out []models.ItemAvailability

	tickets, err := database.Queries(ctx, repo.q).ListTicketTypeAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range tickets {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeTicket,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	variants, err := database.Queries(ctx, repo.q).ListProductVariantAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range variants {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeProduct,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	occurrences, err := database.Queries(ctx, repo.q).ListProgramOccurrenceAvailability(ctx, editionID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	for _, row := range occurrences {
		out = append(out, models.ItemAvailability{
			ItemType:         models.PurchaseItemTypeProgramOccurrence,
			ItemID:           row.ItemID,
			BaseQuantity:     row.BaseQuantity,
			ReservedQuantity: row.ReservedQuantity,
		})
	}

	return out, nil
}
