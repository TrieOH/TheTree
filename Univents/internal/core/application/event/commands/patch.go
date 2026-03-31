package commands

import (
	"bytes"
	"context"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"go.opentelemetry.io/otel/attribute"
)

func (uc *CommandService) PatchEvent(ctx context.Context, in domain.PatchEventSpec) (out *domain.Event, warns []string, err error) {
	ctx, span := uc.tracer.Start(ctx, "EventService.PatchEvent")
	defer span.End()
	defer func() {
		span.SetAttributes(attribute.Bool("patch.success", err == nil))
	}()

	span.SetAttributes(attribute.String("event.id", in.ID.String()))

	ga := uc.gaClient

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return nil, nil, err
	}

	var event *domain.Event
	event, err = uc.events.GetByID(ctx, in.ID)
	if err != nil {
		return nil, nil, err
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("events").
		Action("edit").
		Scope(event.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !allowed {
		return nil, nil, errx.Forbidden("event").SetMessage("insufficient permissions")
	}

	if event.Name != in.Name {
		event.Name = in.Name
	}

	if in.Acronym != nil && (event.Acronym == nil || *in.Acronym != *event.Acronym) {
		event.Acronym = in.Acronym
	} else if in.Acronym == nil && event.Acronym != nil {
		event.Acronym = in.Acronym
	}

	if event.Slug != in.Slug {
		event.Slug = in.Slug
	}

	if in.Tagline != nil && (event.Tagline == nil || *in.Tagline != *event.Tagline) {
		event.Tagline = in.Tagline
	} else if in.Tagline == nil && event.Tagline != nil {
		event.Tagline = in.Tagline
	}

	if in.Description != nil && (event.Description == nil || *in.Description != *event.Description) {
		event.Description = in.Description
	} else if in.Description == nil && event.Description != nil {
		event.Description = in.Description
	}

	if event.IsSeries != in.IsSeries {
		if !in.IsSeries && event.EditionsCount > 1 {
			warns = append(warns, "Cannot convert to non-series when multiple editions exist")
		} else {
			event.IsSeries = in.IsSeries
		}
	}

	if in.LogoUrl != nil && (event.LogoUrl == nil || *in.LogoUrl != *event.LogoUrl) {
		event.LogoUrl = in.LogoUrl
	} else if in.LogoUrl == nil && event.LogoUrl != nil {
		event.LogoUrl = in.LogoUrl
	}

	if in.BannerUrl != nil && (event.BannerUrl == nil || *in.BannerUrl != *event.BannerUrl) {
		event.BannerUrl = in.BannerUrl
	} else if in.BannerUrl == nil && event.BannerUrl != nil {
		event.BannerUrl = in.BannerUrl
	}

	if event.HasGallery != in.HasGallery {
		event.HasGallery = in.HasGallery
	}

	if in.ContactEmail != nil && (event.ContactEmail == nil || *in.ContactEmail != *event.ContactEmail) {
		event.ContactEmail = in.ContactEmail
	} else if in.ContactEmail == nil && event.ContactEmail != nil {
		event.ContactEmail = in.ContactEmail
	}

	if in.SocialLinks != nil && (event.SocialLinks == nil || !bytes.Equal(*in.SocialLinks, *event.SocialLinks)) {
		event.SocialLinks = in.SocialLinks
	} else if in.SocialLinks == nil && event.SocialLinks != nil {
		event.SocialLinks = in.SocialLinks
	}

	var patched *domain.Event
	patched, err = uc.events.PatchEvent(ctx, event)
	if err != nil {
		return nil, warns, err
	}

	return patched, warns, nil
}
