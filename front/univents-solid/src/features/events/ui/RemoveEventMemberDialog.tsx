import type { JSX } from "@solidjs/web";
import { Show, createEffect, createSignal } from "solid-js";

import LoaderCircleIcon from "~icons/lucide/loader-circle";
import TrashIcon from "~icons/lucide/trash";

import { Button, Dialog, Field, Input } from "@trieoh/ui-solid";
import type { EventMemberWithEmailI } from "../model/member";

const LoaderCircle = LoaderCircleIcon as unknown as (props: { class?: string }) => JSX.Element;
const Trash = TrashIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface RemoveEventMemberDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  member: EventMemberWithEmailI | null;
  onRemove: (userId: string, email: string) => Promise<boolean>;
}

export function RemoveEventMemberDialog(props: RemoveEventMemberDialogProps): JSX.Element {
  const [confirmEmail, setConfirmEmail] = createSignal("");
  const [submitting, setSubmitting] = createSignal(false);

  createEffect(
    () => props.open,
    (open) => {
      if (open) setConfirmEmail("");
    },
  );

  const memberEmail = () => props.member?.email?.trim().toLowerCase() ?? "";

  const canSubmit = () => {
    if (submitting()) return false;
    const required = memberEmail();
    if (required) {
      return confirmEmail().trim().toLowerCase() === required;
    }
    return true;
  };

  const handleRemove = async () => {
    if (!props.member || !canSubmit()) return;

    setSubmitting(true);
    try {
      const emailToSend = memberEmail() || confirmEmail().trim().toLowerCase();
      const ok = await props.onRemove(props.member.user_id, emailToSend);
      if (ok) {
        props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={props.onOpenChange}
      title="Remover membro?"
      description={
        memberEmail()
          ? `Digite o e-mail do membro (${memberEmail()}) para confirmar a remoção da equipe.`
          : "Tem certeza de que deseja remover este membro da equipe deste evento?"
      }
      footer={
        <div class="flex items-center justify-end gap-2 pt-2">
          <Button
            type="button"
            variant="ghost"
            onClick={() => props.onOpenChange(false)}
            disabled={submitting()}
          >
            Cancelar
          </Button>
          <Button
            type="button"
            variant="destructive"
            onClick={handleRemove}
            disabled={!canSubmit()}
            class="gap-2"
          >
            <Show when={submitting()} fallback={<Trash class="size-4" />}>
              <LoaderCircle class="size-4 animate-spin" />
            </Show>
            Remover membro
          </Button>
        </div>
      }
    >
      <Show when={memberEmail()}>
        <div class="py-2">
          <Field label="Confirmar e-mail do membro">
            {({ id, describedBy }) => (
              <Input
                id={id}
                aria-describedby={describedBy}
                type="email"
                value={confirmEmail()}
                onInput={(e) => setConfirmEmail(e.currentTarget.value)}
                placeholder={memberEmail()}
                autocomplete="email"
              />
            )}
          </Field>
        </div>
      </Show>
    </Dialog>
  );
}
