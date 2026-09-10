import type { Product, ProductVariant } from "@trieoh/univents-api/schemas";
import { For, createSignal, untrack } from "solid-js";
import type { JSX } from "@solidjs/web";
import ChevronLeftIcon from "~icons/lucide/chevron-left";
import ChevronRightIcon from "~icons/lucide/chevron-right";
import PackageIcon from "~icons/lucide/package";
import UserCheckIcon from "~icons/lucide/user-check";
import ShoppingCartIcon from "~icons/lucide/shopping-cart";
import { useCart } from "../hooks/use-cart";

const ChevronLeft = ChevronLeftIcon as unknown as () => JSX.Element;
const ChevronRight = ChevronRightIcon as unknown as () => JSX.Element;
const Package = PackageIcon as unknown as () => JSX.Element;
const UserCheck = UserCheckIcon as unknown as () => JSX.Element;
const ShoppingCart = ShoppingCartIcon as unknown as () => JSX.Element;

const price = (cents: number) =>
  new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(cents / 100);

export function ProductCard(props: {
  product: Product;
  variants: ProductVariant[];
  stock: ReadonlyMap<string, number | null>;
  editionId: string;
}) {
  return (
    <article class="w-88 max-w-full space-y-2">
      <div class="relative flex h-10 items-center border-b border-border/50 px-1">
        <p class="text-[11px] font-medium uppercase tracking-wider text-muted-foreground">
          {props.product.vendor_code}
        </p>
        {props.product.requires_registration && (
          <span
            title="Necessita cadastro no evento"
            class="absolute right-0 top-1/2 flex size-7 -translate-y-1/2 items-center justify-center rounded-full bg-card text-primary"
          >
            <UserCheck />
          </span>
        )}
      </div>
      <div class="flex flex-col gap-2">
        <For each={props.variants.slice(0, 3)}>
          {(variant) => (
            <VariantCard
              variant={variant}
              stock={props.stock.get(variant.id)}
              editionId={props.editionId}
              productCode={props.product.vendor_code}
            />
          )}
        </For>
      </div>
    </article>
  );
}

function VariantCard(props: {
  variant: ProductVariant;
  stock?: number | null;
  editionId: string;
  productCode: string;
}) {
  const cart = useCart(untrack(() => props.editionId));
  const [image, setImage] = createSignal(0);
  const images = untrack(() => props.variant.gallery_urls);
  const move = (direction: number) =>
    setImage(
      (current) => (current + direction + images.length) % images.length,
    );
  const soldOut = () => props.stock === 0;
  const inCart = () =>
    cart
      .items()
      .find((item) => item.type === "product" && item.id === props.variant.id);

  return (
    <div class="group rounded-xl border border-border/60 bg-card p-4">
      <div class="flex h-16 items-start gap-3">
        <div class="relative size-16 shrink-0 overflow-hidden rounded-md bg-muted ring-1 ring-border/60">
          {images[image()] ? (
            <img
              src={images[image()]}
              alt={props.variant.name}
              class="size-full object-cover"
            />
          ) : (
            <span class="flex size-full items-center justify-center text-muted-foreground/50">
              <Package />
            </span>
          )}
          {images.length > 1 && (
            <>
              <button
                type="button"
                aria-label="Imagem anterior"
                class="absolute left-0.5 top-1/2 flex size-5 -translate-y-1/2 items-center justify-center rounded-full bg-black/50 text-white opacity-0 transition-opacity group-hover:opacity-100"
                onClick={() => move(-1)}
              >
                <ChevronLeft />
              </button>
              <button
                type="button"
                aria-label="Próxima imagem"
                class="absolute right-0.5 top-1/2 flex size-5 -translate-y-1/2 items-center justify-center rounded-full bg-black/50 text-white opacity-0 transition-opacity group-hover:opacity-100"
                onClick={() => move(1)}
              >
                <ChevronRight />
              </button>
            </>
          )}
        </div>
        <div class="flex h-16 min-w-0 flex-1 flex-col justify-center">
          <p class="truncate text-sm font-semibold leading-5 text-foreground">
            {props.variant.name}
          </p>
          <p class="truncate text-xs leading-5 text-muted-foreground">
            {props.variant.description}
          </p>
        </div>
        <div class="flex h-16 shrink-0 flex-col justify-center text-right">
          <span class="text-sm font-bold tabular-nums">
            {price(props.variant.price)}
          </span>
          <span
            class={`mt-0.5 text-[11px] ${soldOut()
              ? "font-medium text-destructive"
              : "text-muted-foreground"
              }`}
          >
            {props.stock === null
              ? "Ilimitado"
              : soldOut()
                ? "Esgotado"
                : props.stock}
          </span>
        </div>
      </div>
      <button
        type="button"
        disabled={soldOut()}
        class="mt-3 flex h-9 w-full items-center justify-center gap-2 rounded-md bg-primary text-xs font-semibold text-primary-foreground disabled:bg-muted disabled:text-muted-foreground"
        onClick={() =>
          untrack(() =>
            cart.add({
              id: props.variant.id,
              type: "product",
              name: `${props.productCode} - ${props.variant.name}`,
              price_cents: props.variant.price,
              stock: props.stock ?? props.variant.stock ?? null,
            })
          )
        }
      >
        <span class="size-4">
          <ShoppingCart />
        </span>
        {soldOut()
          ? "Indisponível"
          : inCart()
            ? `Adicionado (${inCart()?.quantity})`
            : "Adicionar ao carrinho"}
      </button>
    </div>
  );
}
