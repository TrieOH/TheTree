import { createServerAuth } from "@soramux/node-auth-sdk/server";
import { env } from "@/env";

export const serverAuth = createServerAuth(env.VITE_AUTH_API_URL);
