import { useQueryClient } from "@trieoh/front-core/solid";
import type {
  MyParticipation,
  MyTicket,
  Product,
  ProductVariant,
  Program,
  ProgramOccurrence,
  StoreStockItem,
  TicketType,
} from "@trieoh/univents-api/schemas";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { Loading, createMemo, untrack } from "solid-js";
import {
  productsQueryOptions,
  productVariantsQueryOptions,
  storeStockQueryOptions,
} from "@/features/products/api";
import { useInventoryStream } from "@/features/products/hooks/use-inventory-stream";
import { ProductsSection } from "@/features/products/ui/ProductsSection";
import {
  myParticipationsQueryOptions,
  occurrencesQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import { ProgramSection } from "@/features/programs/ui/ProgramSection";
import {
  myTicketQueryOptions,
  ticketsQueryOptions,
} from "@/features/tickets/api";
import { TicketsSection } from "@/features/tickets/ui/TicketsSection";

export function EventCatalog(props: { editionId: string; eventSlug: string }) {
  const queryClient = useQueryClient();
  const { isAuthenticated } = useAuth();
  const catalog = createMemo(async () => {
    const authenticated = isAuthenticated();
    const [
      tickets,
      programs,
      occurrences,
      products,
      stock,
      heldTicket,
      participations,
    ] = await Promise.all([
      queryClient.fetchQuery(ticketsQueryOptions(props.editionId)),
      queryClient.fetchQuery(programsQueryOptions(props.editionId)),
      queryClient.fetchQuery(occurrencesQueryOptions(props.editionId)),
      queryClient.fetchQuery(productsQueryOptions(props.editionId)),
      queryClient.fetchQuery(storeStockQueryOptions(props.editionId)),
      authenticated
        ? queryClient.fetchQuery(myTicketQueryOptions(props.editionId))
        : null,
      authenticated
        ? queryClient.fetchQuery(myParticipationsQueryOptions(props.editionId))
        : [],
    ]);
    const productsWithVariants = await Promise.all(
      products.map(async (product) => ({
        product,
        variants: await queryClient.fetchQuery(
          productVariantsQueryOptions(product.id),
        ),
      })),
    );
    return {
      tickets,
      programs,
      occurrences,
      productsWithVariants,
      stock,
      heldTicket,
      participations,
      authenticated,
    };
  });

  return (
    <Loading fallback={<CatalogSkeleton />}>
      <LiveCatalog
        editionId={props.editionId}
        eventSlug={props.eventSlug}
        catalog={catalog()}
      />
    </Loading>
  );
}

function LiveCatalog(props: {
  editionId: string;
  eventSlug: string;
  catalog: Catalog;
}) {
  const { editionId, initialStock } = untrack(() => ({
    editionId: props.editionId,
    initialStock: props.catalog.stock,
  }));
  const stock = useInventoryStream(editionId, initialStock);
  return (
    <>
      <TicketsSection
        tickets={props.catalog.tickets}
        stock={stock()}
        eventSlug={props.eventSlug}
        editionId={props.editionId}
        heldTicket={props.catalog.heldTicket}
      />
      <ProgramSection
        programs={props.catalog.programs}
        occurrences={props.catalog.occurrences}
        eventSlug={props.eventSlug}
        editionId={props.editionId}
        authenticated={props.catalog.authenticated}
        heldTicket={props.catalog.heldTicket}
        participations={props.catalog.participations}
      />
      <ProductsSection
        products={props.catalog.productsWithVariants}
        stock={stock()}
        eventSlug={props.eventSlug}
        editionId={props.editionId}
      />
    </>
  );
}

type Catalog = {
  tickets: TicketType[];
  programs: Program[];
  occurrences: ProgramOccurrence[];
  productsWithVariants: { product: Product; variants: ProductVariant[] }[];
  stock: StoreStockItem[];
  heldTicket: MyTicket | null;
  participations: MyParticipation[];
  authenticated: boolean;
};

function CatalogSkeleton() {
  return (
    <div class="grid gap-4 py-10 sm:grid-cols-3">
      <div class="h-48 animate-pulse rounded-2xl bg-muted" />
      <div class="h-48 animate-pulse rounded-2xl bg-muted" />
      <div class="h-48 animate-pulse rounded-2xl bg-muted" />
    </div>
  );
}
