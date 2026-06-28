// import { tanstackQueryFetcher } from "@/shared/lib/api/fetch";
// import { queryOptions } from "@tanstack/react-query";
// import { createClientOnlyFn } from "@tanstack/react-start";
// import type { ActorI } from "../model";


// /**
//  * Fetches all actors from the server.
//  * @param orgId - The organization ID to filter actors by (optional).
//  * @returns A promise that resolves to an array of ActorI objects.
//  */
// export const getActorsFn = createClientOnlyFn(async (orgId?: string) => {
//   if (orgId)
//     return await tanstackQueryFetcher<ActorI[]>(`/organizations/${orgId}/actors`);
//   return await tanstackQueryFetcher<ActorI[]>("/actors");
// });

// /**
//  * Query options for fetching actors, compatible with React Query's useQuery hook.
//  * @param orgId - The organization ID to filter actors by (optional).
//  * @returns An object containing the query key and query function for fetching actors.
//  */
// export const allActorsQueryOptions = (orgId?: string) => {
//   return queryOptions({
//     queryKey: ["organizations", orgId, "actors"],
//     queryFn: () => getActorsFn(orgId),
//   });
// };