import { TanStackQueryProvider, createQueryClient } from "@trieoh/front-core"
import type { ReactNode } from "react"

let context: { queryClient: ReturnType<typeof createQueryClient> } | undefined

export function getContext() {
  if (context) return context

  context = { queryClient: createQueryClient() }

  return context
}

export function Provider({ children }: { children: ReactNode }) {
  const { queryClient } = getContext()

  return <TanStackQueryProvider queryClient={queryClient}>{children}</TanStackQueryProvider>
}
