import { spicedb } from "@soramux/node-perm-sdk";
import { env } from "#/env";

export const serverPerm = spicedb.permission({
  url: env.TRIEOH_AUTHZED_URL,
  token: env.TRIEOH_AUTHZED_TOKEN
})

export const serverRelationship = spicedb.relationship({
  url: env.TRIEOH_AUTHZED_URL,
  token: env.TRIEOH_AUTHZED_TOKEN
})
