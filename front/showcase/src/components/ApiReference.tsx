import { For, Show } from "solid-js";
import apiReference from "virtual:api-reference";
import type { Entry } from "../catalog";

export function ApiReference(props: { entry: Entry }) {
  return (
    <section class="mt-10 min-w-0 border-t border-border pt-8 [&_code]:rounded [&_code]:bg-muted [&_code]:px-1.25 [&_code]:py-0.5 [&_code]:font-mono [&_code]:text-[11px]">
      <h2 class="mb-4 text-lg font-semibold tracking-[-0.4px]">
        API Reference
      </h2>
      <h3 class="mt-6 mb-3 text-base font-semibold">{props.entry.name}</h3>
      <p class="my-3.5 text-xs leading-[1.7] text-muted-foreground">
        {props.entry.description}
      </p>
      <Show
        when={apiReference[props.entry.name]?.length}
        fallback={
          <p class="my-3.5 text-xs leading-[1.7] text-muted-foreground">
            Este componente não recebe propriedades.
          </p>
        }
      >
        <div
          class="mt-6 overflow-x-auto"
          tabindex={0}
          role="region"
          aria-label={`API de ${props.entry.name}`}
        >
          <table class="w-full border-collapse text-left text-xs">
            <thead>
              <tr class="border-b border-border">
                <th class="py-3 pr-3 font-medium" scope="col">
                  Propriedade
                </th>
                <th class="py-3 pr-3 font-medium" scope="col">
                  Tipo
                </th>
                <th class="py-3 pr-3 font-medium" scope="col">
                  Padrão
                </th>
                <th class="py-3 pr-3 font-medium" scope="col">
                  Obrigatória
                </th>
              </tr>
            </thead>
            <tbody>
              <For each={apiReference[props.entry.name]}>
                {(property) => (
                  <tr class="border-b border-border">
                    <td class="py-3.5 pr-3 align-top min-w-32.5">
                      <code>{property.name}</code>
                      <Show when={property.description}>
                        <p class="mt-2 text-[11px] leading-[1.6] text-muted-foreground">
                          {property.description}
                        </p>
                      </Show>
                    </td>
                    <td class="py-3.5 pr-3 align-top min-w-40 max-w-90 wrap-anywhere">
                      <code>{property.type}</code>
                    </td>
                    <td class="py-3.5 pr-3 align-top min-w-22.5">
                      <Show when={property.defaultValue} fallback="—">
                        {(value) => <code>{value()}</code>}
                      </Show>
                    </td>
                    <td class="py-3.5 pr-3 align-top min-w-22.5">
                      {property.required ? "Sim" : "Não"}
                    </td>
                  </tr>
                )}
              </For>
            </tbody>
          </table>
        </div>
      </Show>
    </section>
  );
}
