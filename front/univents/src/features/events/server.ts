import { createServerFn } from "@tanstack/react-start";
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";
import { env } from "@/env";
import { preventResponseCaching } from "@/shared/lib/http-cache";

const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.IDENTITYX_ACCESS_API_KEY,
});

export const getActorEmailsServerFn = createServerFn({ method: "GET" })
  .validator((data: { actorIds: string[] }) => data)
  .handler(async ({ data }) => {
    preventResponseCaching();
    const actors = await Promise.all(
      [...new Set(data.actorIds)].map(async (actorId) => {
        const actor = await identityXAccessClient.actors.getById(
          env.VITE_TRIEOH_AUTH_PROJECT_ID,
          actorId,
        );
        if (!actor.success) {
          console.warn("[actor-emails] actor lookup failed", {
            actorId,
            code: actor.code,
            errorId: actor.error_id,
          });
          return null;
        }
        if (!actor.data.email) {
          console.warn("[actor-emails] actor has no email", { actorId });
          return null;
        }
        return [actorId, actor.data.email] as const;
      }),
    );
    return Object.fromEntries(
      actors.filter(
        (actor): actor is readonly [string, string] => actor !== null,
      ),
    );
  });
