export function getSubdomain(hostname: string): string | null {
  const parts = hostname.split('.')

  // localhost
  if (hostname === 'localhost') return null

  // slug.localhost
  if (parts.length === 2 && parts[1] === 'localhost') return parts[0]

  // IP
  if (/^(?:\d{1,3}\.){3}\d{1,3}$/.test(hostname)) return null

  // Default: subdomain.mydomain.com
  if (parts.length > 2) {
    // Return everything except the last two parts (domain and TLD)
    return parts.slice(0, -2).join('.')
  }

  return null
}
