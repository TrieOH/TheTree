package app

import (
	"context"
	"lib/errx"
	libriver "lib/river"
	"log/slog"
	"net/http"
	"univents/internal/authz"
	"univents/internal/features/badges"
	"univents/internal/features/certifications"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/signatures"
	"univents/internal/features/ticket_types"
	"univents/internal/sqlc"
	"univents/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"riverqueue.com/riverui"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	events            ports.EventRepo
	editions          ports.EditionRepo
	ticketTypes       ports.TicketTypeRepo
	products          ports.ProductRepo
	programs          ports.ProgramRepo
	occurrences       ports.ProgramOccurrenceRepo
	badges            ports.BadgeTemplateRepo
	signatures        ports.SignatureRepo
	signatureRequests ports.SignatureRequestRepo
	certs             ports.CertificationRepo
}

type queries struct {
	events      *events.Queries
	editions    *editions.Queries
	ticketTypes *ticket_types.Queries
	products    *products.Queries
	programs    *programs.Queries
	badges      *badges.Queries
	signatures  *signatures.Queries
	certs       *certifications.Queries
}

type commands struct {
	events      *events.Commands
	editions    *editions.Commands
	ticketTypes *ticket_types.Commands
	products    *products.Commands
	programs    *programs.Commands
	badges      *badges.Commands
	signatures  *signatures.Commands
	certs       *certifications.Commands
}

type middlewares struct {
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

type handlers struct {
	events      *events.Handlers
	editions    *editions.Handlers
	ticketTypes *ticket_types.Handlers
	products    *products.Handlers
	programs    *programs.Handlers
	badges      *badges.Handler
	signatures  *signatures.Handlers
	certs       *certifications.Handlers
}

// ── Init methods ──────────────────────────────────────────────────────────

// TODO: Split the duplicate repos
func (app *Univents) initRepos() repos {
	q := sqlc.New(app.db)
	r := repos{
		events:            events.NewRepos(q),
		editions:          editions.NewRepos(q),
		ticketTypes:       ticket_types.NewRepos(q),
		products:          products.NewRepos(q),
		programs:          programs.NewRepos(q),
		occurrences:       programs.NewRepos(q),
		badges:            badges.NewRepos(q),
		signatures:        signatures.NewRepos(q),
		signatureRequests: signatures.NewRepos(q),
		certs:             certifications.NewRepos(q),
	}
	authz.Service = authz.New(r.events)
	return r
}

func (app *Univents) initQueries(r repos) queries {
	return queries{
		events:      events.NewQueries(r.events),
		editions:    editions.NewQueries(r.events, r.editions),
		ticketTypes: ticket_types.NewQueries(r.editions, r.ticketTypes),
		products:    products.NewQueries(r.editions, r.products),
		programs:    programs.NewQueries(r.programs, r.occurrences),
		badges:      badges.NewQueries(r.badges),
		signatures:  signatures.NewQueries(r.editions, r.signatures, r.signatureRequests),
		certs:       certifications.NewQueries(r.certs),
	}
}

func (app *Univents) initCommands(r repos) commands {
	return commands{
		events:      events.NewCommands(r.events, app.objStorage, app.idxClient),
		editions:    editions.NewCommands(r.events, r.editions),
		ticketTypes: ticket_types.NewCommands(r.events, r.editions, r.ticketTypes),
		products:    products.NewCommands(r.events, r.editions, r.products),
		programs:    programs.NewCommands(r.events, r.editions, r.programs, r.occurrences),
		badges:      badges.NewCommands(r.badges),
		signatures:  signatures.NewCommands(r.events, r.editions, r.signatures, r.signatureRequests, app.emailClient, app.cfg.HmacSecret),
		certs:       certifications.NewCommands(r.events, r.editions, r.certs, r.programs, app.emailClient),
	}
}

func (app *Univents) initHandlers(q queries, c commands) handlers {
	return handlers{
		events:      events.NewHandlers(c.events, q.events),
		editions:    editions.NewHandlers(c.editions, q.editions),
		ticketTypes: ticket_types.NewHandlers(c.ticketTypes, q.ticketTypes),
		products:    products.NewHandlers(c.products, q.products),
		programs:    programs.NewHandlers(c.programs, q.programs),
		badges:      badges.NewHandlers(c.badges, q.badges),
		signatures:  signatures.NewHandlers(c.signatures, q.signatures),
		certs:       certifications.NewHandlers(c.certs, q.certs),
	}
}

func (app *Univents) initMiddlewares() middlewares {
	var mw middlewares
	authMW := SetupAuthMiddlewares()

	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *Univents) initRiver(ctx context.Context, r repos) (*river.Client[pgx.Tx], *riverui.Handler) {
	libriver.Migrate(ctx, app.db)

	client := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register(certifications.NewGrantCertsWorker(r.certs, r.editions, r.events, app.emailClient)),
		libriver.Register(certifications.NewGrantCertsForOccurrenceWorker(r.certs, r.editions, r.events, app.emailClient)),
	), nil, nil)
	// TODO: schedule GrantCertsForEdition on edition end and GrantCertsForOccurrence on occurrence end

	err := client.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river client")
	}

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:   false,
		Endpoints: riverui.NewEndpoints[pgx.Tx](client, nil),
		Logger:    slog.Default(),
		Prefix:    "/riverui",
	})
	if err != nil {
		errx.Exit(err, "failed to create river ui handler")
	}
	err = riverUIHandler.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river ui handler")
	}

	return client, riverUIHandler
}
