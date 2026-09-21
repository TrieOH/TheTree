import { orvalData } from "@trieoh/api-client";
import { listEditionSignatures } from "@trieoh/univents-api";
import type { SignatureI } from "../model";

export const allSignaturesQueryOptions = (editionId: string) => ({
  queryKey: ["signatures", "by-edition", editionId] as const,
  queryFn: async () => {
    return listEditionSignatures(editionId, { public: true }).then(
      orvalData<SignatureI[]>,
    );
  },
});
