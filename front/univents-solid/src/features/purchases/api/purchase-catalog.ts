import type { QueryClient } from "@tanstack/query-core";
import {
  occurrenceQueryOptions,
  programsQueryOptions,
} from "@/features/programs/api";
import {
  productsQueryOptions,
  productVariantsQueryOptions,
} from "@/features/products/api";
import { ticketsQueryOptions } from "@/features/tickets/api";

export type PurchaseCatalog = Record<
  string,
  { name?: string; description?: string | null; image?: string | null }
>;

export type PurchaseCatalogItem = {
  item_type: string;
  item_id: string;
  edition_id: string;
};

const UUID_PATTERN =
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

const isUuid = (value: string) => UUID_PATTERN.test(value);

/**
 * Resolves the display data (name/description/image) for the items of a
 * purchase. Every lookup goes through the shared TanStack Query cache via
 * the existing per-feature query options, so a catalog never refetches a
 * list another screen already loaded.
 */
export async function resolvePurchaseCatalog(
  queryClient: QueryClient,
  items: PurchaseCatalogItem[],
): Promise<PurchaseCatalog> {
  const getVariants = async (editionId: string) => {
    const products = await queryClient.fetchQuery(
      productsQueryOptions(editionId),
    );
    const groups = await Promise.all(
      products.map((product) =>
        queryClient.fetchQuery(productVariantsQueryOptions(product.id)),
      ),
    );
    return groups.flat();
  };

  const entries = await Promise.all(
    items
      .filter((item) => isUuid(item.item_id))
      .map(async (item) => {
        try {
          if (item.item_type === "ticket") {
            const tickets = await queryClient.fetchQuery(
              ticketsQueryOptions(item.edition_id),
            );
            const ticket = tickets.find(({ id }) => id === item.item_id);
            if (!ticket) return null;

            return [
              `ticket:${item.item_id}`,
              { name: ticket.name, description: ticket.description },
            ] as const;
          }

          if (item.item_type === "product") {
            const variant = (await getVariants(item.edition_id)).find(
              ({ id }) => id === item.item_id,
            );
            if (!variant) return null;

            return [
              `product:${item.item_id}`,
              {
                name: variant.name,
                description: variant.description,
                image: variant.gallery_urls[0],
              },
            ] as const;
          }

          const occurrence = await queryClient.fetchQuery(
            occurrenceQueryOptions(item.item_id),
          );
          const programs = await queryClient.fetchQuery(
            programsQueryOptions(item.edition_id),
          );
          const program = programs.find((c) => c.id === occurrence.program_id);

          return [
            `program_occurrence:${item.item_id}`,
            {
              name: program?.name,
              description: program?.description,
              image: program?.banner_url,
            },
          ] as const;
        } catch {
          return null;
        }
      }),
  );

  return Object.fromEntries(
    entries.filter(
      (entry): entry is NonNullable<typeof entry> => entry !== null,
    ),
  );
}
