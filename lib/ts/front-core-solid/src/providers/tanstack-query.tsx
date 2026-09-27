import { MutationObserver, QueryClient, QueryObserver, type MutationObserverOptions, type QueryObserverOptions } from "@tanstack/query-core"
import { createComponent, createContext, createEffect, createSignal, onCleanup, useContext, untrack, type Accessor, type ParentProps } from "solid-js"

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

function retryQuery(failureCount: number, error: unknown) {
  const value = error as {
    code?: number
    envelope?: { code?: number }
  }
  const code = value?.envelope?.code ?? value?.code
  return !(typeof code === "number" && code >= 400 && code < 500) && failureCount < 3
}

export function createQueryClient(config?: QueryClientConfig) {
  return new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: config?.staleTime ?? 1000 * 60 * 5,
        retry: config?.maxRetries === undefined ? retryQuery : config.maxRetries,
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
  const observer = untrack(() => new MutationObserver<TData, TError, TVariables>(client, options))
  const [state, setState] = createSignal(untrack(() => observer.getCurrentResult()), { ownedWrite: true })
  const unsubscribe = untrack(() => observer.subscribe((result) => setState(() => result)))
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
  const isFn = typeof options === "function"
  const observer = untrack(() => {
    const initialOptions = isFn ? (options as Accessor<QueryObserverOptions<TData, TError>>)() : options
    return new QueryObserver<TData, TError>(client, initialOptions)
  })
  const [state, setState] = createSignal(untrack(() => observer.getCurrentResult()), { ownedWrite: true })
  const unsubscribe = untrack(() => observer.subscribe((result) => setState(() => result)))

  if (isFn) {
    createEffect(
      () => (options as Accessor<QueryObserverOptions<TData, TError>>)(),
      (next) => untrack(() => observer.setOptions(next)),
    )
  }

  onCleanup(() => unsubscribe())
  return state
}
