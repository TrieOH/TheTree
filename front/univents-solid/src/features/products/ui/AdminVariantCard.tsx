import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";

import BoxesIcon from "~icons/lucide/boxes";
import ImagePlusIcon from "~icons/lucide/image-plus";
import InfinityIcon from "~icons/lucide/infinity";
import PencilIcon from "~icons/lucide/pencil";
import TrashIcon from "~icons/lucide/trash";

import { Button, cn } from "@trieoh/ui-solid";
import { Reveal } from "@/shared/ui/Reveal";
import type { VariantI } from "../model";

const Boxes = BoxesIcon as unknown as (props: { class?: string }) => JSX.Element;
const ImagePlus = ImagePlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const InfinityLucide = InfinityIcon as unknown as (props: { class?: string }) => JSX.Element;
const Pencil = PencilIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminVariantCardProps {
  variant: VariantI;
  index?: number;
  animate?: boolean;
  onEdit: (variant: VariantI) => void;
  onDelete: (variant: VariantI) => void;
  onManageGallery?: (variant: VariantI) => void;
}

function formatPrice(cents: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);
}

export function AdminVariantCard(props: AdminVariantCardProps): JSX.Element {
  const isDeleted = () => props.variant.deleted_at !== null;
  const imageCount = () => props.variant.gallery_urls?.length ?? 0;
  const coverImage = () => props.variant.gallery_urls?.[0];

  return (
    <Reveal delay={(props.index ?? 0) * 0.05} animate={props.animate} class="h-full">
      <div
        role="region"
        aria-label={`Variação ${props.variant.name}`}
        class={cn(
          "group relative flex h-full w-full min-w-0 flex-col justify-between overflow-hidden rounded-xl border border-border border-t-[3px] border-t-primary bg-card bg-linear-to-b from-primary/[0.04] via-card to-card p-4 text-left shadow-xs transition-all duration-200",
          "hover:-translate-y-0.5 hover:border-foreground/20 hover:border-t-primary hover:shadow-md",
          "focus-within:ring-2 focus-within:ring-ring",
          isDeleted() && "opacity-60 grayscale",
        )}
      >
        <div class="space-y-2">
          {/* Header with thumbnail, name, code and actions */}
          <div class="flex items-start justify-between gap-3">
            <div class="flex items-center gap-2.5 min-w-0 flex-1">
              <div class="relative flex size-11 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-border/70 bg-muted/60 transition-transform group-hover:scale-105">
                <Show
                  when={coverImage()}
                  fallback={
                    <div class="flex size-full items-center justify-center bg-primary/10 text-primary font-bold uppercase select-none">
                      {props.variant.name.charAt(0) || "V"}
                    </div>
                  }
                >
                  {(url) => (
                    <img
                      src={url()}
                      alt={props.variant.name}
                      class="size-full object-cover"
                    />
                  )}
                </Show>
              </div>

              <div class="min-w-0 flex-1">
                <h3
                  class="truncate text-sm font-semibold tracking-tight text-foreground transition-colors group-hover:text-primary"
                  title={props.variant.name}
                >
                  {props.variant.name}
                </h3>
                <div class="mt-0.5 flex items-center gap-1.5 text-xs text-muted-foreground">
                  <span class="truncate font-mono" title={props.variant.vendor_code}>
                    {props.variant.vendor_code}
                  </span>
                  <span>•</span>
                  <span class="shrink-0 text-[11px]">
                    {imageCount()} {imageCount() === 1 ? "foto" : "fotos"}
                  </span>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div class="flex items-center gap-0.5 shrink-0">
              <Show when={props.onManageGallery}>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  aria-label={`Galeria de ${props.variant.name}`}
                  title="Gerenciar fotos da galeria"
                  onClick={() => props.onManageGallery?.(props.variant)}
                  class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100 cursor-pointer"
                >
                  <ImagePlus class="size-3.5" />
                </Button>
              </Show>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Editar ${props.variant.name}`}
                title="Editar variação"
                onClick={() => props.onEdit(props.variant)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-muted hover:text-foreground group-hover:opacity-100 cursor-pointer"
              >
                <Pencil class="size-3.5" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                aria-label={`Excluir ${props.variant.name}`}
                title="Excluir variação"
                onClick={() => props.onDelete(props.variant)}
                class="size-7 text-muted-foreground opacity-60 transition-opacity hover:bg-destructive/15 hover:text-destructive group-hover:opacity-100 cursor-pointer"
              >
                <Trash class="size-3.5" />
              </Button>
            </div>
          </div>

          <p
            class="line-clamp-2 text-xs leading-relaxed text-muted-foreground min-h-[2rem]"
            title={props.variant.description ?? ""}
          >
            {props.variant.description || "Nenhuma descrição fornecida para esta variação."}
          </p>
        </div>

        {/* Specs footer: Price and Stock */}
        <div class="mt-3 flex items-center justify-between border-t border-border/60 pt-2.5 text-[11px] text-muted-foreground">
          <div class="flex items-center gap-1.5 font-semibold text-foreground">
            <span>{formatPrice(props.variant.price)}</span>
          </div>

          <div class="flex items-center gap-1.5">
            <Show
              when={props.variant.stock != null}
              fallback={
                <span class="inline-flex items-center gap-1 text-emerald-600 dark:text-emerald-400 font-medium">
                  <InfinityLucide class="size-3" />
                  <span>Ilimitado</span>
                </span>
              }
            >
              <span class="inline-flex items-center gap-1 text-foreground/80 font-medium">
                <Boxes class="size-3 text-muted-foreground shrink-0" />
                <span>{props.variant.stock} un. em estoque</span>
              </span>
            </Show>
          </div>
        </div>
      </div>
    </Reveal>
  );
}
