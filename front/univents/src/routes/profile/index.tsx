import { createFileRoute, redirect } from "@tanstack/react-router";
import { z } from "zod";
import { requireAuth } from "@/features/auths/lib/route-guard";

const profileSearchSchema = z.object({
  tab: z.enum(["about", "badges", "certificates", "purchases"]).catch("about"),
});

export const Route = createFileRoute("/profile/")({
  validateSearch: profileSearchSchema,
  beforeLoad: async (args) => {
    requireAuth(args);
    const sessionAuth = args.context.auth?.auth;
    const actorId = sessionAuth?.profile()?.id;

    if (actorId) {
      const response = await sessionAuth?.getActorProfile(actorId);
      const profileId =
        response?.success && response.data?.handle
          ? response.data.handle
          : actorId;

      throw redirect({
        to: "/profile/$actorId",
        params: { actorId: profileId },
        search: args.search,
      });
    }
  },
});
