import { Link, useLocation } from "@tanstack/react-router";
import { ChevronRight } from "lucide-react";
import { Fragment } from "react";
import { cn } from "#/shared/lib/utils";

export function Breadcrumb() {
  const { pathname } = useLocation();

  const segments = pathname.split("/").filter(Boolean);

  return (
    <nav
      className={cn(
        "flex items-center space-x-2 text-muted-foreground",
        "font-bold uppercase tracking-[0.2em] text-[10px]",
        "px-6 h-16 border-b border-border/60",
        "bg-background/95 backdrop-blur-md",
        "overflow-x-auto whitespace-nowrap",
      )}
    >
      {segments.map((segment, index) => {
        const isLast = index === segments.length - 1;
        const path = `/${segments.slice(0, index + 1).join("/")}`;

        const label = segment.charAt(0).toUpperCase() + segment.slice(1);

        const isUUID =
          /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
            segment,
          );
        const displayLabel =
          isUUID || label.length > 20
            ? `${label.slice(0, 4)}...${label.slice(-2)}`
            : label;

        return (
          <Fragment key={path}>
            {index > 0 && (
              <ChevronRight className="h-3 w-3 shrink-0 text-muted-foreground/40" />
            )}
            {isLast ? (
              <span className="max-w-37.5 truncate text-foreground">
                {displayLabel}
              </span>
            ) : (
              <Link
                to={path}
                className="max-w-37.5 truncate transition-colors hover:text-primary"
              >
                {displayLabel}
              </Link>
            )}
          </Fragment>
        );
      })}
    </nav>
  );
}
