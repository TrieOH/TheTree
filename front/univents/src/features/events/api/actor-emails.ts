import { createServerFn } from "@tanstack/react-start";
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";
import { env } from "@/env";

const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.IDENTITYX_ACCESS_API_KEY,
});

export const getActorEmailsServerFn = createServerFn({ method: "GET" })
  .validator((data: { actorIds: string[] }) => data)
  .handler(async ({ data }) => {
    const actors = await identityXAccessClient.actors.list(
      env.VITE_TRIEOH_AUTH_PROJECT_ID,
    );

    if (!actors.success) return {};

    const requestedIds = new Set(data.actorIds);
    return Object.fromEntries(
      actors.data
        .filter((actor) => requestedIds.has(actor.id) && actor.email)
        .map((actor) => [actor.id, actor.email as string]),
    );
  });
