import { ArrowRight, Calendar, MapPin } from "lucide-react";
import { formatDateRange } from "@/shared/lib/date";
import { getInitials } from "@/shared/lib/share";
import { cn } from "@/shared/lib/utils";
import type { EditionI, EditionStatus } from "../model";

function statusBadge(status: EditionStatus): {
  label: string;
  className: string;
} {
  switch (status) {
    case "future":
      return {
        label: "UPCOMING",
        className:
          "bg-blue-100 text-blue-700 dark:bg-blue-950/40 dark:text-blue-400",
      };
    case "active":
      return {
        label: "LIVE",
        className:
          "bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:text-amber-400",
      };
    default:
      return {
        label: "CLOSED",
        className:
          "bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400",
      };
  }
}

function actionLabel(status: EditionStatus): string {
  switch (status) {
    case "future":
      return "Learn More";
    case "active":
      return "Join Now";
    default:
      return "View Archives";
  }
}

interface EditionSummaryCardProps {
  edition: EditionI;
}

export function EditionSummaryCard({ edition }: EditionSummaryCardProps) {
  const initials = getInitials(edition.name);
  const badge = statusBadge(edition.status);

  return (
    <div
      className={cn(
        "group flex rounded-xl bg-card border border-border/60 overflow-hidden",
        "hover:border-border hover:shadow-md transition-all duration-200",
        "min-w-64 flex-1 max-w-96",
      )}
    >
      {/* Image left */}
      <div className="relative w-24 shrink-0 overflow-hidden bg-muted">
        {edition.banner_url ? (
          <img
            src={edition.banner_url}
            alt={edition.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : edition.logo_url ? (
          <img
            src={edition.logo_url}
            alt={edition.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-muted to-muted/50">
            <span className="text-xl sm:text-2xl font-bold text-muted-foreground/50">
              {initials}
            </span>
          </div>
        )}
      </div>

      {/* Content right */}
      <div className="flex flex-col justify-between flex-1 p-4 min-w-0">
        {/* Top: Title + Badge */}
        <div className="flex items-start justify-between gap-2">
          <h3 className="text-sm sm:text-base font-semibold text-foreground leading-tight line-clamp-2">
            {edition.name}
          </h3>
          <span
            className={cn(
              "shrink-0 inline-block px-1.5 py-0.5 rounded text-[9px] sm:text-[10px] font-bold tracking-wide",
              badge.className,
            )}
          >
            {badge.label}
          </span>
        </div>

        {/* Meta */}
        <div className="mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-muted-foreground">
          <div className="flex items-center gap-1">
            <Calendar className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
            <span className="text-[11px] sm:text-xs">
              {formatDateRange(edition.starts_at, edition.ends_at)}
            </span>
          </div>
          {edition.location_name && (
            <div className="flex items-center gap-1">
              <MapPin className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
              <span className="text-[11px] sm:text-xs truncate max-w-25 sm:max-w-35">
                {edition.location_name}
              </span>
            </div>
          )}
        </div>

        {/* Action */}
        <div className="mt-2 sm:mt-3">
          <span className="inline-flex items-center gap-1 text-xs sm:text-sm font-semibold text-primary group-hover:gap-1.5 transition-all duration-200">
            {actionLabel(edition.status)}
            <ArrowRight className="w-3 h-3 sm:w-3.5 sm:h-3.5" />
          </span>
        </div>
      </div>
    </div>
  );
}
