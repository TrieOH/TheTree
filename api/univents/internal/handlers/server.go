// Package handlers implements the generated StrictServerInterface by
// delegating to one subpackage per feature. Auth, validation, and error
// mapping run in the strict middleware stack (see internal/app); the
// handlers here are pure domain logic + fun envelope construction.
package handlers

import (
	"context"

	"univents/internal/handlers/badges"
	"univents/internal/handlers/certifications"
	"univents/internal/handlers/editions"
	"univents/internal/handlers/events"
	"univents/internal/handlers/products"
	"univents/internal/handlers/programs"
	"univents/internal/handlers/signatures"
	"univents/internal/handlers/ticket_types"
	"univents/internal/openapi"
	"univents/internal/services"
)

// Server implements openapi.StrictServerInterface.
type Server struct {
	events      *events.Handlers
	editions    *editions.Handlers
	ticketTypes *ticket_types.Handlers
	products    *products.Handlers
	programs    *programs.Handlers
	badges      *badges.Handlers
	signatures  *signatures.Handlers
	certs       *certifications.Handlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		events:      events.New(ops.Events),
		editions:    editions.New(ops.Editions),
		ticketTypes: ticket_types.New(ops.TicketTypes),
		products:    products.New(ops.Products),
		programs:    programs.New(ops.Programs),
		badges:      badges.New(ops.Badges),
		signatures:  signatures.New(ops.Signatures),
		certs:       certifications.New(ops.Certs),
	}
}

// ── StrictServerInterface ────────────────────────────────────────────────

func (s *Server) DeleteBadgeTemplate(ctx context.Context, req openapi.DeleteBadgeTemplateRequestObject) (openapi.DeleteBadgeTemplateResponseObject, error) {
	return s.badges.DeleteBadgeTemplate(ctx, req)
}

func (s *Server) GetBadgeTemplate(ctx context.Context, req openapi.GetBadgeTemplateRequestObject) (openapi.GetBadgeTemplateResponseObject, error) {
	return s.badges.GetBadgeTemplate(ctx, req)
}

func (s *Server) ListMyCertifications(ctx context.Context, req openapi.ListMyCertificationsRequestObject) (openapi.ListMyCertificationsResponseObject, error) {
	return s.certs.ListMyCertifications(ctx, req)
}

func (s *Server) DeleteCertificationTemplate(ctx context.Context, req openapi.DeleteCertificationTemplateRequestObject) (openapi.DeleteCertificationTemplateResponseObject, error) {
	return s.certs.DeleteCertificationTemplate(ctx, req)
}

func (s *Server) GetCertificationTemplate(ctx context.Context, req openapi.GetCertificationTemplateRequestObject) (openapi.GetCertificationTemplateResponseObject, error) {
	return s.certs.GetCertificationTemplate(ctx, req)
}

func (s *Server) UpdateCertificationTemplate(ctx context.Context, req openapi.UpdateCertificationTemplateRequestObject) (openapi.UpdateCertificationTemplateResponseObject, error) {
	return s.certs.UpdateCertificationTemplate(ctx, req)
}

func (s *Server) UnlinkCertificationTemplate(ctx context.Context, req openapi.UnlinkCertificationTemplateRequestObject) (openapi.UnlinkCertificationTemplateResponseObject, error) {
	return s.certs.UnlinkCertificationTemplate(ctx, req)
}

func (s *Server) LinkCertificationTemplate(ctx context.Context, req openapi.LinkCertificationTemplateRequestObject) (openapi.LinkCertificationTemplateResponseObject, error) {
	return s.certs.LinkCertificationTemplate(ctx, req)
}

func (s *Server) ListCertificationTemplateLinks(ctx context.Context, req openapi.ListCertificationTemplateLinksRequestObject) (openapi.ListCertificationTemplateLinksResponseObject, error) {
	return s.certs.ListCertificationTemplateLinks(ctx, req)
}

func (s *Server) GetCertification(ctx context.Context, req openapi.GetCertificationRequestObject) (openapi.GetCertificationResponseObject, error) {
	return s.certs.GetCertification(ctx, req)
}

func (s *Server) InvalidateCertification(ctx context.Context, req openapi.InvalidateCertificationRequestObject) (openapi.InvalidateCertificationResponseObject, error) {
	return s.certs.InvalidateCertification(ctx, req)
}

func (s *Server) ListBadgeTemplates(ctx context.Context, req openapi.ListBadgeTemplatesRequestObject) (openapi.ListBadgeTemplatesResponseObject, error) {
	return s.badges.ListBadgeTemplates(ctx, req)
}

func (s *Server) CreateBadgeTemplate(ctx context.Context, req openapi.CreateBadgeTemplateRequestObject) (openapi.CreateBadgeTemplateResponseObject, error) {
	return s.badges.CreateBadgeTemplate(ctx, req)
}

func (s *Server) ListEditionCertifications(ctx context.Context, req openapi.ListEditionCertificationsRequestObject) (openapi.ListEditionCertificationsResponseObject, error) {
	return s.certs.ListEditionCertifications(ctx, req)
}

func (s *Server) ListCertificationEmissionErrors(ctx context.Context, req openapi.ListCertificationEmissionErrorsRequestObject) (openapi.ListCertificationEmissionErrorsResponseObject, error) {
	return s.certs.ListCertificationEmissionErrors(ctx, req)
}

func (s *Server) ListCertificationTemplates(ctx context.Context, req openapi.ListCertificationTemplatesRequestObject) (openapi.ListCertificationTemplatesResponseObject, error) {
	return s.certs.ListCertificationTemplates(ctx, req)
}

func (s *Server) CreateCertificationTemplate(ctx context.Context, req openapi.CreateCertificationTemplateRequestObject) (openapi.CreateCertificationTemplateResponseObject, error) {
	return s.certs.CreateCertificationTemplate(ctx, req)
}

func (s *Server) ListEditionOccurrences(ctx context.Context, req openapi.ListEditionOccurrencesRequestObject) (openapi.ListEditionOccurrencesResponseObject, error) {
	return s.programs.ListEditionOccurrences(ctx, req)
}

func (s *Server) ListEditionProducts(ctx context.Context, req openapi.ListEditionProductsRequestObject) (openapi.ListEditionProductsResponseObject, error) {
	return s.products.ListEditionProducts(ctx, req)
}

func (s *Server) CreateInitialProduct(ctx context.Context, req openapi.CreateInitialProductRequestObject) (openapi.CreateInitialProductResponseObject, error) {
	return s.products.CreateInitialProduct(ctx, req)
}

func (s *Server) GetProductByVendorCode(ctx context.Context, req openapi.GetProductByVendorCodeRequestObject) (openapi.GetProductByVendorCodeResponseObject, error) {
	return s.products.GetProductByVendorCode(ctx, req)
}

func (s *Server) ListEditionPrograms(ctx context.Context, req openapi.ListEditionProgramsRequestObject) (openapi.ListEditionProgramsResponseObject, error) {
	return s.programs.ListEditionPrograms(ctx, req)
}

func (s *Server) CreateProgram(ctx context.Context, req openapi.CreateProgramRequestObject) (openapi.CreateProgramResponseObject, error) {
	return s.programs.CreateProgram(ctx, req)
}

func (s *Server) ListEditionSignatureRequests(ctx context.Context, req openapi.ListEditionSignatureRequestsRequestObject) (openapi.ListEditionSignatureRequestsResponseObject, error) {
	return s.signatures.ListEditionSignatureRequests(ctx, req)
}

func (s *Server) CreateSignatureRequest(ctx context.Context, req openapi.CreateSignatureRequestRequestObject) (openapi.CreateSignatureRequestResponseObject, error) {
	return s.signatures.CreateSignatureRequest(ctx, req)
}

func (s *Server) ListEditionSignatures(ctx context.Context, req openapi.ListEditionSignaturesRequestObject) (openapi.ListEditionSignaturesResponseObject, error) {
	return s.signatures.ListEditionSignatures(ctx, req)
}

func (s *Server) CreateSignature(ctx context.Context, req openapi.CreateSignatureRequestObject) (openapi.CreateSignatureResponseObject, error) {
	return s.signatures.CreateSignature(ctx, req)
}

func (s *Server) ListTicketTypes(ctx context.Context, req openapi.ListTicketTypesRequestObject) (openapi.ListTicketTypesResponseObject, error) {
	return s.ticketTypes.ListTicketTypes(ctx, req)
}

func (s *Server) CreateTicketType(ctx context.Context, req openapi.CreateTicketTypeRequestObject) (openapi.CreateTicketTypeResponseObject, error) {
	return s.ticketTypes.CreateTicketType(ctx, req)
}

func (s *Server) GetVariantByVendorCode(ctx context.Context, req openapi.GetVariantByVendorCodeRequestObject) (openapi.GetVariantByVendorCodeResponseObject, error) {
	return s.products.GetVariantByVendorCode(ctx, req)
}

func (s *Server) ListPublicEvents(ctx context.Context, req openapi.ListPublicEventsRequestObject) (openapi.ListPublicEventsResponseObject, error) {
	return s.events.ListPublicEvents(ctx, req)
}

func (s *Server) CreateEvent(ctx context.Context, req openapi.CreateEventRequestObject) (openapi.CreateEventResponseObject, error) {
	return s.events.CreateEvent(ctx, req)
}

func (s *Server) ListJoinedEvents(ctx context.Context, req openapi.ListJoinedEventsRequestObject) (openapi.ListJoinedEventsResponseObject, error) {
	return s.events.ListJoinedEvents(ctx, req)
}

func (s *Server) ListOwnedEvents(ctx context.Context, req openapi.ListOwnedEventsRequestObject) (openapi.ListOwnedEventsResponseObject, error) {
	return s.events.ListOwnedEvents(ctx, req)
}

func (s *Server) PatchEvent(ctx context.Context, req openapi.PatchEventRequestObject) (openapi.PatchEventResponseObject, error) {
	return s.events.PatchEvent(ctx, req)
}

func (s *Server) DiscontinueEvent(ctx context.Context, req openapi.DiscontinueEventRequestObject) (openapi.DiscontinueEventResponseObject, error) {
	return s.events.DiscontinueEvent(ctx, req)
}

func (s *Server) ListPublicEditions(ctx context.Context, req openapi.ListPublicEditionsRequestObject) (openapi.ListPublicEditionsResponseObject, error) {
	return s.editions.ListPublicEditions(ctx, req)
}

func (s *Server) CreateEdition(ctx context.Context, req openapi.CreateEditionRequestObject) (openapi.CreateEditionResponseObject, error) {
	return s.editions.CreateEdition(ctx, req)
}

func (s *Server) GetActiveEdition(ctx context.Context, req openapi.GetActiveEditionRequestObject) (openapi.GetActiveEditionResponseObject, error) {
	return s.editions.GetActiveEdition(ctx, req)
}

func (s *Server) ListDraftEditions(ctx context.Context, req openapi.ListDraftEditionsRequestObject) (openapi.ListDraftEditionsResponseObject, error) {
	return s.editions.ListDraftEditions(ctx, req)
}

func (s *Server) ListPastEditions(ctx context.Context, req openapi.ListPastEditionsRequestObject) (openapi.ListPastEditionsResponseObject, error) {
	return s.editions.ListPastEditions(ctx, req)
}

func (s *Server) ListUpcomingEditions(ctx context.Context, req openapi.ListUpcomingEditionsRequestObject) (openapi.ListUpcomingEditionsResponseObject, error) {
	return s.editions.ListUpcomingEditions(ctx, req)
}

func (s *Server) PatchEdition(ctx context.Context, req openapi.PatchEditionRequestObject) (openapi.PatchEditionResponseObject, error) {
	return s.editions.PatchEdition(ctx, req)
}

func (s *Server) PublishEdition(ctx context.Context, req openapi.PublishEditionRequestObject) (openapi.PublishEditionResponseObject, error) {
	return s.editions.PublishEdition(ctx, req)
}

func (s *Server) ListEventMembers(ctx context.Context, req openapi.ListEventMembersRequestObject) (openapi.ListEventMembersResponseObject, error) {
	return s.events.ListEventMembers(ctx, req)
}

func (s *Server) AddEventMember(ctx context.Context, req openapi.AddEventMemberRequestObject) (openapi.AddEventMemberResponseObject, error) {
	return s.events.AddEventMember(ctx, req)
}

func (s *Server) RemoveEventMember(ctx context.Context, req openapi.RemoveEventMemberRequestObject) (openapi.RemoveEventMemberResponseObject, error) {
	return s.events.RemoveEventMember(ctx, req)
}

func (s *Server) PublishEvent(ctx context.Context, req openapi.PublishEventRequestObject) (openapi.PublishEventResponseObject, error) {
	return s.events.PublishEvent(ctx, req)
}

func (s *Server) GetEventBySlug(ctx context.Context, req openapi.GetEventBySlugRequestObject) (openapi.GetEventBySlugResponseObject, error) {
	return s.events.GetEventBySlug(ctx, req)
}

func (s *Server) GetEditionBySlug(ctx context.Context, req openapi.GetEditionBySlugRequestObject) (openapi.GetEditionBySlugResponseObject, error) {
	return s.editions.GetEditionBySlug(ctx, req)
}

func (s *Server) DeleteOccurrence(ctx context.Context, req openapi.DeleteOccurrenceRequestObject) (openapi.DeleteOccurrenceResponseObject, error) {
	return s.programs.DeleteOccurrence(ctx, req)
}

func (s *Server) GetOccurrence(ctx context.Context, req openapi.GetOccurrenceRequestObject) (openapi.GetOccurrenceResponseObject, error) {
	return s.programs.GetOccurrence(ctx, req)
}

func (s *Server) PatchOccurrence(ctx context.Context, req openapi.PatchOccurrenceRequestObject) (openapi.PatchOccurrenceResponseObject, error) {
	return s.programs.PatchOccurrence(ctx, req)
}

func (s *Server) DeleteProduct(ctx context.Context, req openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	return s.products.DeleteProduct(ctx, req)
}

func (s *Server) GetProduct(ctx context.Context, req openapi.GetProductRequestObject) (openapi.GetProductResponseObject, error) {
	return s.products.GetProduct(ctx, req)
}

func (s *Server) PatchProduct(ctx context.Context, req openapi.PatchProductRequestObject) (openapi.PatchProductResponseObject, error) {
	return s.products.PatchProduct(ctx, req)
}

func (s *Server) ListProductVariants(ctx context.Context, req openapi.ListProductVariantsRequestObject) (openapi.ListProductVariantsResponseObject, error) {
	return s.products.ListProductVariants(ctx, req)
}

func (s *Server) CreateProductVariant(ctx context.Context, req openapi.CreateProductVariantRequestObject) (openapi.CreateProductVariantResponseObject, error) {
	return s.products.CreateProductVariant(ctx, req)
}

func (s *Server) DeleteProgram(ctx context.Context, req openapi.DeleteProgramRequestObject) (openapi.DeleteProgramResponseObject, error) {
	return s.programs.DeleteProgram(ctx, req)
}

func (s *Server) GetProgram(ctx context.Context, req openapi.GetProgramRequestObject) (openapi.GetProgramResponseObject, error) {
	return s.programs.GetProgram(ctx, req)
}

func (s *Server) PatchProgram(ctx context.Context, req openapi.PatchProgramRequestObject) (openapi.PatchProgramResponseObject, error) {
	return s.programs.PatchProgram(ctx, req)
}

func (s *Server) ListProgramOccurrences(ctx context.Context, req openapi.ListProgramOccurrencesRequestObject) (openapi.ListProgramOccurrencesResponseObject, error) {
	return s.programs.ListProgramOccurrences(ctx, req)
}

func (s *Server) CreateProgramOccurrence(ctx context.Context, req openapi.CreateProgramOccurrenceRequestObject) (openapi.CreateProgramOccurrenceResponseObject, error) {
	return s.programs.CreateProgramOccurrence(ctx, req)
}

func (s *Server) DenySignatureRequest(ctx context.Context, req openapi.DenySignatureRequestRequestObject) (openapi.DenySignatureRequestResponseObject, error) {
	return s.signatures.DenySignatureRequest(ctx, req)
}

func (s *Server) FulfillSignatureRequest(ctx context.Context, req openapi.FulfillSignatureRequestRequestObject) (openapi.FulfillSignatureRequestResponseObject, error) {
	return s.signatures.FulfillSignatureRequest(ctx, req)
}

func (s *Server) GetSignatureRequest(ctx context.Context, req openapi.GetSignatureRequestRequestObject) (openapi.GetSignatureRequestResponseObject, error) {
	return s.signatures.GetSignatureRequest(ctx, req)
}

func (s *Server) CancelSignatureRequest(ctx context.Context, req openapi.CancelSignatureRequestRequestObject) (openapi.CancelSignatureRequestResponseObject, error) {
	return s.signatures.CancelSignatureRequest(ctx, req)
}

func (s *Server) RevokeSignature(ctx context.Context, req openapi.RevokeSignatureRequestObject) (openapi.RevokeSignatureResponseObject, error) {
	return s.signatures.RevokeSignature(ctx, req)
}

func (s *Server) DeleteSignature(ctx context.Context, req openapi.DeleteSignatureRequestObject) (openapi.DeleteSignatureResponseObject, error) {
	return s.signatures.DeleteSignature(ctx, req)
}

func (s *Server) GetSignature(ctx context.Context, req openapi.GetSignatureRequestObject) (openapi.GetSignatureResponseObject, error) {
	return s.signatures.GetSignature(ctx, req)
}

func (s *Server) GetTicketType(ctx context.Context, req openapi.GetTicketTypeRequestObject) (openapi.GetTicketTypeResponseObject, error) {
	return s.ticketTypes.GetTicketType(ctx, req)
}

func (s *Server) PatchTicketType(ctx context.Context, req openapi.PatchTicketTypeRequestObject) (openapi.PatchTicketTypeResponseObject, error) {
	return s.ticketTypes.PatchTicketType(ctx, req)
}

func (s *Server) DeleteProductVariant(ctx context.Context, req openapi.DeleteProductVariantRequestObject) (openapi.DeleteProductVariantResponseObject, error) {
	return s.products.DeleteProductVariant(ctx, req)
}

func (s *Server) PatchProductVariant(ctx context.Context, req openapi.PatchProductVariantRequestObject) (openapi.PatchProductVariantResponseObject, error) {
	return s.products.PatchProductVariant(ctx, req)
}

func (s *Server) VerifyCertification(ctx context.Context, req openapi.VerifyCertificationRequestObject) (openapi.VerifyCertificationResponseObject, error) {
	return s.certs.VerifyCertification(ctx, req)
}
