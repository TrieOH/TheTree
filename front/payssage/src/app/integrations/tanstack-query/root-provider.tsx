import { TanStackQueryProvider, createQueryClient } from "@trieoh/front-core"
import type { ReactNode } from "react"

export function getContext() {
  return { queryClient: createQueryClient() }
}

export function Provider({ children }: { children: ReactNode }) {
  const { queryClient } = getContext()

  return <TanStackQueryProvider queryClient={queryClient}>{children}</TanStackQueryProvider>
}
