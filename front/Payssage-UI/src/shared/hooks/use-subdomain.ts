import { useEffect, useState } from 'react'
import { getSubdomain } from '../lib/subdomain'

export function useSubdomain() {
  const [subdomain, setSubdomain] = useState<string | null>(null)

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const detected = getSubdomain(window.location.hostname)
      setSubdomain(detected)
    }
  }, [])

  return subdomain
}
