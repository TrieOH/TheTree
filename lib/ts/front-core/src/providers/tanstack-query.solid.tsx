import { QueryClient } from "@tanstack/query-core"
import { createComponent, createContext, useContext, type ParentProps } from "solid-js"

const QueryClientContext = createContext<QueryClient>()

export { QueryClient }

export interface QueryClientConfig {
  staleTime?: number
  maxRetries?: number
  onError?: (error: unknown) => void
}

export class QueryError extends Error {
  readonly envelope: { code?: number; message: string }

  constructor(message: string, code?: number) {
    super(message)
    this.name = "QueryError"
    this.envelope = { message, code }
  }
}

export function queryError(message: string, code?: number) {
  return new QueryError(message, code)
}

export function createQueryClient(config?: QueryClientConfig) {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: config?.staleTime ?? 1000 * 60 * 5,
        retry: config?.maxRetries ?? 3,
      },
      mutations: { onError: config?.onError },
    },
  })
}

export function TanStackQueryProvider(props: {
  client: QueryClient
  children: ParentProps["children"]
}) {
  return createComponent(QueryClientContext, {
    value: props.client,
    get children() {
      return props.children
    },
  })
}

export function useQueryClient() {
  const client = useContext(QueryClientContext)
  if (!client) throw new Error("useQueryClient must be used inside <TanStackQueryProvider>")
  return client
}
