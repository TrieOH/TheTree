import { authFetcher } from "@/shared/lib/api/fetch";
import { createClientOnlyFn } from "@tanstack/react-start";

export const rotateApiKey = createClientOnlyFn((project_id: string) => {
  return authFetcher<{api_key: string}>(`/projects/${project_id}/api-keys/rotate`, {
    method: "POST",
  });
});

export const revokeApiKey = createClientOnlyFn((project_id: string) => {
  return authFetcher<null>(`/projects/${project_id}/api-keys`, {
    method: "DELETE",
  });
});