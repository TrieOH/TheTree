package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"lib/database"
	"log"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/authzed/authzed-go/v1"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/trace"
)

type AsynqHandlers struct {
	tickets     ports.TicketsRepository
	products    ports.ProductsRepository
	activities  ports.ActivitiesRepository
	checkpoints ports.CheckpointsRepository
	tracer      trace.Tracer
	az          *authzed.Client
	tx          database.TxRunner
}

func NewAsynqService(
	tickets ports.TicketsRepository,
	products ports.ProductsRepository,
	activities ports.ActivitiesRepository,
	checkpoints ports.CheckpointsRepository,
	tracer trace.Tracer,
	az *authzed.Client,
	tx database.TxRunner,
) *AsynqHandlers {
	return &AsynqHandlers{
		tickets:     tickets,
		products:    products,
		activities:  activities,
		checkpoints: checkpoints,
		tracer:      tracer,
		az:          az,
		tx:          tx,
	}
}

func (uc *AsynqHandlers) HandleGrantTicketPermissions(ctx context.Context, t *asynq.Task) error {
	var p contracts.GrantTicketPermissionsPayload
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
			case contracts.PermissionTypeCheckpoint:
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
			case contracts.PermissionTypeActivity:
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
			case contracts.PermissionTypeProduct:
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
