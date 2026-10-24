import { createFileRoute } from '@tanstack/solid-router'
import { z } from "zod";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ownProfileQueryOptions, cacheProfile } from "@/features/profile/api";
import { PROFILE_TABS } from "@/features/profile/model/profile-tabs";
import { appQueryClient } from "@/shared/lib/query-client";

const profileSearchSchema = z.object({
  tab: z.enum(PROFILE_TABS).catch("about"),
});

export const Route = createFileRoute('/profile/')({
  validateSearch: profileSearchSchema,
  beforeLoad: async (args) => {
    requireAuth(args);
    const sessionAuth = args.context.session?.service;
    const actorId = sessionAuth?.profile()?.id;

    if (actorId && sessionAuth) {
      const profile = await appQueryClient.fetchQuery(
        ownProfileQueryOptions(actorId, sessionAuth),
      );
      const profileId = profile?.handle ?? actorId;
      if (profile) {
        cacheProfile(appQueryClient, { success: true, data: profile }, profileId);
      }

      throw Route.redirect({
        to: "/profile/$actorId",
        params: { actorId: profileId },
        search: args.search,
      });
    }
  },
})
