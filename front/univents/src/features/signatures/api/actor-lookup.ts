import { createServerFn } from "@tanstack/react-start";
import { createIdentityXAccessClient } from "@trieoh/identityx-access-sdk-ts";
import { z } from "zod";
import { env } from "@/env";

const identityXAccessClient = createIdentityXAccessClient({
  baseURL: env.VITE_AUTH_API_URL,
  apiKey: env.IDENTITYX_ACCESS_API_KEY,
});

export const findActorIdByEmailServerFn = createServerFn({ method: "GET" })
  .validator(z.object({ email: z.email() }))
  .handler(async ({ data }) => {
    const result = await identityXAccessClient.actors.getByEmail(
      env.VITE_TRIEOH_AUTH_PROJECT_ID,
      data.email,
    );

    return result.success ? result.data.id : undefined;
  });
