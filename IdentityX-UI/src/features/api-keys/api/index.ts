import { authFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";

export const rotateApiKey = createClientOnlyFn((project_id: string) => {
  return authFetcher.post<string>(`/projects/${project_id}/api-keys/rotate`);
});

export const revokeApiKey = createClientOnlyFn((project_id: string) => {
  return authFetcher.delete<null>(`/projects/${project_id}/api-keys`);
});