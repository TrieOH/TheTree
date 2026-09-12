import { MutationObserver, QueryClient, QueryObserver, type MutationObserverOptions, type QueryObserverOptions } from "@tanstack/query-core"
import { createComponent, createContext, createEffect, createSignal, onCleanup, useContext, type Accessor, type ParentProps } from "solid-js"

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

export function useMutation<TData, TError = Error, TVariables = void>(
  options: Omit<MutationObserverOptions<TData, TError, TVariables>, "mutationKey">,
) {
  const client = useQueryClient()
  const observer = new MutationObserver<TData, TError, TVariables>(client, options)
  const [state, setState] = createSignal(observer.getCurrentResult())
  const unsubscribe = observer.subscribe((result) => setState(() => result))
  onCleanup(() => unsubscribe())
  return {
    mutateAsync: (variables: TVariables) => observer.mutate(variables),
    result: state,
  }
}

export function useQuery<TData, TError = Error>(
  options: QueryObserverOptions<TData, TError> | Accessor<QueryObserverOptions<TData, TError>>,
) {
  const client = useQueryClient()
  const getOptions = typeof options === "function" ? options : () => options
  const observer = new QueryObserver<TData, TError>(client, getOptions())
  const [state, setState] = createSignal(observer.getCurrentResult())
  const unsubscribe = observer.subscribe((result) => setState(() => result))
  createEffect(
    getOptions,
    (next) => observer.setOptions(next),
  )
  void observer.refetch()
  onCleanup(() => unsubscribe())
  return state
}
