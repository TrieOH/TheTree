import { createSignal } from "solid-js";

// What's typed into the try-it panel, kept in memory for this tab only:
// credentials are never written to storage. The base URL and credentials
// are shared by every endpoint; parameters and bodies are per endpoint.

const [baseUrl, setBaseUrl] = createSignal("", { name: "try-base-url" });
const [credentials, setCredentials] = createSignal<Record<string, string>>({}, { name: "try-credentials" });
const [params, setParams] = createSignal<Record<string, Record<string, string>>>({}, { name: "try-params" });
const [bodies, setBodies] = createSignal<Record<string, string>>({}, { name: "try-bodies" });

export const tryState = {
  baseUrl,
  setBaseUrl,
  credential: (scheme: string) => credentials()[scheme] ?? "",
  setCredential: (scheme: string, value: string) => setCredentials((all) => ({ ...all, [scheme]: value })),
  credentials,
  params: (op: string) => params()[op] ?? {},
  setParam: (op: string, name: string, value: string) =>
    setParams((all) => ({ ...all, [op]: { ...all[op], [name]: value } })),
  body: (op: string): string | null => bodies()[op] ?? null,
  setBody: (op: string, value: string) => setBodies((all) => ({ ...all, [op]: value })),
};
