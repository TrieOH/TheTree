import { orvalData } from "@trieoh/api-client";
import {
  getOccurrence,
  listEditionProducts,
  listEditionPrograms,
  listProductVariants,
  listTicketTypes,
} from "@trieoh/univents-api";
import type {
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
  /^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(
    value,
  );

export const resolvePurchaseCatalog = async (
  items: Array<{ item_type: string; item_id: string; edition_id: string }>,
) => {
  const ticketsByEdition = new Map<string, Promise<TicketType[]>>();
  const variantsByEdition = new Map<string, Promise<VariantI[]>>();
  const programsByEdition = new Map<string, Promise<Program[]>>();

  const getTickets = (editionId: string) => {
    let promise = ticketsByEdition.get(editionId);
    if (!promise) {
      promise = listTicketTypes(editionId, { public: true }).then(
        orvalData<TicketType[]>,
      );
      ticketsByEdition.set(editionId, promise);
    }
    return promise;
  };

  const getVariants = (editionId: string) => {
    let promise = variantsByEdition.get(editionId);
    if (!promise) {
      promise = listEditionProducts(editionId, { public: true })
        .then(orvalData<Array<{ id: string }>>)
        .then((products) =>
          Promise.all(
            products.map((product) =>
              listProductVariants(product.id, { public: true }).then(
                orvalData<VariantI[]>,
              ),
            ),
          ),
        )
        .then((variants) => variants.flat());
      variantsByEdition.set(editionId, promise);
    }
    return promise;
  };

  const getPrograms = (editionId: string) => {
    let promise = programsByEdition.get(editionId);
    if (!promise) {
      promise = listEditionPrograms(editionId, { public: true }).then(
        orvalData<Program[]>,
      );
      programsByEdition.set(editionId, promise);
    }
    return promise;
  };

  const entries = await Promise.all(
    items
      .filter((item) => isUuid(item.item_id))
      .map(async (item) => {
        try {
          if (item.item_type === "ticket") {
            const tickets = await getTickets(item.edition_id);
            const ticket = tickets.find(({ id }) => id === item.item_id);
            if (!ticket) return null;
            return [
              `${item.item_type}:${item.item_id}`,
              { name: ticket.name, description: ticket.description },
            ] as const;
          }

          if (item.item_type === "product") {
            const variant = (await getVariants(item.edition_id)).find(
              ({ id }) => id === item.item_id,
            );
            if (!variant) return null;
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
            const programs = await getPrograms(item.edition_id);
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
