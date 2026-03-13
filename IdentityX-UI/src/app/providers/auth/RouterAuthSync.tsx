import { useRouter } from "@tanstack/react-router"
import { useAuth } from "@soramux/node-auth-sdk/react"
import { useEffect, useRef, useState } from "react"

export function AuthSynchronizer({ children }: { children: React.ReactNode }) {
  const auth = useAuth()
  const router = useRouter()

  const authRef = useRef(auth)
  const routerRef = useRef(router)

  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    authRef.current = auth
    routerRef.current = router
  }, [auth, router])

  useEffect(() => {
    let mounted = true

    const sync = async () => {
      const r = routerRef.current
      const a = authRef.current

      const currentAuth = r.options.context?.auth
      if (currentAuth?.isAuthenticated === a.isAuthenticated) {
        if (mounted) setIsLoading(false)
        return
      }

      try {
        r.update({ context: { ...r.options.context, auth: a } })
        await r.invalidate()
      } finally {
        if (mounted) setIsLoading(false)
      }
    }

    sync()
    return () => { mounted = false }
  }, [])

  if (isLoading) return null
  return <>{children}</>
}
