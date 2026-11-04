import { CodeBlock, Tabs, cn } from "@trieoh/ui-solid";
import { For, Show, createMemo, createSignal } from "solid-js";
import SendIcon from "~icons/lucide/send";
import type { ApiAuth, ApiOperation } from "../model";
import { PLACEHOLDER_BASE, authFor, buildRequest, sendCheck, toFetch, type TryValues } from "../try/request";
import { LANGUAGES, snippet, type LanguageId } from "../try/snippets";
import { tryState } from "../try/state";
import { MethodBadge, PathText } from "./MethodBadge";

const [language, setLanguage] = createSignal<LanguageId>("curl", { name: "try-language" });

const input =
  "h-8 w-full rounded-md border border-border bg-background px-2 font-mono text-[12px] outline-none placeholder:text-muted-foreground/60 focus:border-ring";
const label = "mb-1 flex items-baseline gap-1.5 text-[11.5px] text-muted-foreground";

interface Result {
  status: number | null;
  statusText: string;
  ms: number;
  body: string;
}

function pretty(text: string): string {
  try {
    return JSON.stringify(JSON.parse(text), null, 2);
  } catch {
    return text;
  }
}

/**
 * Fill in an endpoint's request, get it as code, and send it. Nothing is sent
 * until Send is pressed; the base URL starts empty, so the page itself never
 * reaches an API.
 */
export function TryPanel(props: { op: ApiOperation; auth: readonly ApiAuth[]; class?: string }) {
  const [result, setResult] = createSignal<Result | null>(null, { name: "try-result" });
  const [sending, setSending] = createSignal(false, { name: "try-sending" });

  const values = (): TryValues => ({
    baseUrl: tryState.baseUrl(),
    credentials: tryState.credentials(),
    params: tryState.params(props.op.id),
    body: tryState.body(props.op.id),
  });
  const request = createMemo(() => buildRequest(props.op, props.auth, values()));
  const check = createMemo(() => sendCheck(props.op, request(), values()));
  const scheme = () => authFor(props.op, props.auth, values());
  const fields = () => props.op.parameters.filter((param) => param.in !== "cookie");

  const send = async () => {
    if (check().blocker) return;
    const { url, init } = toFetch(request());
    setSending(true);
    const started = performance.now();
    try {
      const response = await fetch(url, init);
      setResult({
        status: response.status,
        statusText: response.statusText,
        ms: Math.round(performance.now() - started),
        body: pretty(await response.text()),
      });
    } catch (error) {
      setResult({
        status: null,
        statusText: "No response",
        ms: Math.round(performance.now() - started),
        body: `${error instanceof Error ? error.message : String(error)}\n\nThe API may be down, or not allow requests from this origin (CORS).`,
      });
    } finally {
      setSending(false);
    }
  };

  return (
    <div class={cn("space-y-4 text-[13px]", props.class)}>
      <div class="flex items-center gap-2">
        <MethodBadge method={props.op.method} />
        <PathText path={props.op.path} />
      </div>

      <div class="space-y-3">
        <div>
          <label class={label} for="try-base">Base URL</label>
          <input
            id="try-base"
            class={input}
            placeholder={`${PLACEHOLDER_BASE}, e.g. http://localhost:8080`}
            value={tryState.baseUrl()}
            onInput={(event) => tryState.setBaseUrl(event.currentTarget.value)}
          />
        </div>

        <Show when={scheme()}>
          {(auth) => (
            <div>
              <label class={label} for={`try-auth-${auth().name}`}>
                {auth().header}
                <span class="font-mono">{auth().prefix.trim()}</span>
              </label>
              <input
                id={`try-auth-${auth().name}`}
                type="password"
                autocomplete="off"
                class={input}
                placeholder={auth().placeholder}
                value={tryState.credential(auth().name)}
                onInput={(event) => tryState.setCredential(auth().name, event.currentTarget.value)}
              />
            </div>
          )}
        </Show>

        <For each={fields()}>
          {(param) => (
            <div>
              <label class={label} for={`try-${props.op.id}-${param.name}`}>
                <span class="font-mono text-foreground">{param.name}</span>
                {param.in}
                <Show when={param.required}>
                  <span class="text-warn">required</span>
                </Show>
              </label>
              <input
                id={`try-${props.op.id}-${param.name}`}
                class={input}
                placeholder={`<${param.name}>`}
                value={tryState.params(props.op.id)[param.name] ?? ""}
                onInput={(event) => tryState.setParam(props.op.id, param.name, event.currentTarget.value)}
              />
            </div>
          )}
        </For>

        <Show when={props.op.requestBody}>
          {(body) => (
            <div>
              <label class={label} for={`try-${props.op.id}-body`}>
                Body <span class="font-mono">{body().contentType}</span>
              </label>
              <textarea
                id={`try-${props.op.id}-body`}
                rows={6}
                spellcheck={false}
                class={cn(input, "h-auto py-1.5")}
                value={tryState.body(props.op.id) ?? body().example ?? ""}
                onInput={(event) => tryState.setBody(props.op.id, event.currentTarget.value)}
              />
            </div>
          )}
        </Show>
      </div>

      <div>
        <Tabs
          ariaLabel="Snippet language"
          items={LANGUAGES.map((lang) => ({ id: lang.id, label: lang.label }))}
          value={language()}
          onChange={(id) => setLanguage(id as LanguageId)}
          size="sm"
        />
        <CodeBlock
          class="mt-2"
          code={snippet(language(), request())}
          filename={LANGUAGES.find((lang) => lang.id === language())!.filename}
        />
      </div>

      <div class="space-y-1.5">
        <button
          type="button"
          onClick={send}
          disabled={!!check().blocker || sending()}
          class="inline-flex h-8 w-full items-center justify-center gap-2 rounded-md bg-primary text-[13px] font-medium text-primary-foreground transition-opacity disabled:opacity-40"
        >
          <SendIcon class="size-3.5" />
          {sending() ? "Sending…" : "Send"}
        </button>
        <p class="text-[11.5px] text-muted-foreground">
          {check().blocker ?? check().hint ?? "Sends from your browser. A sleeping API wakes up to answer."}
        </p>
      </div>

      <Show when={result()}>
        {(res) => (
          <div>
            <p class="mb-1 flex items-baseline gap-2 text-[12px]">
              <span
                class={cn(
                  "font-mono font-medium",
                  res().status === null ? "text-destructive-foreground" : res().status! < 400 ? "text-ok" : "text-warn",
                )}
              >
                {res().status ?? "—"} {res().statusText}
              </span>
              <span class="text-muted-foreground">{res().ms} ms</span>
            </p>
            <CodeBlock code={res().body || "(empty body)"} filename="response" />
          </div>
        )}
      </Show>
    </div>
  );
}
