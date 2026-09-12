import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Same contract as `cn` in `@trieoh/ui-react`: conditional classes, later
 * Tailwind utilities winning over earlier ones. Framework-free on purpose, so
 * both packages can share it (and so a `ui-core` can take it later without a
 * rewrite).
 */
export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}
