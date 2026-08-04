import { createServerFn } from "@tanstack/react-start";
import { createTanStackIdentityXBff } from "@trieoh/front-core/auth/tanstack/server";
import { z } from "zod";
import { env } from "@/env";

const bff = createTanStackIdentityXBff({
  identityX: {
    baseURL: env.VITE_AUTH_API_URL,
    projectId: env.VITE_TRIEOH_AUTH_PROJECT_ID,
  },
  session: {
    password: env.AUTH_SESSION_PASSWORD,
    name: "payssage-auth",
    secure: import.meta.env.PROD,
  },
  apiBaseURL: env.VITE_API_URL,
});

export const loginServerFn = createServerFn({ method: "POST" })
  .validator(z.object({ email: z.email(), password: z.string().min(1) }))
  .handler(({ data }) => bff.login(data.email, data.password));

export const logoutServerFn = createServerFn({ method: "POST" }).handler(() =>
  bff.logout(),
);

export const refreshServerFn = createServerFn({ method: "POST" }).handler(() =>
  bff.refresh(),
);

export const restoreSessionServerFn = createServerFn({ method: "GET" }).handler(
  () => bff.restore(),
);

export const authenticatedProxyServerFn = createServerFn({ method: "POST" })
  .validator(
    z.object({
      path: z
        .string()
        .startsWith("/")
        .refine((path) => !path.startsWith("//")),
      method: z.enum(["GET", "POST", "PUT", "PATCH", "DELETE"]).optional(),
      body: z.json().optional(),
      headers: z.record(z.string(), z.string()).optional(),
    }),
  )
  .handler(({ data }) => bff.request(data));
