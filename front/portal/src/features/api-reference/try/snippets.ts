import { goNative } from "@scalar/snippetz/plugins/go/native";
import { jsFetch } from "@scalar/snippetz/plugins/js/fetch";
import { pythonRequests } from "@scalar/snippetz/plugins/python/requests";
import { shellCurl } from "@scalar/snippetz/plugins/shell/curl";
import type { HarRequest } from "./request";

// Snippets come from @scalar/snippetz; only these four plugins are bundled.
export const LANGUAGES = [
  { id: "curl", label: "cURL", filename: "request.sh", plugin: shellCurl },
  { id: "js", label: "JavaScript", filename: "request.js", plugin: jsFetch },
  { id: "python", label: "Python", filename: "request.py", plugin: pythonRequests },
  { id: "go", label: "Go", filename: "main.go", plugin: goNative },
] as const;

export type LanguageId = (typeof LANGUAGES)[number]["id"];

export function snippet(language: LanguageId, request: HarRequest): string {
  const entry = LANGUAGES.find((lang) => lang.id === language) ?? LANGUAGES[0];
  return entry.plugin.generate(request as never);
}
