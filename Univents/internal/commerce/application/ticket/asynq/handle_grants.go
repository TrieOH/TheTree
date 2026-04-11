package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"univents/internal/commerce/domain"
	"univents/internal/shared/authz"

	"github.com/hibiken/asynq"
)

func (uc *AsynqHandlers) HandleGrantTicketPermissions(ctx context.Context, t *asynq.Task) error {
	var p domain.GrantTicketPermissionsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	log.Printf("[task] granting permissions for %d tickets", len(p.Grants))

	for _, grant := range p.Grants {
		ticket, err := uc.tickets.GetByID(ctx, grant.TicketID)
		if err != nil {
			log.Printf("[task] error getting ticket: " + grant.TicketID.String())
			return err
		}

		permissions, err := uc.tickets.GetPermissions(ctx, ticket.ID)
		if err != nil {
			log.Printf("[task] error getting ticket permissions for ticket: " + ticket.Name)
			return err
		}

		for _, perm := range permissions {
			switch perm.PermissionType {
			case domain.PermissionTypeCheckpoint:
				if perm.CheckpointID != nil {
					checkpoint, err := uc.checkpoints.GetByID(ctx, *perm.CheckpointID)
					if err != nil {
						log.Printf("[task] error getting checkpoint: " + err.Error())
						return err
					}
					if err = authz.GrantPerm(ctx, uc.az, "checkpoint:"+checkpoint.ID.String()+"#access@user:"+grant.UserID.String()); err != nil {
						log.Printf("[task] error giving permission: " + err.Error())
						return err
					}
				}
			case domain.PermissionTypeActivity:
				if perm.ActivityID != nil {
					activity, err := uc.activities.GetByID(ctx, *perm.ActivityID)
					if err != nil {
						log.Printf("[task] error getting activity: " + err.Error())
						return err
					}
					if err = authz.GrantPerm(ctx, uc.az, "activity:"+activity.ID.String()+"#attend@user:"+grant.UserID.String()); err != nil {
						log.Printf("[task] error giving permission: " + err.Error())
						return err
					}
				}
			case domain.PermissionTypeProduct:
				if perm.ProductID != nil {
					product, err := uc.products.GetByID(ctx, *perm.ProductID)
					if err != nil {
						log.Printf("[task] error getting product: " + err.Error())
						return err
					}
					if err = authz.GrantPerm(ctx, uc.az, "product:"+product.ID.String()+"#purchase@user:"+grant.UserID.String()); err != nil {
						log.Printf("[task] error giving permission: " + err.Error())
						return err
					}
				}
			default:
				log.Println("Grant to invalid Permission type "+perm.PermissionType, perm)
			}
		}
	}

	return nil
}
