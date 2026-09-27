import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import CalendarIcon from "~icons/lucide/calendar";
import LayersIcon from "~icons/lucide/layers";
import PackageIcon from "~icons/lucide/package";
import PencilIcon from "~icons/lucide/pencil";
import ShieldCheckIcon from "~icons/lucide/shield-check";
import TrashIcon from "~icons/lucide/trash";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { ProductI, VariantI } from "../model";

const Calendar = CalendarIcon as unknown as (props: { class?: string }) => JSX.Element;
const Layers = LayersIcon as unknown as (props: { class?: string }) => JSX.Element;
const Package = PackageIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminProductCardProps {
  product: ProductI;
  variants?: VariantI[];
  coverImage?: string | null;
  index?: number;
  animate?: boolean;
  onEdit: (product: ProductI) => void;
  onDelete: (product: ProductI) => void;
  onManageVariants: (product: ProductI) => void;
}

function formatDate(iso: string): string {
  try {
    return new Intl.DateTimeFormat("pt-BR", {
      day: "2-digit",
      month: "short",
      year: "numeric",
    }).format(new Date(iso));
  } catch {
    return "";
  }
}

export function AdminProductCard(props: AdminProductCardProps): JSX.Element {
  const isDeleted = () => props.product.deleted_at !== null;

  const firstImage = () => {
    if (props.coverImage) return props.coverImage;
    if (props.variants) {
      for (const v of props.variants) {
        if (v.gallery_urls && v.gallery_urls.length > 0 && v.gallery_urls[0]) {
          return v.gallery_urls[0];
        }
      }
    }
    return null;
  };

  const variantsCount = () => props.variants?.length;

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div
        role="region"
        aria-label={`Produto ${props.product.vendor_code}`}
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border border-t-[3px] border-t-primary bg-card bg-linear-to-b from-primary/[0.04] via-card to-card p-4 text-left shadow-xs transition-all duration-200",
          "hover:-translate-y-0.5 hover:border-foreground/20 hover:border-t-primary hover:shadow-md",
          "focus-within:ring-2 focus-within:ring-ring",
          isDeleted() && "opacity-60 grayscale",
        )}
      >
        <div class="space-y-3">
          {/* Header with thumbnail, code, badges and actions */}
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-2.5 min-w-0 flex-1">
              <div class="relative flex size-11 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-border/70 bg-muted/60 transition-transform group-hover:scale-105">
                <Show
                  when={firstImage()}
                  fallback={
                    <div class="flex size-full items-center justify-center bg-primary/10 text-primary">
                      <Package class="size-5" />
                    </div>
                  }
                >
                  {(url) => (
                    <img
                      src={url()}
                      alt={props.product.vendor_code}
                      class="size-full object-cover"
                    />
                  )}
                </Show>
              </div>

              <div class="min-w-0 flex-1">
                <h3
                  class="truncate text-sm font-bold font-mono tracking-tight text-foreground transition-colors group-hover:text-primary"
                  title={props.product.vendor_code}
                >
                  {props.product.vendor_code}
                </h3>
                <div class="mt-0.5 flex flex-wrap items-center gap-1.5">
                  <Show
                    when={props.product.requires_registration}
                    fallback={
                      <span class="inline-flex items-center text-[10px] font-medium text-muted-foreground">
                        Livre
                      </span>
                    }
                  >
                    <span class="inline-flex items-center gap-1 text-[10px] font-medium text-primary">
                      <ShieldCheck class="size-3 shrink-0" />
                      Exige cadastro
                    </span>
                  </Show>

                  <Show when={variantsCount() != null && variantsCount()! > 0}>
                    <span class="text-[10px] text-muted-foreground/60">•</span>
                    <span class="inline-flex items-center gap-1 text-[10px] font-medium text-muted-foreground">
                      <Layers class="size-2.5" />
                      {variantsCount()} {variantsCount() === 1 ? "variação" : "variações"}
                    </span>
                  </Show>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div class="flex items-center gap-0.5 shrink-0">
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Editar ${props.product.vendor_code}`}
                title="Editar produto"
                onClick={() => props.onEdit(props.product)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100 cursor-pointer"
              >
                <Pencil class="size-3.5" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Excluir ${props.product.vendor_code}`}
                title="Excluir produto"
                onClick={() => props.onDelete(props.product)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-destructive/15 hover:text-destructive group-hover:opacity-100 cursor-pointer"
              >
                <Trash class="size-3.5" />
              </Button>
            </div>
          </div>

          <div class="flex items-center justify-between text-[11px] text-muted-foreground min-h-[1.25rem]">
            <span class="inline-flex items-center gap-1.5">
              <Calendar class="size-3 shrink-0 opacity-70" />
              <span>Criado em {formatDate(props.product.created_at)}</span>
            </span>
          </div>
        </div>

        {/* Footer: Manage Variants button */}
        <div class="mt-3 border-t border-border/60 pt-2.5">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => props.onManageVariants(props.product)}
            class="w-full gap-1.5 text-xs font-medium cursor-pointer hover:bg-accent hover:text-accent-foreground"
          >
            <Layers class="size-3.5 text-muted-foreground" />
            Gerenciar variações
          </Button>
        </div>
      </div>
    </Reveal>
  );
}
