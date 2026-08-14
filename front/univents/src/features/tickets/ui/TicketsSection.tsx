import { Link } from "@tanstack/react-router";
import type { MyTicket } from "@trieoh/univents-api/schemas";
import { ArrowRight } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import type { TicketI } from "../model";
import { TicketCard } from "./TicketCard";

interface TicketsSectionProps {
  tickets: TicketI[];
  eventSlug: string;
  editionId?: string;
  heldTicket?: MyTicket | null;
}

export function TicketsSection({
  tickets,
  eventSlug,
  editionId,
  heldTicket = null,
}: TicketsSectionProps) {
  if (tickets.length === 0) return null;

  const sortedTickets = [...tickets]
    .sort((a, b) => b.access_level - a.access_level)
    .slice(0, 3);

  const visibilityClasses = ["", "hidden sm:block", "hidden lg:block"];

  return (
    <section className="w-full py-10">
      <div className="px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-semibold text-foreground tracking-tight">
            Ingressos
          </h2>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed max-w-xl mx-auto">
            Escolha o tipo de ingresso que melhor se encaixa na sua experiência
            no evento.
          </p>
        </div>

        <div className="flex justify-center gap-4 max-w-5xl mx-auto">
          {sortedTickets.map((ticket, index) => (
            <div
              key={ticket.id}
              className={cn("shrink-0", visibilityClasses[index])}
            >
              <TicketCard
                ticket={ticket}
                isFeatured={index === 1}
                editionId={editionId}
                heldTicket={heldTicket}
              />
            </div>
          ))}
        </div>

        <div className="mt-6 text-center">
          <Link
            to="/events/$slug/tickets"
            params={{ slug: eventSlug }}
            className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:gap-2.5 transition-all duration-200"
          >
            Ver todos os ingressos
            <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      </div>
    </section>
  );
}
