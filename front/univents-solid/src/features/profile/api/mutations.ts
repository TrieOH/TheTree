import { useQueryClient } from "@trieoh/front-core-solid";
import type { ProfileData } from "@trieoh/identityx-sdk-ts-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import { Effect } from "effect";

import { preprocessImageUpload } from "@/features/storage/api/index";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { invalidateProfileCaches } from "./cache";

export type ProfileSetupInput = {
  actorId: string;
  auth: {
    upsertActorProfile: (
      actorId: string,
      data: {
        handle: string;
        pfp_url: string | null;
        profile: ProfileData;
      },
    ) => Promise<{ success: boolean; message?: string }>;
  };
  handle: string;
  legalName: string;
  preferredName: string;
  photo?: File;
};

export const useProfileSetupMutation = () => {
  const queryClient = useQueryClient();

  return useEffectMutation({
    retryTransient: true,
    mutationEffect: (input: ProfileSetupInput) =>
      apiEffect(async () => {
        const now = new Date().toISOString();
        const pfpUrl = input.photo
          ? await preprocessImageUpload(
            input.photo,
            "profiles/images",
            crypto.randomUUID(),
            256,
          )
          : undefined;
        return { now, pfpUrl };
      }).pipe(
        Effect.flatMap(({ now, pfpUrl }) =>
          apiEffect(() =>
            withSpan("action:profile-setup", async () => {
              const response = await input.auth.upsertActorProfile(
                input.actorId,
                {
                  handle: input.handle,
                  pfp_url: pfpUrl ?? null,
                  profile: {
                    createdAt: now,
                    updatedAt: now,
                    legalName: input.legalName,
                    preferredName: input.preferredName || input.legalName,
                    ...(pfpUrl ? { pfpUrl } : {}),
                  },
                },
              );
              if (!response.success) {
                throw new Error(
                  response.message || "Não foi possível salvar o perfil.",
                );
              }
            }),
          ),
        ),
        Effect.tap(() =>
          Effect.promise(() => invalidateProfileCaches(queryClient)),
        ),
      ),
  });
};
