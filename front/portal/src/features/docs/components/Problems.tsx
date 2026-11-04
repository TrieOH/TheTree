import { Badge } from "@trieoh/ui-solid";
import { For, Show } from "solid-js";
import problems from "virtual:problems";
import type { ServiceId } from "@/features/catalog/services";

/**
 * `<Problems service="identityx" />`: every RFC 9457 problem the service can
 * answer, from api/problems.yml and api/<svc>/problems.yml at build time.
 * Slugs aren't links: `{API_URL}/problems/{slug}` would wake the API.
 */
export function Problems(props: { service: ServiceId }) {
  const list = () => problems[props.service] ?? [];
  return (
    <div class="not-prose my-5 overflow-x-auto">
      <table class="w-full text-left text-[13px]">
        <thead class="text-[12px] text-muted-foreground">
          <tr class="border-b border-border">
            <th class="py-2 pr-3 font-normal">Status</th>
            <th class="py-2 pr-3 font-normal">Slug</th>
            <th class="py-2 font-normal">Meaning</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <For each={list()}>
            {(problem) => (
              <tr class="align-top">
                <td class="py-2 pr-3 tabular-nums">{problem.status}</td>
                <td class="py-2 pr-3">
                  <code class="font-mono text-[12px]">{problem.slug}</code>
                  <Show when={problem.shared}>
                    <Badge variant="outline" class="ml-2">shared</Badge>
                  </Show>
                </td>
                <td class="py-2">
                  <span class="font-medium">{problem.title}</span>
                  <Show when={problem.description}>
                    <span class="text-muted-foreground"> — {problem.description}</span>
                  </Show>
                  <Show when={problem.members.length > 0}>
                    <span class="mt-0.5 block text-[12px] text-muted-foreground">
                      adds <For each={problem.members}>{(member, i) => <>{i() > 0 ? ", " : ""}<code>{member}</code></>}</For>
                    </span>
                  </Show>
                </td>
              </tr>
            )}
          </For>
        </tbody>
      </table>
    </div>
  );
}
