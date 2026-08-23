import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Calendar, MapPin, Share2 } from "lucide-react";
import { useMemo } from "react";
import {
  activeEditionQueryOptions,
  pastEditionsQueryOptions,
  upcomingEditionsQueryOptions,
} from "@/features/editions/api";
import { resolveEditionVisuals } from "@/features/editions/model/resolve-edition-visuals";
import { OtherEditionsSection } from "@/features/editions/ui/OtherEditionsSection";
import { ContactSection } from "@/features/events/ui/ContactSection";
import { storeStockQueryOptions } from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import { EventCart } from "@/features/products/ui/EventCart";
import { ProductsSection } from "@/features/products/ui/ProductsSection";
import { ProgramSection } from "@/features/programs/ui/ProgramSection";
import { UserActivitiesContent } from "@/features/programs/ui/UserActivitiesContent";
import {
  allTicketsQueryOptions,
  myTicketQueryOptions,
} from "@/features/tickets/api";
import { TicketsSection } from "@/features/tickets/ui/TicketsSection";
import { formatDateRange } from "@/shared/lib/date";
import { getInitials, handleShare } from "@/shared/lib/share";
import { LocationMap } from "@/widgets/ui/map-embed";

export const Route = createLazyFileRoute("/events/$slug/")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = Route.useNavigate();
  const event = Route.useLoaderData();
  const { isAuthenticated } = useAuth();
  const { data: activeEdition } = useSuspenseQuery(
    activeEditionQueryOptions(event.id),
  );

  const { data: upcomingEditions = [] } = useSuspenseQuery(
    upcomingEditionsQueryOptions(event.id),
  );

  const { data: pastEditions = [] } = useSuspenseQuery(
    pastEditionsQueryOptions(event.id),
  );

  const { data: tickets = [] } = useQuery(
    allTicketsQueryOptions(activeEdition?.id ?? ""),
  );

  const { data: stock = [] } = useQuery(
    storeStockQueryOptions(activeEdition?.id ?? ""),
  );

  const { data: heldTicket = null } = useQuery(
    myTicketQueryOptions(activeEdition?.id ?? "", isAuthenticated),
  );

  const initials = getInitials(event.full_name);
  const visuals = resolveEditionVisuals(
    event,
    activeEdition,
    upcomingEditions,
    pastEditions,
  );
  const stockById = new Map(stock.map((item) => [item.id, item.stock]));
  const mapLocation = useMemo(
    () =>
      activeEdition?.location_name
        ? {
            name: activeEdition.location_name,
            address: activeEdition.location_description ?? "",
          }
        : null,
    [activeEdition?.location_description, activeEdition?.location_name],
  );
  useInventoryStream(activeEdition?.id ?? "");

  return (
    <div className="min-h-screen bg-background pb-24">
      {/* Banner Section */}
      <div className="relative">
        <div className="relative h-40 w-full border-b-4 border-b-accent min-[300px]:h-48 sm:h-52 md:h-64">
          {visuals.banner_url ? (
            <img
              src={visuals.banner_url}
              alt={event.full_name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-linear-to-br from-muted via-primary/25 to-secondary/25" />
          )}
        </div>

        <button
          type="button"
          onClick={() =>
            handleShare?.(
              activeEdition?.name ?? event.full_name,
              window.location.href,
            )
          }
          className="absolute right-4 top-4 z-10 shrink-0 rounded-full border border-border/50 bg-background/80 p-2 shadow-sm backdrop-blur-sm transition-all duration-200 hover:scale-105 hover:bg-background active:scale-95 sm:right-6 sm:top-6 sm:p-2.5"
          aria-label="Compartilhar evento"
        >
          <Share2 className="h-4 w-4 text-foreground min-[300px]:h-5 min-[300px]:w-5" />
        </button>

        <div className="absolute inset-x-0 top-full z-10 flex -translate-y-1/2 justify-center">
          {/* Logo / Initials */}
          <div className="relative shrink-0">
            {visuals.logo_url ? (
              <div className="relative">
                <img
                  src={visuals.logo_url}
                  alt={event.full_name}
                  className="aspect-square h-37.5 w-37.5 rounded-full border-4 border-accent bg-muted object-cover shadow-lg sm:h-40 sm:w-40"
                />
              </div>
            ) : (
              <div className="relative">
                <div className="flex aspect-square h-37.5 w-37.5 items-center justify-center rounded-full border-4 border-accent bg-primary shadow-lg sm:h-40 sm:w-40">
                  <span className="text-lg min-[300px]:text-xl sm:text-2xl md:text-3xl font-bold text-primary-foreground">
                    {initials}
                  </span>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="mx-auto max-w-6xl bg-background px-4 pt-24 sm:px-6 sm:pt-28 md:px-8">
        <div className="relative text-center">
          <h1 className="text-xl font-bold leading-tight text-foreground min-[300px]:text-2xl sm:text-3xl md:text-4xl">
            {event.full_name}
          </h1>

          {activeEdition && (
            <div className="mt-3 flex flex-wrap items-center justify-center gap-x-4 gap-y-2 sm:mt-4">
              <div className="flex items-center gap-1.5 text-muted-foreground">
                <Calendar className="w-3.5 h-3.5 min-[300px]:w-4 min-[300px]:h-4 sm:w-4 sm:h-4" />
                <span className="text-xs min-[300px]:text-sm sm:text-sm font-medium">
                  {formatDateRange(
                    activeEdition.starts_at,
                    activeEdition.ends_at,
                  )}
                </span>
              </div>

              {activeEdition.location_name && (
                <div className="flex items-center gap-1.5 text-muted-foreground">
                  <MapPin className="w-3.5 h-3.5 min-[300px]:w-4 min-[300px]:h-4 sm:w-4 sm:h-4" />
                  <span className="text-xs min-[300px]:text-sm sm:text-sm font-medium">
                    {activeEdition.location_name}
                  </span>
                </div>
              )}
            </div>
          )}

          {event.description && (
            <p className="mx-auto mt-5 md:mx-4! border-l-2 border-primary/50 pl-4 text-left text-xs leading-relaxed text-muted-foreground min-[300px]:text-sm sm:mt-6 sm:text-base">
              {event.description}
            </p>
          )}
        </div>
      </div>
      <main className="flex flex-col justify-center items-center w-full mt-6 px-4 sm:px-8!">
        <div className="max-w-6xl w-full">
          <TicketsSection
            tickets={tickets}
            eventSlug={event.slug}
            editionId={activeEdition?.id}
            stockById={stockById}
            heldTicket={heldTicket}
          />
          <ProgramSection
            editionId={activeEdition?.id}
            eventId={event.id}
            eventSlug={event.slug}
          />
          {isAuthenticated && activeEdition ? (
            <section className="w-full py-10">
              <div className="mb-6 text-center">
                <h2 className="text-3xl font-semibold tracking-tight">
                  Minhas atividades
                </h2>
                <p className="mt-2 text-sm text-muted-foreground">
                  Atividades em que você está inscrito neste evento.
                </p>
              </div>
              <UserActivitiesContent edition={activeEdition} />
            </section>
          ) : null}
          <ProductsSection
            editionId={activeEdition?.id}
            eventSlug={event.slug}
            stockById={stockById}
          />
          <OtherEditionsSection
            editions={[...pastEditions, ...upcomingEditions]}
            currentEditionId={activeEdition?.id}
            eventSlug={event.slug}
            maxDisplay={5}
          />
          {mapLocation && (
            <section className="relative z-0 isolate mt-8 w-full sm:mt-10">
              <div className="mb-4">
                <h2 className="text-xl font-semibold text-foreground sm:text-2xl">
                  Local do evento
                </h2>
                <p className="mt-1 text-sm text-muted-foreground">
                  {mapLocation.name}
                  {mapLocation.address ? ` — ${mapLocation.address}` : ""}
                </p>
              </div>
              <LocationMap
                location={mapLocation}
                height="320px"
                className="overflow-hidden rounded-xl border border-border shadow-sm"
              />
            </section>
          )}
          <ContactSection event={event} />
        </div>
      </main>
      {activeEdition && (
        <EventCart
          editionId={activeEdition.id}
          onCheckout={() =>
            navigate({
              to: "/events/$slug/checkout",
              params: { slug: event.slug },
            })
          }
          onExplore={() =>
            navigate({
              to: "/events/$slug/store",
              params: { slug: event.slug },
              search: { tab: "products" },
            })
          }
        />
      )}
    </div>
  );
}
