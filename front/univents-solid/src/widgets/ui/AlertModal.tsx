import type { JSX } from "@solidjs/web";
import { Show } from "solid-js";
import { Button, Dialog } from "@trieoh/ui-solid";

export interface AlertModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  onConfirm: () => void | Promise<void>;
  variant?: "default" | "destructive";
  loading?: boolean;
}

export function AlertModal(props: AlertModalProps): JSX.Element {
  const confirmLabel = () => props.confirmLabel ?? "Confirmar";
  const cancelLabel = () => props.cancelLabel ?? "Cancelar";
  const variant = () => props.variant ?? "default";

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title={props.title}
      description={props.description}
      footer={
        <div class="flex w-full items-center justify-end gap-2 pt-2">
          <Button
            variant="outline"
            disabled={props.loading}
            onClick={() => props.onOpenChange(false)}
          >
            {cancelLabel()}
          </Button>
          <Button
            variant={variant() === "destructive" ? "destructive" : "default"}
            disabled={props.loading}
            onClick={() => props.onConfirm()}
          >
            <Show when={props.loading} fallback={confirmLabel()}>
              Processando...
            </Show>
          </Button>
        </div>
      }
    />
  );
}
