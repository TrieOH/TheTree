import type { Event } from "@trieoh/univents-api/schemas";
import { useNavigate } from "@tanstack/solid-router";
import ArrowUpRight from "~icons/lucide/arrow-up-right";
import { animate } from "motion/mini";
import type { JSX } from "@solidjs/web";

const ArrowIcon = ArrowUpRight as unknown as () => JSX.Element;

export function EventCard(props: {
  event: Event;
  index?: number;
  class?: string;
}) {
  const navigate = useNavigate();
  let card!: HTMLElement;
  const index = () => props.index ?? 0;
  const visual = () => props.event.banner_url ?? props.event.logo_url;
  const createdDate = () =>
    new Date(props.event.created_at)
      .toLocaleDateString("pt-BR", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      })
      .replace(".", "");
  const handleClick = () =>
    void navigate({
      to: "/events/$slug",
      params: { slug: props.event.slug },
    } as never);
  return (
    <article
      ref={(element) => {
        card = element;
        requestAnimationFrame(() =>
          animate(
            card,
            {
              opacity: [0, 1],
              transform: ["translateY(20px)", "translateY(0)"],
            },
            {
              delay: index() * 0.06,
              duration: 0.4,
              ease: [0.25, 0.1, 0.25, 1],
            },
          ),
        );
      }}
      class={`group relative min-w-0 cursor-pointer overflow-hidden rounded-2xl border border-transparent bg-card transition-all duration-300 ease-out hover:-translate-y-1 hover:border-border hover:shadow-lg hover:shadow-foreground/5 ${props.class ?? ""}`}
      onClick={handleClick}
      role="link"
      tabindex="0"
      onKeyDown={(event) => {
        if (event.key === "Enter") handleClick();
      }}
    >
      <div class="aspect-4/3 overflow-hidden bg-muted">
        {visual() ? (
          <img
            class="h-full w-full object-cover transition-transform duration-700 group-hover:scale-105"
            src={visual() ?? ""}
            alt=""
            loading={props.index && props.index < 4 ? "eager" : "lazy"}
          />
        ) : (
          <div class="flex h-full items-center justify-center bg-linear-to-br from-muted to-muted/50 text-2xl font-semibold text-muted-foreground/30">
            <div class="flex size-20 items-center justify-center rounded-full border-2 border-dashed border-border/50">
              <span class="text-2xl font-semibold text-muted-foreground/30">
                {props.event.acronym ?? props.event.full_name.charAt(0)}
              </span>
            </div>
          </div>
        )}
        <div class="absolute right-3 top-3 opacity-0 transition-opacity duration-300 group-hover:opacity-100 md:right-4 md:top-4">
          <div class="flex size-8 items-center justify-center rounded-full bg-background/90 backdrop-blur-sm">
              <span class="flex size-4 items-center justify-center">
              <ArrowIcon />
            </span>
          </div>
        </div>
      </div>
      <div class="space-y-2.5 p-4 md:p-5">
        <p class="text-xs text-muted-foreground">Criado em {createdDate()}</p>
        <h3 class="line-clamp-2 text-base font-medium leading-snug transition-colors group-hover:text-primary md:text-lg">
          {props.event.full_name}
        </h3>
        {props.event.description && (
          <p class="line-clamp-2 text-sm text-muted-foreground">
            {props.event.description}
          </p>
        )}
      </div>
    </article>
  );
}
