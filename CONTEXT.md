# TrieOH

The words we use for TheTree's products and code, and what each one means. When two words could name the same thing, the one here wins. Why things are built this way lives in `docs/adr/`.

## IdentityX

**Actor**:
Anything that can sign in: a person, a service or a machine.
_Avoid_: user, account (except for a project user, below)

**Organization**:
The top-level tenant. Its members are a member, an admin or an owner.

**Project**:
A workspace inside an organization. API keys and profiles belong to a project.

**Project user**:
An actor that belongs to a project, like someone who signed up to an app built on it. It has no role in the project.

**Project member**:
An actor with a role in the project: member, admin or owner. Someone with a role in the project's organization counts as a project member with the same role.

**Platform role**:
A role on IdentityX itself: support, admin or super admin. It's not a role in any organization or project.

**Access check**:
The one place an API asks "may this caller do this?". An unknown resource answers 404 and a missing role answers 403.

**Scope**:
A limit on which kind of caller may use a route, set with `x-scope` in the spec. Today there's one: `platform-only`, for callers that aren't in a project.

**OAuth provider**:
A login a project turns on, like Google or GitHub.

**Profile**:
A person's public face in IdentityX: name, picture and socials.

**Profile schema**:
The versioned shape a project's profiles follow. When a new version arrives, profiles that fit it move to it, and the rest are flagged for an admin.

**API key**:
A long-lived secret a service uses instead of signing in. It carries the capabilities it was given.

**Action token**:
A single-use link token, like the one in a verify or reset email.

**Signing keys**:
Every project, and the platform, always has signing and encryption keys. They rotate on their own, and new projects get theirs when they're created.

## Univents

**Event**:
Something people attend. Its members are an owner, an admin or staff.

**Edition**:
One run of an event, with its own dates and registration window.

**Program**:
The schedule of an edition. Staff are assigned to programs to mark attendance, and attendance earns certificates.

**Check-in**:
A record that someone arrived at an edition. Any staff can make one, and one is enough for the edition's certificate.

**Ticket type, product, product variant**:
What people can buy at an edition. Each person holds at most one active ticket per edition.

**Purchase**:
An order at an edition. Its items stay reserved until the payment confirms, or expire if it doesn't.
_Avoid_: order

**Donation**:
A gift to an event. The donor picks the amount, at least R$ 1. It's not a purchase and reserves nothing.
_Avoid_: doação in code comments and docs

**Certificate**:
Proof of attendance, issued per edition or per program from a certificate template.
_Avoid_: certification

**Badge**:
A person's badge for an edition, with a QR code that opens their Univents profile. Its design comes from the edition's badge templates when it's shown, so a design change shows up everywhere at once.

**Signature request**:
A request, opened by a link, to sign a document.

Univents never reads IdentityX profiles. The apps fetch names and pictures from IdentityX's public profile route.

## Payssage

**Wallet**:
What owns payments: its intents, sellers and webhooks hang off it. Univents uses one wallet for every event, with a 5% fee.

**Collector**:
The wallet's own provider account, like a MercadoPago account. It receives the wallet's fee.

**Seller**:
A provider account a payment is split to, connected to a wallet through the provider's login. Each Univents event has its own seller.

**Intent**:
One payment attempt. Its status comes from the provider.
_Avoid_: payment, transaction

**Webhook endpoint**:
A tenant's URL that Payssage sends payment events to.

**Webhook delivery**:
One try at sending one event to one endpoint. The delivery row is the truth about whether it arrived.

## Informd

**Form**:
A survey made of steps, each with fields.

**Namespace**:
A group of forms with its own members.

**Response**:
One person's answers to a form.

## Code

**Backend**:
One of the four Go APIs. Each is its own process.
_Avoid_: service, microservice

**Feature**:
One slice of a backend, like `intents`. It's a folder in `internal/features/` with its routes in `spec.yml`, plus `handlers.go`, `repo.go` and a `service/` folder.

**Named error**:
A Go error for something a caller did wrong, like `errs.ErrProjectNotFound`. Code returns it and never picks an HTTP status. One mapper per backend turns it into a problem.

**Problem**:
An error as the caller sees it: an RFC 9457 JSON object, declared in the spec.

**Backend kit**:
`lib/backend`, the code every backend starts from. Shared pieces like Postgres or email are modules, and each backend wires the rest itself.

**Harness**:
`lib/httpserver`, which runs the HTTP server: middleware, telemetry, `/health` and shutdown.

**Job carrier**:
What moves background work. Locally it runs in the same process. On Cloudflare it goes through a queue and comes back to the container at `/_jobs`.

**Edge cache**:
Cloudflare's cache in front of each API. Only routes marked `x-cache: shared` are cached, and a write purges what it changed.

**Gateway**:
The front door of each API's Worker. It answers CORS and runs purges, and isn't cached.

**TS client**:
The TypeScript client generated from each backend's spec. It's plain functions, with no React or Solid in it.

**Adapter**:
The concrete thing behind an interface, like the MercadoPago provider or a SQL repo.

## Releases

**Release tag**:
A tag like `payssage/v0.7.11` or `univents-ui/v0.10.6`. Pushing one runs every check, then deploys to production. Only people push tags. Agents never do. If a release is asked for, an agent gets everything ready and leaves the tag to a person.

**Hotfix tag**:
A release tag ending in `-hotfix.N`. It skips everything but the build and unit tests, for real emergencies only. The next normal release runs the full checks.

**Hooks**:
The local checks in `.husky/`: commit message and scope, a secret scan, lint and tests for what changed, and the branch name on push. Never skip them.

**Artifact**:
What a release builds once and deploys as is: an image for a backend, a Worker version for an app, an npm version for an SDK. Rolling back means deploying the previous one.
