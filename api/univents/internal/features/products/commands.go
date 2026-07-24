package products

import (
	"context"
	"errors"
	ports2 "univents/ports"

	"lib/database"
	"lib/objectstorage"
	"univents/internal/shared/authz"
	"univents/internal/shared/contracts"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type CommandService struct {
	editions  ports2.EditionRepo
	products  ports.ProductsRepository
	purchases ports.PurchaseRepository
	obj       *objectstorage.Client
	logger    *zap.Logger
	tracer    trace.Tracer
	tx        database.TxRunner
}

func NewCommandService(
	editions ports2.EditionRepo,
	products ports.ProductsRepository,
	purchases ports.PurchaseRepository,
	obj *objectstorage.Client,
	logger *zap.Logger,
	tracer trace.Tracer,
	tx database.TxRunner,
) *CommandService {
	return &CommandService{
		editions:  editions,
		products:  products,
		purchases: purchases,
		obj:       obj,
		logger:    logger,
		tracer:    tracer,
		tx:        tx,
	}
}

func (uc *CommandService) Create(ctx context.Context, in contracts.CreateProductSpec) (out *contracts.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Create")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("create.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	var validProduct *contracts.Product
	validProduct, err = contracts.NewProduct(sub.ID, in)
	if err != nil {
		return nil, err
	}

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, in.EditionID)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("create_products"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return nil, err
	}

	var created *contracts.Product
	created, err = uc.products.Create(ctx, *validProduct)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *CommandService) Publish(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Publish")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("publish.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var product *contracts.Product
	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("publish"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return err
	}

	if product.Status != contracts.ProductStatusDraft {
		return errors.New("can't publish products on statuses different than draft")
	}

	// TODO: ADD ASYNQ TASKS FOR PRODUCT AVAILABILITY?

	if err = uc.products.Publish(ctx, product.ID); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) SetThumbnail(ctx context.Context, id uuid.UUID, url string) (product *contracts.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.SetThumbnail")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("set_thumbnail.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	product, err = uc.products.SetThumbnail(ctx, product.ID, url)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *CommandService) UnsetThumbnail(ctx context.Context, id uuid.UUID) (product *contracts.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.UnsetThumbnail")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("unset_thumbnail.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	product, err = uc.products.UnsetThumbnail(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *CommandService) AddGalleryImage(ctx context.Context, id uuid.UUID, url string) (product *contracts.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.AddGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("add_gallery.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	product, err = uc.products.AddGalleryImage(ctx, product.ID, url)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *CommandService) RemoveGalleryImage(ctx context.Context, id uuid.UUID, url string) (product *contracts.Product, err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.RemoveGalleryImage")
	defer span.End()
	defer func() { span.SetAttributes(attribute.Bool("remove_gallery.success", err == nil)) }()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("edit"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return nil, err
	}

	bucket, key, err := objectstorage.ParseURL(url)
	if err != nil {
		return nil, errx.Invalid("product").SetMessage("invalid image url")
	}

	if err = uc.obj.RemoveObject(ctx, bucket, key); err != nil {
		return nil, errx.Internal("product").SetMessage("failed to delete image from storage: " + err.Error())
	}

	product, err = uc.products.RemoveGalleryImage(ctx, product.ID, url)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *CommandService) Delete(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Delete")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("delete.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var product *contracts.Product
	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("delete"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return err
	}

	var hasPurchase bool
	hasPurchase, err = uc.products.ItemHasCompletedPurchases(ctx, product.ID)
	if err != nil {
		return err
	}

	if hasPurchase {
		return errx.Forbidden("product").SetMessage("cannot delete a product that was already purchased")
	}

	if err = uc.products.Delete(ctx, product.ID); err != nil {
		return err
	}

	return nil
}

func (uc *CommandService) Restore(ctx context.Context, id uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "ProductService.Restore")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("restore.success", err == nil))
	}()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var product *contracts.Product
	product, err = uc.products.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("restore"),
		authz.Resource("product", product.ID.String()),
	); err != nil {
		return err
	}

	if err = uc.products.Restore(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
