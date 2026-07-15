import { createClientOnlyFn } from "@tanstack/react-start";
import { queryOptions } from "@tanstack/react-query";
import type { ProductCreateOutputI, ProductI } from "../model";
import type { ImageURLUploadI } from "@/shared/model/generic";
import { authFetcher, authQueryFetcher, publicQueryFetcher } from "@/shared/lib/api/fetch";

/**
 * Creates a new Product on the server.
 * @param productData - The data for the new product.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to the API response containing the newly created product.
 */
export const createProductFn = createClientOnlyFn((
  productData: ProductCreateOutputI, eventId: string, editionId: string
) => {
  return authFetcher.post<ProductI>(
    `/events/${eventId}/editions/${editionId}/products`,
    productData
  );
});

/**
 * Patches a Product on the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @param productData - The data for the product.
 * @returns A promise that resolves to the API response containing the updated product.
 */
export const patchProductFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string, productData: ProductCreateOutputI
) => {
  return authFetcher.patch<ProductI>(
    `/events/${eventId}/editions/${editionId}/products/${productId}`,
    productData
  );
});

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

/**
 * Restores a soft deleted product. Only works if the product has not been hard deleted yet.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API null response.
 */
export const restoreSoftDeletedProductFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string
) => {
  return authFetcher.post<null>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/restore`
  );
});

/**
 * Soft Delete a product
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API null response.
 */
export const softDeleteProductFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string
) => {
  return authFetcher.delete<null>(
    `/events/${eventId}/editions/${editionId}/products/${productId}`
  );
});

/**
 * Adds a MinIO URL to the product's gallery_urls array.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API ProductI response.
 */
export const addImageToTheProductGalleryFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.post<ProductI>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/gallery`,
    urlData
  );
});

/**
 * Removes a URL from the product's gallery_urls array and deletes the object from MinIO.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API ProductI response.
 */
export const removeImageToTheProductGalleryFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.delete<ProductI>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/gallery`,
    urlData
  );
});

/**
 * Sets the product thumbnail URL. If the URL is not already in gallery_urls it is added automatically.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API ProductI response.
 */
export const setProductThumbnailFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string, urlData: ImageURLUploadI
) => {
  return authFetcher.put<ProductI>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/thumbnail`,
    urlData
  );
});

/**
 * Clears the product thumbnail. The image remains in gallery_urls.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @returns A promise that resolves to the API ProductI response.
 */
export const unsetProductThumbnailFn = createClientOnlyFn((
  eventId: string, editionId: string, productId: string
) => {
  return authFetcher.delete<ProductI>(
    `/events/${eventId}/editions/${editionId}/products/${productId}/thumbnail`
  );
});

/**
 * Fetches all products for a specific edition from the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to an array of Product objects.
 */
export const getAllPublicProductsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  return publicQueryFetcher<ProductI[]>(`/events/${eventId}/editions/${editionId}/products`);
});

/**
 * Query options for fetching all products for a specific edition, using TanStack Query.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns An object containing the query key and query function for fetching all products for a specific edition.
 */
export const allPublicProductsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['products', 'public', eventId, editionId],
    queryFn: () => getAllPublicProductsFn(eventId, editionId),
  })
};

/**
 * Fetches all admin products for a specific edition from the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns A promise that resolves to an array of Product objects.
 */
export const getAllAdminProductsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  return authQueryFetcher<ProductI[]>(`/events/${eventId}/editions/${editionId}/products/admin`);
});

/**
 * Query options for fetching all admin products for a specific edition, using TanStack Query.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @returns An object containing the query key and query function for fetching all admin products for a specific edition.
 */
export const allAdminProductsQueryOptions = (eventId: string, editionId: string) => {
  return queryOptions({
    queryKey: ['products', 'admin', eventId, editionId],
    queryFn: () => getAllAdminProductsFn(eventId, editionId),
  })
};


// FIXME: I NEED TO DELETE EVERYTHING BELOW THIS LINE AND REPLACE IT

/**
 * Fetches all products for a specific edition from the server.
 * @param eventId - The event id
 * @param editionId - The edition id
 * @param productId - The product id
 * @param productData - The data for the product.
 * @returns A promise that resolves to an array of Product objects.
 */
export const getAllProductsFn = createClientOnlyFn(async (eventId: string, editionId: string) => {
  authQueryFetcher<ProductI[]>(`/events/${eventId}/editions/${editionId}/products`);
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

// FIXME: I NEED TO MOVE THIS TO OTHER FILE

/**
 * Get websocket token.
 * @returns A promise that resolves to the API null response.
 */
export const getWebsocketAuthToken = createClientOnlyFn(() => {
  return authFetcher.get<{ token: string }>("/ws/token");
});
