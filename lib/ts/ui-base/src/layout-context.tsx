import {
  createContext,
  useContext,
  useLayoutEffect,
} from 'react'
import type { ReactNode } from 'react';

interface LayoutContextValue {
  /** Replace the header slot content. Pass null to clear. */
  setHeader: (node: ReactNode) => void
}

export const LayoutContext =
  createContext<LayoutContextValue | null>(null)

/**
 * Call inside any child page to inject content into the layout header slot.
 *
 * @example
 * function FormsPage() {
 *   useLayoutHeader(
 *     <div className="flex items-center justify-between">
 *       <h1 className="text-xl font-semibold">Forms</h1>
 *       <Button>New Form</Button>
 *     </div>
 *   )
 *   return <FormsTable />
 * }
 */
export function useLayoutHeader(node: ReactNode) {
  const ctx = useContext(LayoutContext)
  if (!ctx) throw new Error('useLayoutHeader must be used inside Layout')

  const { setHeader } = ctx

  useLayoutEffect(() => {
    setHeader(node)

    return () => {
      setHeader(null)
    }
  }, [node, setHeader])
}
