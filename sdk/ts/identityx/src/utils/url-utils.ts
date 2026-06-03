export function joinUrl(base: string, path: string): string {
  if (path.startsWith("http")) return path;
  
  const cleanBase = base.replace(/\/$/, "");
  const cleanPath = path.replace(/^\//, "");
  
  return `${cleanBase}/${cleanPath}`;
}
