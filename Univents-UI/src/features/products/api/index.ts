import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { ProductCreateI, ProductI } from "../model";
import { authFetcher, tanstackQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Product on the server.
 * @param productData - The data for the new product.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to the API response containing the newly created product.
 */
export const createProductFn = createClientOnlyFn((
  productData: ProductCreateI, eventId: string, editionId: string
) => {
  return authFetcher.post<ProductI>(
    `/events/${eventId}/editions/${editionId}/products`,
    productData
  );
});

/**
 * Fetches all products for a specific edition from the server.
 * @returns A promise that resolves to an array of Product objects.
 */
export const getAllProductsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<ProductI[]>(`/events/${eventId}/editions/${editionId}/products`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all products for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all products for a specific edition.
 */
export const allProductsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['products', 'public', eventId, editionId],
    queryFn: () => getAllProductsFn(eventId, editionId),
  })
}

/**
 * Fetches all admin products for a specific edition from the server.
 * @returns A promise that resolves to an array of Product objects.
 */
export const getAllAdminProductsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  try {
    return await tanstackQueryFetcher<ProductI[]>(`/events/${eventId}/editions/${editionId}/products/admin`);
  } catch {
    return [];
  }
});

/**
 * Query options for fetching all admin products for a specific edition, using TanStack Query.
 * @returns An object containing the query key and query function for fetching all admin products for a specific edition.
 */
export const allAdminProductsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['products', 'admin', eventId, editionId],
    queryFn: () => getAllAdminProductsFn(eventId, editionId),
  })
};

/**
 * Publish a Product on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API null response.
 */
export const publishProductFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/publish`
  );
});
