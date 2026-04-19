import { spicedb } from "@soramux/node-perm-sdk";
// import { createServerAuth } from "@soramux/node-auth-sdk/server";
import { env } from "@/env";

// export const serverAuth = createServerAuth(env.VITE_AUTH_API_URL);

export const serverPerm = spicedb.permission({
  url: env.TRIEOH_AUTHZED_URL,
  token: env.TRIEOH_AUTHZED_TOKEN
})
