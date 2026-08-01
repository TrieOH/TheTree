package app

import (
	"context"
	"lib/errx"
	libriver "lib/river"
	"log/slog"
	"net/http"
	"univents/internal/authz"
	"univents/internal/features/badges"
	badgesHandlers "univents/internal/features/badges/handlers"
	badgesRepos "univents/internal/features/badges/repos"
	"univents/internal/features/certifications"
	certificationsHandlers "univents/internal/features/certifications/handlers"
	certsJobs "univents/internal/features/certifications/jobs"
	certificationsRepos "univents/internal/features/certifications/repos"
	"univents/internal/features/editions"
	editionsHandlers "univents/internal/features/editions/handlers"
	editionsRepos "univents/internal/features/editions/repos"
	"univents/internal/features/events"
	eventsHandlers "univents/internal/features/events/handlers"
	eventsRepos "univents/internal/features/events/repos"
	"univents/internal/features/products"
	productsHandlers "univents/internal/features/products/handlers"
	productsRepos "univents/internal/features/products/repos"
	"univents/internal/features/programs"
	programsHandlers "univents/internal/features/programs/handlers"
	programsRepos "univents/internal/features/programs/repos"
	"univents/internal/features/signatures"
	signaturesHandlers "univents/internal/features/signatures/handlers"
	signaturesRepos "univents/internal/features/signatures/repos"
	"univents/internal/features/ticket_types"
	ticketTypesHandlers "univents/internal/features/ticket_types/handlers"
	ticketTypesRepos "univents/internal/features/ticket_types/repos"
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

type operations struct {
	events      *events.Operations
	editions    *editions.Operations
	ticketTypes *ticket_types.Operations
	products    *products.Operations
	programs    *programs.Operations
	badges      *badges.Operations
	signatures  *signatures.Operations
	certs       *certifications.Operations
}

type middlewares struct {
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

type handlers struct {
	events      *eventsHandlers.Handlers
	editions    *editionsHandlers.Handlers
	ticketTypes *ticketTypesHandlers.Handlers
	products    *productsHandlers.Handlers
	programs    *programsHandlers.Handlers
	badges      *badgesHandlers.Handlers
	signatures  *signaturesHandlers.Handlers
	certs       *certificationsHandlers.Handlers
}

// ── Init methods ──────────────────────────────────────────────────────────

func (app *Univents) initRepos() repos {
	q := sqlc.New(app.db)
	r := repos{
		events:            eventsRepos.NewRepo(q),
		editions:          editionsRepos.NewRepo(q),
		ticketTypes:       ticketTypesRepos.NewRepo(q),
		products:          productsRepos.NewRepo(q),
		programs:          programsRepos.NewRepo(q),
		occurrences:       programsRepos.NewRepo(q),
		badges:            badgesRepos.NewRepo(q),
		signatures:        signaturesRepos.NewRepo(q),
		signatureRequests: signaturesRepos.NewRepo(q),
		certs:             certificationsRepos.NewRepo(q),
	}
	authz.Service = authz.New(r.events)
	return r
}

func (app *Univents) initOperations(r repos) operations {
	return operations{
		events:      events.NewOperations(r.events, app.objStorage, app.idxClient),
		editions:    editions.NewOperations(r.events, r.editions),
		ticketTypes: ticket_types.NewOperations(r.events, r.editions, r.ticketTypes),
		products:    products.NewOperations(r.events, r.editions, r.products),
		programs:    programs.NewOperations(r.events, r.editions, r.programs, r.occurrences),
		badges:      badges.NewOperations(r.badges),
		signatures:  signatures.NewOperations(r.events, r.editions, r.signatures, r.signatureRequests, app.emailClient, app.cfg.HmacSecret),
		certs:       certifications.NewOperations(r.events, r.editions, r.certs, r.programs, app.emailClient),
	}
}

func (app *Univents) initHandlers(ops operations) handlers {
	return handlers{
		events:      eventsHandlers.NewHandlers(ops.events),
		editions:    editionsHandlers.NewHandlers(ops.editions),
		ticketTypes: ticketTypesHandlers.NewHandlers(ops.ticketTypes),
		products:    productsHandlers.NewHandlers(ops.products),
		programs:    programsHandlers.NewHandlers(ops.programs),
		badges:      badgesHandlers.NewHandler(ops.badges),
		signatures:  signaturesHandlers.NewHandlers(ops.signatures),
		certs:       certificationsHandlers.NewHandlers(ops.certs),
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
		libriver.Register(certsJobs.NewGrantCertsWorker(r.certs, r.editions, r.events, app.emailClient)),
		libriver.Register(certsJobs.NewGrantCertsForOccurrenceWorker(r.certs, r.editions, r.events, app.emailClient)),
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
