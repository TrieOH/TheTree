import { createServerFn } from "@tanstack/react-start";
import { createTanStackIdentityXBff } from "@trieoh/front-core/auth/tanstack/server";
import { z } from "zod";
import { env } from "@/env";

const credentialsSchema = z.object({
  email: z.email(),
  password: z.string().min(1),
});

const providerSchema = z.enum(["github", "google"]);

const bff = createTanStackIdentityXBff({
  identityX: { baseURL: env.VITE_API_URL },
  apiBaseURL: env.VITE_API_URL,
  session: {
    password: env.AUTH_SESSION_PASSWORD,
    name: "identityx-auth",
    secure: import.meta.env.PROD,
  },
});

export const isSetupDoneServerFn = createServerFn({ method: "GET" })
  .handler(() => bff.isSetupDone());

export const setupServerFn = createServerFn({ method: "POST" })
  .validator(credentialsSchema)
  .handler(({ data }) => bff.setup(data.email, data.password));

export const loginServerFn = createServerFn({ method: "POST" })
  .validator(credentialsSchema)
  .handler(({ data }) => bff.login(data.email, data.password));

export const registerServerFn = createServerFn({ method: "POST" })
  .validator(credentialsSchema)
  .handler(({ data }) => bff.register(data.email, data.password));

export const loginWithProviderServerFn = createServerFn({ method: "POST" })
  .validator(z.object({ provider: providerSchema }))
  .handler(({ data }) => bff.loginWithProvider(data.provider));

export const completeProviderLoginServerFn = createServerFn({ method: "POST" })
  .validator(z.object({ provider: providerSchema, code: z.string().min(1) }))
  .handler(({ data }) => bff.completeProviderLogin(data.provider, data.code));

export const logoutServerFn = createServerFn({ method: "POST" })
  .handler(() => bff.logout());

export const refreshServerFn = createServerFn({ method: "POST" })
  .handler(() => bff.refresh());

export const restoreSessionServerFn = createServerFn({ method: "GET" })
  .handler(() => bff.restore());

export const authenticatedProxyServerFn = createServerFn({ method: "POST" })
  .validator(z.object({
    path: z.string().startsWith("/").refine((path) => !path.startsWith("//")),
    method: z.enum(["GET", "POST", "PUT", "PATCH", "DELETE"]).optional(),
    body: z.json().optional(),
    headers: z.record(z.string(), z.string()).optional(),
  }))
  .handler(({ data }) => bff.request(data));
