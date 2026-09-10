import type { JSX } from "@solidjs/web";
import { untrack } from "solid-js";
import BadgeIcon from "~icons/lucide/badge";
import FileCheckIcon from "~icons/lucide/file-check-2";
import ShoppingBagIcon from "~icons/lucide/shopping-bag";

const Badge = BadgeIcon as unknown as () => JSX.Element;
const FileCheck = FileCheckIcon as unknown as () => JSX.Element;
const ShoppingBag = ShoppingBagIcon as unknown as () => JSX.Element;

type Tab = "badges" | "certificates" | "purchases";

export function ProfileCollectionItem(props: { tab: Tab; item: unknown }) {
  const { tab, rawItem } = untrack(() => ({
    tab: props.tab,
    rawItem: props.item,
  }));
  const item =
    rawItem && typeof rawItem === "object"
      ? (rawItem as Record<string, unknown>)
      : {};
  const title = String(item.name ?? item.title ?? item.edition_name ?? "Item");
  const description = item.description ?? item.status ?? item.event_name;
  const image = typeof item.image === "string" ? item.image : undefined;
  const Icon =
    tab === "badges" ? Badge : tab === "certificates" ? FileCheck : ShoppingBag;

  return (
    <article class="group overflow-hidden rounded-md border border-border bg-card shadow-sm transition hover:-translate-y-0.5 hover:border-primary/50 hover:shadow-md">
      {image && <img src={image} alt="" class="h-32 w-full object-cover" />}
      <div class="flex gap-3 p-4">
        <span class="flex size-9 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
          <Icon />
        </span>
        <div class="min-w-0">
          <h3 class="truncate text-sm font-semibold text-foreground">
            {title}
          </h3>
          {description && (
            <p class="mt-1 line-clamp-2 text-xs text-muted-foreground">
              {String(description)}
            </p>
          )}
        </div>
      </div>
    </article>
  );
}
