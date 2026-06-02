export function obfuscate(data: unknown): string {
  try {
    const jsonStr = encodeURIComponent(JSON.stringify(data));
    const reversed = jsonStr.split("").reverse().join("");
    return "tr_" + btoa(reversed);
  } catch {
    return "";
  }
}

export function deobfuscate<T>(obfuscated: string | null): T | null {
  if (!obfuscated) return null;
  try {
    if (!obfuscated.startsWith("tr_")) return null;
    const cleanBase64 = obfuscated.replace(/^tr_/, "");
    const reversed = atob(cleanBase64);
    const jsonStr = reversed.split("").reverse().join("");
    return JSON.parse(decodeURIComponent(jsonStr));
  } catch {
    return null;
  }
}
