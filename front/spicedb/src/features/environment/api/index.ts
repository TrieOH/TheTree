import { env } from "#/env"
import { createServerFn } from "@tanstack/react-start"

export const getEnvNames = createServerFn({
  method: "GET",
}).handler(() => {
  const envs = env.TRIEOH_AUTHZED_ENVIRONMENTS
  return envs.map(e => e.name)
})