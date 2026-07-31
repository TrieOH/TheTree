import { Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import type { EditionI } from "../model";
import { EditionSummaryCard } from "./EditionSummaryCard";

function sortEditions(editions: EditionI[]): EditionI[] {
  const now = Date.now();

  return [...editions].sort((a, b) => {
    const aEnd = new Date(a.ends_at).getTime();
    const bEnd = new Date(b.ends_at).getTime();
    const aStart = new Date(a.starts_at).getTime();
    const bStart = new Date(b.starts_at).getTime();

    const aIsFuture = aEnd >= now;
    const bIsFuture = bEnd >= now;

    if (aIsFuture && !bIsFuture) return -1;
    if (!aIsFuture && bIsFuture) return 1;
    if (aIsFuture && bIsFuture) return aStart - bStart;

    return bEnd - aEnd;
  });
}

interface OtherEditionsSectionProps {
  editions: EditionI[];
  currentEditionId?: string;
  eventSlug: string;
  maxDisplay?: number;
}

export function OtherEditionsSection({
  editions,
  currentEditionId,
  eventSlug,
  maxDisplay = 6,
}: OtherEditionsSectionProps) {
  const filtered = editions
    .filter((ed) => ed.id !== currentEditionId)
    .slice(0, maxDisplay);

  const sorted = sortEditions(filtered);
  if (sorted.length === 0) return null;

  return (
    <section className="w-full">
      {/* Header */}
      <div className="mb-5 sm:mb-6">
        <h2 className="text-2xl sm:text-3xl font-normal text-primary tracking-tight">
          Outras Edições
        </h2>
        <p className="mt-1 text-sm sm:text-base text-muted-foreground">
          Conheça os próximos eventos ou relembre os destaques das edições
          anteriores.
        </p>
      </div>

      <div className="flex flex-wrap gap-4">
        {sorted.map((edition) => (
          <EditionSummaryCard key={edition.id} edition={edition} />
        ))}
      </div>

      {/* View all */}
      <div className="mt-5">
        <Link
          to="/events/$slug/editions"
          params={{ slug: eventSlug }}
          className="inline-flex items-center gap-1.5 text-sm font-semibold text-primary hover:gap-2.5 transition-all duration-200"
        >
          Ver todas as Edições
          <ArrowRight className="w-4 h-4" />
        </Link>
      </div>
    </section>
  );
}
