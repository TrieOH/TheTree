import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createLazyFileRoute } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Calendar, MapPin, Share2 } from "lucide-react";
import {
  activeEditionQueryOptions,
  pastEditionsQueryOptions,
  upcomingEditionsQueryOptions,
} from "@/features/editions/api";
import { OtherEditionsSection } from "@/features/editions/ui/OtherEditionsSection";
import { ContactSection } from "@/features/events/ui/ContactSection";
import { storeStockQueryOptions } from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import { EventCart } from "@/features/products/ui/EventCart";
import { ProductsSection } from "@/features/products/ui/ProductsSection";
import { ProgramSection } from "@/features/programs/ui/ProgramSection";
import {
  allTicketsQueryOptions,
  myTicketQueryOptions,
} from "@/features/tickets/api";
import { TicketsSection } from "@/features/tickets/ui/TicketsSection";
import { formatDateRange } from "@/shared/lib/date";
import { getInitials, handleShare } from "@/shared/lib/share";

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
  const stockById = new Map(stock.map((item) => [item.id, item.stock]));
  useInventoryStream(activeEdition?.id ?? "");

  return (
    <div className="min-h-screen bg-background pb-24">
      {/* Banner Section */}
      <div className="relative">
        <div className="relative w-full h-40 min-[300px]:h-48 sm:h-52 md:h-64">
          {event.banner_url ? (
            <img
              src={event.banner_url}
              alt={event.full_name}
              className="w-full h-full object-cover"
            />
          ) : (
            <div className="w-full h-full bg-linear-to-br from-muted via-primary/25 to-secondary/25" />
          )}

          {/* Top gradient overlay */}
          <div className="absolute inset-x-0 top-0 h-32 bg-linear-to-b from-muted/20 via-primary/10 to-transparent" />

          {/* Middle gradient overlay */}
          <div className="absolute inset-x-0 bottom-0 h-2/3 bg-linear-to-t from-background/80 via-secondary/20 to-transparent" />

          {/* Bottom fade to background */}
          <div className="absolute inset-x-0 bottom-0 h-24 bg-linear-to-t from-background to-transparent" />
        </div>

        {/* Content overlay - positioned absolutely over the banner */}
        <div className="absolute inset-x-0 bottom-0 px-4 sm:px-6 md:px-8 pb-4 sm:pb-6">
          <div className="flex items-end gap-3 sm:gap-4">
            {/* Logo / Initials */}
            <div className="relative shrink-0">
              {event.logo_url ? (
                <div className="relative">
                  <img
                    src={event.logo_url}
                    alt={event.full_name}
                    className="w-16 h-16 min-[300px]:w-20 min-[300px]:h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-xl object-cover border-2 border-background shadow-lg"
                  />
                  {/* Blur backdrop for logo */}
                  <div className="absolute -inset-1 rounded-xl bg-primary/10 blur-md -z-10" />
                </div>
              ) : (
                <div className="relative">
                  <div className="w-16 h-16 min-[300px]:w-20 min-[300px]:h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-xl bg-primary flex items-center justify-center border-2 border-background shadow-lg">
                    <span className="text-lg min-[300px]:text-xl sm:text-2xl md:text-3xl font-bold text-primary-foreground">
                      {initials}
                    </span>
                  </div>
                  {/* Blur backdrop for initials */}
                  <div className="absolute -inset-1 rounded-xl bg-primary/20 blur-md -z-10" />
                </div>
              )}
            </div>

            {/* Text content */}
            <div className="flex-1 min-w-0 pb-1">
              {/* Tagline / Badge */}
              {(activeEdition?.tagline ?? event.acronym) && (
                <div className="flex items-center gap-2 mb-1 sm:mb-2">
                  <span className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-primary/10 text-primary text-[10px] min-[300px]:text-xs sm:text-sm font-medium">
                    <span className="w-1.5 h-1.5 rounded-full bg-primary animate-pulse" />
                    {activeEdition?.tagline ?? event.acronym}
                  </span>
                </div>
              )}

              {/* Event Name */}
              <h1 className="text-xl min-[300px]:text-2xl sm:text-3xl md:text-4xl font-bold text-foreground leading-tight">
                {event.full_name}
              </h1>
            </div>

            {/* Share Button */}
            <button
              type="button"
              onClick={() =>
                handleShare?.(
                  activeEdition?.name ?? event.full_name,
                  window.location.href,
                )
              }
              className="shrink-0 p-2 sm:p-2.5 rounded-full bg-background/80 backdrop-blur-sm border border-border/50 shadow-sm hover:bg-background hover:scale-105 transition-all duration-200 active:scale-95"
              aria-label="Compartilhar evento"
            >
              <Share2 className="w-4 h-4 min-[300px]:w-5 min-[300px]:h-5 sm:w-5 sm:h-5 text-foreground" />
            </button>
          </div>

          {/* Description */}
          {event.description && (
            <p className="mt-2 sm:mt-3 text-xs min-[300px]:text-sm sm:text-base text-muted-foreground line-clamp-2 sm:line-clamp-3 max-w-2xl">
              {event.description}
            </p>
          )}

          {/* Meta info row */}
          {activeEdition ? (
            <div className="mt-3 sm:mt-4 flex flex-wrap items-center gap-x-4 gap-y-2">
              {/* Date */}
              <div className="flex items-center gap-1.5 text-muted-foreground">
                <Calendar className="w-3.5 h-3.5 min-[300px]:w-4 min-[300px]:h-4 sm:w-4 sm:h-4" />
                <span className="text-xs min-[300px]:text-sm sm:text-sm font-medium">
                  {formatDateRange(
                    activeEdition.starts_at,
                    activeEdition.ends_at,
                  )}
                </span>
              </div>

              {/* Location */}
              {activeEdition.location_name && (
                <div className="flex items-center gap-1.5 text-muted-foreground">
                  <MapPin className="w-3.5 h-3.5 min-[300px]:w-4 min-[300px]:h-4 sm:w-4 sm:h-4" />
                  <span className="text-xs min-[300px]:text-sm sm:text-sm font-medium">
                    {activeEdition.location_name}
                  </span>
                </div>
              )}
            </div>
          ) : (
            <div className="mt-3 sm:mt-4 flex flex-wrap items-center gap-x-4 gap-y-2">
              {/* Acronym */}
              {event.acronym && (
                <div className="flex items-center gap-1.5 text-muted-foreground">
                  <span className="text-xs min-[300px]:text-sm sm:text-sm font-medium">
                    {event.acronym}
                  </span>
                </div>
              )}
            </div>
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
            eventSlug={event.slug}
          />
          <ProductsSection
            editionId={activeEdition?.id}
            eventSlug={event.slug}
          />
          <OtherEditionsSection
            editions={[...pastEditions, ...upcomingEditions]}
            currentEditionId={activeEdition?.id}
            eventSlug={event.slug}
            maxDisplay={5}
          />
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
