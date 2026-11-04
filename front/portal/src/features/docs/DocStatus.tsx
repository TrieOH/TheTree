import { Badge } from "@trieoh/ui-solid";
import { Show } from "solid-js";
import { statusWord } from "./catalog";

export function DocStatus(props: { status?: string }) {
  const word = () => statusWord(props.status);
  return (
    <Show when={word()}>
      {(status) => (
        <Badge variant={status() === "accepted" ? "accent" : "outline"} class="capitalize">
          {status()}
        </Badge>
      )}
    </Show>
  );
}
