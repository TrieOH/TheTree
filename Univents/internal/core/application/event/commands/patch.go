package commands

import (
	"bytes"
	"context"
	"univents/internal/core/domain"
	"univents/internal/plataform/telemetry"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

func (uc *CommandService) PatchEvent(ctx context.Context, in domain.PatchEventSpec) (out *domain.Event, warns []string, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.PatchEvent")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("patch.success", err == nil))
	}()

	// FIXME send to BG Worker and return
	var auditor *domain.AuditBuilder
	defer func() {
		if auditor != nil {
			auditor.Emit()
			audits := auditor.GetAudits()
			if err != nil {
				for i := range audits {
					if audits[i].State == domain.ActionStateUnset {
						audits[i].State = domain.ActionStateFailed
					}
				}
			} else {
				for i := range audits {
					if audits[i].State == domain.ActionStateUnset {
						audits[i].State = domain.ActionStateSucceeded
					}
				}
			}
			ae := uc.events.AppendEventAudits(ctx, audits) // FIXME make this outbox later
			if ae != nil {
				telemetry.Log().Error("failed to insert patch event audits", zap.Error(ae))
			}
		}
	}()

	span.SetAttributes(attribute.String("event.id", in.ID.String()))

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, nil, err
	}

	auditor = domain.StartAudit(in.ID, domain.ActorTypeUnknown, &sub.ID).
		State(domain.ActionStateUnset)

	var event *domain.Event
	event, err = uc.events.GetEventByID(ctx, in.ID)
	if err != nil {
		auditor.Action(string(domain.EventAuditActionEdited)).AddMetadata("reason", "failed to get event")
		return nil, nil, err
	}

	isOwner := sub.ID == event.CreatedBy
	if isOwner {
		auditor.Actor(domain.ActorTypeOwner)
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("edit").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		auditor.Action(string(domain.EventAuditActionEdited)).AddMetadata("reason", "Internal System Error")
		return nil, nil, fail.AsFail(err).System().RecordCtx(ctx)
	}
	if !allowed {
		auditor.Action(string(domain.EventAuditActionEdited)).AddMetadata("reason", "Forbidden")
		return nil, nil, fail.New(errx.AuthzInsufficientPermissions)
	}

	if !isOwner {
		auditor.Actor(domain.ActorTypeAdmin)
	}

	if event.Name != in.Name {
		auditor.Action(string(domain.EventAuditActionNameChanged)).
			AddMetadata("from", event.Name).
			AddMetadata("to", in.Name).
			Emit()
		event.Name = in.Name
	}

	if in.Acronym != nil && (event.Acronym == nil || *in.Acronym != *event.Acronym) {
		auditor.Action(string(domain.EventAuditActionAcronymChanged)).
			AddMetadata("from", event.Acronym).
			AddMetadata("to", in.Acronym).
			Emit()
		event.Acronym = in.Acronym
	} else if in.Acronym == nil && event.Acronym != nil {
		auditor.Action(string(domain.EventAuditActionAcronymChanged)).
			AddMetadata("removed", event.Acronym).
			Emit()
		event.Acronym = in.Acronym
	}

	if event.Slug != in.Slug {
		auditor.Action(string(domain.EventAuditActionSlugChanged)).
			AddMetadata("from", event.Slug).
			AddMetadata("to", in.Slug).
			Emit()
		event.Slug = in.Slug
	}

	if in.Tagline != nil && (event.Tagline == nil || *in.Tagline != *event.Tagline) {
		auditor.Action(string(domain.EventAuditActionTaglineChanged)).
			AddMetadata("from", event.Tagline).
			AddMetadata("to", in.Tagline).
			Emit()
		event.Tagline = in.Tagline
	} else if in.Tagline == nil && event.Tagline != nil {
		auditor.Action(string(domain.EventAuditActionTaglineChanged)).
			AddMetadata("removed", event.Tagline).
			Emit()
		event.Tagline = in.Tagline
	}

	if in.Description != nil && (event.Description == nil || *in.Description != *event.Description) {
		auditor.Action(string(domain.EventAuditActionDescriptionChanged)).
			Emit()
		event.Description = in.Description
	} else if in.Description == nil && event.Description != nil {
		auditor.Action(string(domain.EventAuditActionDescriptionChanged)).
			Emit()
		event.Description = in.Description
	}

	if event.IsSeries != in.IsSeries {
		if !in.IsSeries && event.EditionsCount > 1 {
			warns = append(warns, "Cannot convert to non-series when multiple editions exist")
			auditor.Action(string(domain.EventAuditActionIsSeriesChanged)).
				State(domain.ActionStateFailed).
				AddMetadata("reason", "Cannot disable series with multiple editions").
				AddMetadata("editions_count", event.EditionsCount).
				Emit()
		} else {
			auditor.Action(string(domain.EventAuditActionIsSeriesChanged)).
				AddMetadata("from", event.IsSeries).
				AddMetadata("to", in.IsSeries).
				AddMetadata("editions_count", event.EditionsCount).
				Emit()
			event.IsSeries = in.IsSeries
		}
	}

	if in.LogoUrl != nil && (event.LogoUrl == nil || *in.LogoUrl != *event.LogoUrl) {
		auditor.Action(string(domain.EventAuditActionLogoUpdated)).
			AddMetadata("from", event.LogoUrl).
			AddMetadata("to", in.LogoUrl).
			Emit()
		event.LogoUrl = in.LogoUrl
	} else if in.LogoUrl == nil && event.LogoUrl != nil {
		auditor.Action(string(domain.EventAuditActionLogoUpdated)).
			AddMetadata("removed", event.LogoUrl).
			Emit()
		event.LogoUrl = in.LogoUrl
	}

	if in.BannerUrl != nil && (event.BannerUrl == nil || *in.BannerUrl != *event.BannerUrl) {
		auditor.Action(string(domain.EventAuditActionBannerUpdated)).
			AddMetadata("from", event.BannerUrl).
			AddMetadata("to", in.BannerUrl).
			Emit()
		event.BannerUrl = in.BannerUrl
	} else if in.BannerUrl == nil && event.BannerUrl != nil {
		auditor.Action(string(domain.EventAuditActionBannerUpdated)).
			AddMetadata("removed", event.BannerUrl).
			Emit()
		event.BannerUrl = in.BannerUrl
	}

	if event.HasGallery != in.HasGallery {
		auditor.Action(string(domain.EventAuditActionHasGalleryChanged)).
			AddMetadata("from", event.HasGallery).
			AddMetadata("to", in.HasGallery).
			AddMetadata("gallery_urls_count", len(event.GalleryUrls)).
			Emit()
		event.HasGallery = in.HasGallery
	}

	if in.ContactEmail != nil && (event.ContactEmail == nil || *in.ContactEmail != *event.ContactEmail) {
		auditor.Action(string(domain.EventAuditActionContactUpdated)).
			AddMetadata("from", event.ContactEmail).
			AddMetadata("to", in.ContactEmail).
			Emit()
		event.ContactEmail = in.ContactEmail
	} else if in.ContactEmail == nil && event.ContactEmail != nil {
		auditor.Action(string(domain.EventAuditActionContactUpdated)).
			AddMetadata("removed", event.ContactEmail).
			Emit()
		event.ContactEmail = in.ContactEmail
	}

	if in.SocialLinks != nil && (event.SocialLinks == nil || !bytes.Equal(*in.SocialLinks, *event.SocialLinks)) {
		auditor.Action(string(domain.EventAuditActionSocialLinksUpdated)).
			Emit()
		event.SocialLinks = in.SocialLinks
	} else if in.SocialLinks == nil && event.SocialLinks != nil {
		auditor.Action(string(domain.EventAuditActionSocialLinksUpdated)).
			Emit()
		event.SocialLinks = in.SocialLinks
	}

	var patched *domain.Event
	patched, err = uc.events.PatchEvent(ctx, event)
	if err != nil {
		auditor.Action(string(domain.EventAuditActionEdited)).AddMetadata("reason", err.Error())
		return nil, warns, err
	}

	return patched, warns, nil
}
