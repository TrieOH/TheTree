import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import type { OccurrenceI } from "../model";
import { OccurrenceDrawPage } from "./OccurrenceDrawPage";

export interface OccurrenceDrawModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  occurrence: OccurrenceI | null;
  occurrenceId?: string;
  programKind?: "activity" | "checkpoint";
  programName?: string;
}

export function OccurrenceDrawModal(props: OccurrenceDrawModalProps): JSX.Element {
  const occId = () => props.occurrence?.id ?? props.occurrenceId;

  return (
    <Show when={props.open && occId()}>
      {(id) => (
        <OccurrenceDrawPage
          occurrenceId={id()}
          occurrence={props.occurrence}
          programName={props.programName ?? "Atividade"}
          programKind={props.programKind ?? "activity"}
          onBack={() => props.onOpenChange(false)}
        />
      )}
    </Show>
  );
}
