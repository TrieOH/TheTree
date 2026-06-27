package editions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lib/database"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	events   ports.EventsRepository
	editions ports.EditionsRepository
	logger   *zap.Logger
	tracer   trace.Tracer
	tx       database.TxRunner
}

func NewCommandService(
	events ports.EventsRepository,
	editions ports.EditionsRepository,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		events:   events,
		editions: editions,
		logger:   logger,
		tracer:   tracer,
		tx:       tx,
	}
}

func (uc *CommandService) Create(ctx context.Context, in contracts.CreateEditionSpec) (out *contracts.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Create")
	defer span.End()

	if err = uc.tx.WithinTx(ctx, func(ctx context.Context) error {
		out, err = uc.createInternal(ctx, in)
		return err
	}); err != nil {
		return &contracts.Edition{}, err
	}

	return out, nil
}

func (uc *CommandService) createInternal(ctx context.Context, in contracts.CreateEditionSpec) (out *contracts.Edition, err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.createInternal")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validEdition *contracts.Edition
	validEdition, err = contracts.NewEdition(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var event *contracts.Event
	event, err = uc.events.GetByID(ctx, in.EventID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_editions"),
		authz.Resource("event", event.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *contracts.Edition
	created, err = uc.editions.Create(ctx, validEdition) // FIXME if this fails the scope must be undone (SAGA PATTERN)
	if err != nil {
		return nil, err
	}

	err = uc.events.AddEdition(ctx, validEdition.EventID)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) Announce(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.Announce")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("announce.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("announce"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return err
	}

	if edition.Status != contracts.EditionStatusDraft {
		return errors.New("can't announce editions on statuses different than draft")
	}

	var task *asynq.Task
	opensAt := edition.RegistrationOpensAt
	if opensAt == nil {
		opensAt = new(time.Now())
	}
	task, err = contracts.NewOpenEditionTask(edition.ID, *opensAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = contracts.NewStartEditionTask(edition.ID, edition.StartsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	task, err = contracts.NewFinishEditionTask(edition.ID, edition.EndsAt)
	if err != nil {
		return err
	}
	if _, err = uc.asynq.EnqueueContext(ctx, task); err != nil {
		return err
	}

	if err = uc.editions.Announce(ctx, editionID); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) ConnectPayments(ctx context.Context, triePaymentsCredentialID, editionID uuid.UUID, triePaymentsProvider, publicKey string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.ConnectPayments")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return fmt.Errorf("error getting edition: %w", err)
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("connect_payments"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return err
	}

	if edition.TriePaymentsCredentialID != nil {
		return errors.New("payment account already connected")
	}

	if err = uc.editions.ConnectPaymentsAccount(ctx, editionID, triePaymentsCredentialID, triePaymentsProvider, publicKey); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) DisconnectPayments(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.DisconnectPayments")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("disconnect_payments"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return err
	}

	if edition.TriePaymentsCredentialID == nil {
		return errors.New("payment account already disconnected")
	}

	if err = uc.editions.DisconnectPaymentsAccount(ctx, editionID); err != nil {
		return err
	}

	return nil
}
