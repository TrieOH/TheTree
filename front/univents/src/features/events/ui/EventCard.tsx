import { useNavigate } from "@tanstack/react-router";
import { ArrowUpRight, Eye } from "lucide-react";
import { motion } from "motion/react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import type { EventI } from "../model";

interface EventCardProps {
  event: EventI;
  index?: number;
  className?: string;
  onPublish?: (event: EventI) => void;
}

export function EventCard({
  event,
  index = 0,
  className,
  onPublish,
}: EventCardProps) {
  const navigate = useNavigate();

  const handleClick = () => {
    void navigate({
      to: "/events/$slug",
      params: { slug: event.slug },
    });
  };

  const hasVisual = Boolean(event.banner_url ?? event.logo_url);

  const createdDate = new Date(event.created_at)
    .toLocaleDateString("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    })
    .replace(".", "");

  return (
    <motion.article
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.06,
        duration: 0.4,
        ease: [0.25, 0.1, 0.25, 1],
      }}
      onClick={handleClick}
      className={cn(
        "group relative cursor-pointer min-w-0",
        "bg-card rounded-2xl overflow-hidden",
        "border border-transparent hover:border-border",
        "transition-all duration-300 ease-out",
        "focus:outline-none focus-visible:outline-none focus-visible:ring-0",
        "hover:shadow-lg hover:shadow-foreground/5 hover:-translate-y-1",
        className,
      )}
      role="link"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          handleClick();
        }
      }}
    >
      <div className="overflow-hidden relative bg-muted aspect-4/3">
        {hasVisual ? (
          <img
            src={event.banner_url ?? event.logo_url ?? ""}
            alt=""
            className="w-full h-full object-cover transition-transform duration-700 ease-out group-hover:scale-105"
            loading={index < 4 ? "eager" : "lazy"}
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-linear-to-br from-muted to-muted/50">
            <div className="w-20 h-20 rounded-full border-2 border-dashed border-border/50 flex items-center justify-center">
              <span className="text-2xl font-semibold text-muted-foreground/30">
                {event.acronym ?? event.full_name.charAt(0)}
              </span>
            </div>
          </div>
        )}

        <div className="absolute top-3 right-3 md:top-4 md:right-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
          <div className="w-8 h-8 rounded-full bg-background/90 backdrop-blur-sm flex items-center justify-center">
            <ArrowUpRight className="w-4 h-4 text-foreground" />
          </div>
        </div>
      </div>

      <div className="p-4 md:p-5 space-y-2.5">
        <div className="flex items-start justify-between gap-2">
          <div className="text-xs text-muted-foreground">
            Criado em {createdDate}
          </div>
          {event.status === "draft" && onPublish ? (
            <Button
              type="button"
              size="sm"
              variant="outline"
              className="h-7 px-2.5 text-[11px]"
              onClick={(e) => {
                e.stopPropagation();
                onPublish(event);
              }}
            >
              <Eye className="size-3.5" />
              Publicar
            </Button>
          ) : null}
        </div>

        <h3 className="font-medium text-foreground leading-snug line-clamp-2 group-hover:text-primary transition-colors duration-300 text-base md:text-lg">
          {event.full_name}
        </h3>

        {event.description && (
          <p className="text-sm text-muted-foreground line-clamp-2">
            {event.description}
          </p>
        )}
      </div>
    </motion.article>
  );
}
