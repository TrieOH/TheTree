import { createQueryClient, TanStackQueryProvider } from "@trieoh/front-core"

export function getContext() {
  return { queryClient: createQueryClient() }
}

export const Provider = TanStackQueryProvider

