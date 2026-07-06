/**
 * Truncate a string from both sides, adding ellipsis in the middle.
 * Example: "abc...xyz" (for a 6-character string with start=3, end=3)
 */
export function truncateString(str: string, start: number, end: number): string {
  if (str.length <= start + end) return str
  return `${str.slice(0, start)}...${str.slice(str.length - end)}`
}

/**
 * Mask the middle of a string with a custom mask character.
 * Example: "sk_••••••••abc" (for a secret key)
 */
export function maskStringMiddle(
  str: string,
  start: number,
  end: number,
  maskChar = "•••",
): string {
  if (str.length <= start + end) return str

  const hiddenLength = str.length - (start + end)

  return (
    str.slice(0, start) +
    maskChar.repeat(Math.max(3, Math.min(hiddenLength, 8))) +
    str.slice(str.length - end)
  )
}
