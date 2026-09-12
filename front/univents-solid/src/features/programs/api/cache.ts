import type { QueryClient } from "@trieoh/front-core-solid";
import { programKeys } from "./query-keys";

export function invalidateParticipationCache(queryClient: QueryClient, editionId: string) {
  void queryClient.invalidateQueries({ queryKey: programKeys.mine(editionId) });
  void queryClient.invalidateQueries({ queryKey: ["programs", "page"] });
}
