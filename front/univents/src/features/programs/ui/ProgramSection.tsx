import { useSuspenseQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { storeStockQueryOptions } from "@/features/products/api";
import { cn } from "@/shared/lib/utils";
import { occurrencesQueryOptions, programsQueryOptions } from "../api";
import { useProgramRegistration } from "../hooks/use-program-registration";
import type { OccurrenceI, ProgramI } from "../model";
import { ProgramDayCard } from "./ProgramDayCard";

interface ProgramSectionProps {
  editionId?: string;
  eventId: string;
  eventSlug: string;
}

export function groupByDay(programs: ProgramI[], occurrences: OccurrenceI[]) {
  const programById = new Map(programs.map((p) => [p.id, p]));

  const groups = new Map<
    string,
    { date: string; items: { program: ProgramI; occurrence: OccurrenceI }[] }
  >();

  for (const occurrence of occurrences) {
    const program = programById.get(occurrence.program_id);
    if (!program) continue;

    const dateKey = occurrence.starts_at.slice(0, 10);
    const existing = groups.get(dateKey);

    if (existing) {
      existing.items.push({ program, occurrence });
    } else {
      groups.set(dateKey, {
        date: occurrence.starts_at,
        items: [{ program, occurrence }],
      });
    }
  }

  return Array.from(groups.values())
    .sort((a, b) => a.date.localeCompare(b.date))
    .map((group) => ({
      ...group,
      items: group.items.sort((a, b) =>
        a.occurrence.starts_at.localeCompare(b.occurrence.starts_at),
      ),
    }));
}

export function ProgramSection({
  editionId,
  eventId,
  eventSlug,
}: ProgramSectionProps) {
  if (!editionId) return null;

  const { data: programs } = useSuspenseQuery(programsQueryOptions(editionId));
  const { data: occurrences } = useSuspenseQuery(
    occurrencesQueryOptions(editionId),
  );
  const { data: stock } = useSuspenseQuery(storeStockQueryOptions(editionId));
  const stockByOccurrence = new Map(
    stock
      .filter((item) => item.item_type === "program_occurrence")
      .map((item) => [item.id, item.stock]),
  );
  const registration = useProgramRegistration(
    editionId,
    eventId,
    eventSlug,
    `/events/${eventSlug}`,
  );

  const days = groupByDay(programs, occurrences);

  if (days.length === 0) return null;

  const visibleDays = days.slice(0, 3);
  const visibilityClasses = ["flex", "hidden md:flex", "hidden lg:flex"];
  const separatorClasses = ["hidden md:block", "hidden lg:block"];

  return (
    <section className="w-full py-10">
      <div className="px-4">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-semibold text-foreground tracking-tight">
            Programação
          </h2>
          <p className="mt-2 text-sm text-muted-foreground leading-relaxed max-w-xl mx-auto">
            Confira as atividades e checkpoints do evento.
          </p>
        </div>

        {/* Cards de dia com separador */}
        <div className="flex flex-nowrap justify-center gap-5 mx-auto">
          {visibleDays.map((day, index) => (
            <div
              key={day.date.slice(0, 10)}
              className={cn("items-start gap-5", visibilityClasses[index])}
            >
              <ProgramDayCard
                date={day.date}
                items={day.items}
                maxItems={3}
                registration={registration}
                stockByOccurrence={stockByOccurrence}
              />
              {index < visibleDays.length - 1 && (
                <div
                  className={cn(
                    "w-px self-stretch bg-border rounded-full",
                    separatorClasses[index],
                  )}
                />
              )}
            </div>
          ))}
        </div>

        <div className="mt-6 text-center">
          <Link
            to="/events/$slug/programs"
            params={{ slug: eventSlug }}
            className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:gap-2.5 transition-all duration-200"
          >
            Programação Completa
            <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      </div>
    </section>
  );
}
