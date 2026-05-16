import { tanstackQueryFetcher } from "@/shared/lib/api/fetch";
import type { User } from "../model/types";
import { queryOptions } from "@tanstack/react-query";

/**
 * Fetches all users from the server.
 * @returns A promise that resolves to an array of User objects.
 */
export const getUsersFn = (async ({ queryKey,}: { queryKey: ["users", string] }) => {
  const [, projectId] = queryKey;
  try {
    return await tanstackQueryFetcher<User[]>(
      `/projects/${projectId}/users`
    );
  } catch {
    return [] as User[];
  }
});

export const usersQueryOptions = (project_id: string) => {
  return queryOptions({
    queryKey: ['users', project_id],
    queryFn: getUsersFn,
    enabled: !!project_id
  })
}
