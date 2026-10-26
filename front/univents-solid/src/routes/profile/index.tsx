import { createFileRoute } from '@tanstack/solid-router'
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ownProfileQueryOptions, cacheProfile } from "@/features/profile/api";
import type { ActorProfile } from "@/features/profile/model/profile-data";
import { PROFILE_TABS } from "@/features/profile/model/profile-tabs";
import { Schema, createSchema } from "@/shared/lib/effect-schema";
import { appQueryClient } from "@/shared/lib/query-client";
import { peekQuery } from "@/shared/lib/route-data";

const profileSearchSchema = createSchema(
  Schema.Struct({
    tab: Schema.optional(Schema.Literals(PROFILE_TABS)),
  }),
);

export const Route = createFileRoute('/profile/')({
  validateSearch: (search) => {
    const parsed = profileSearchSchema.parse(search);
    return { tab: parsed.tab ?? "about" };
  },
  beforeLoad: async (args) => {
    requireAuth(args);
    const sessionAuth = args.context.session?.service;
    const actorId = sessionAuth?.profile()?.id;

    if (actorId && sessionAuth) {
      // The handle only makes the URL prettier: the profile page resolves an
      // actor id too and canonicalizes the URL itself, so a slow profile
      // service must not hold up (or fail) this redirect.
      const profile = await peekQuery<ActorProfile | null>(
        appQueryClient,
        ownProfileQueryOptions(actorId, sessionAuth),
      );
      const profileId = profile?.handle ?? actorId;
      if (profile) {
        cacheProfile(appQueryClient, { success: true, data: profile }, profileId);
      }

      throw Route.redirect({
        to: "/profile/$actorId",
        params: { actorId: profileId },
        search: (prev) => prev,
        replace: true,
      });
    }

    throw Route.redirect({
      to: "/auth",
      search: { redirect: "/profile" },
      replace: true,
    });
  },
});
