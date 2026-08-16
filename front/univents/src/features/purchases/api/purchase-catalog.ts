import { orvalData } from "@trieoh/api-client";
import {
  getOccurrence,
  getProduct,
  getTicketType,
  listEditionPrograms,
  listProductVariants,
} from "@trieoh/univents-api";
import type {
  Product,
  Program,
  ProgramOccurrence,
  TicketType,
  Uuid,
} from "@trieoh/univents-api/schemas";
import type { VariantI } from "@/features/products/model";

export type PurchaseCatalogItem = {
  name?: string;
  description?: string | null;
  image?: string | null;
};

const isUuid = (value: string) =>
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(
    value,
  );

export const resolvePurchaseCatalog = async (
  items: Array<{ item_type: string; item_id: string; edition_id: string }>,
) => {
  const entries = await Promise.all(
    items
      .filter((item) => isUuid(item.item_id))
      .map(async (item) => {
        try {
          if (item.item_type === "ticket") {
            const ticket = orvalData<TicketType>(
              await getTicketType(item.item_id as Uuid, { public: true }),
            );
            return [
              `${item.item_type}:${item.item_id}`,
              { name: ticket.name, description: ticket.description },
            ] as const;
          }

          if (item.item_type === "product") {
            const product = orvalData<Product>(
              await getProduct(item.item_id as Uuid, { public: true }),
            );
            const variants = orvalData<VariantI[]>(
              await listProductVariants(product.id, { public: true }),
            );
            const variant = variants[0];
            return [
              `${item.item_type}:${item.item_id}`,
              {
                name: variant?.name,
                description: variant?.description,
                image: variant?.gallery_urls[0],
              },
            ] as const;
          }

          if (item.item_type === "program_occurrence") {
            const occurrence = orvalData<ProgramOccurrence>(
              await getOccurrence(item.item_id as Uuid, { public: true }),
            );
            const programs = orvalData<Program[]>(
              await listEditionPrograms(item.edition_id, { public: true }),
            );
            const program = programs.find(
              (candidate) => candidate.id === occurrence.program_id,
            );
            return [
              `${item.item_type}:${item.item_id}`,
              {
                name: program?.name,
                description: program?.description,
                image: program?.banner_url,
              },
            ] as const;
          }
        } catch {
          return null;
        }
        return null;
      }),
  );

  return Object.fromEntries(
    entries.filter(
      (entry): entry is NonNullable<typeof entry> => entry !== null,
    ),
  );
};
