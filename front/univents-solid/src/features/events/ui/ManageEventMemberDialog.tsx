import type { JSX } from "@solidjs/web";
import { For, Show, createSignal } from "solid-js";

import LoaderCircleIcon from "~icons/lucide/loader-circle";
import ShieldCheckIcon from "~icons/lucide/shield-check";
import UserCogIcon from "~icons/lucide/user-cog";
import UserPlusIcon from "~icons/lucide/user-plus";
import UsersIcon from "~icons/lucide/users";

import { Button, Dialog, Field, Input, cn } from "@trieoh/ui-solid";
import type { EventMemberRole } from "../model/member";

const LoaderCircle = LoaderCircleIcon as unknown as (props: { class?: string }) => JSX.Element;
const ShieldCheck = ShieldCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserCog = UserCogIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserPlus = UserPlusIcon as unknown as (props: { class?: string }) => JSX.Element;
const Users = UsersIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface ManageEventMemberValues {
  email: string;
  role: EventMemberRole;
}

export interface ManageEventMemberDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (values: ManageEventMemberValues) => Promise<boolean>;
}

const ROLE_OPTIONS: Array<{
  value: EventMemberRole;
  label: string;
  description: string;
  icon: (props: { class?: string }) => JSX.Element;
}> = [
  {
    value: "staff",
    label: "Equipe",
    description: "Acesso às operações do dia a dia.",
    icon: Users,
  },
  {
    value: "admin",
    label: "Administrador",
    description: "Pode gerenciar o evento e sua equipe.",
    icon: UserCog,
  },
  {
    value: "owner",
    label: "Proprietário",
    description: "Nível máximo de acesso ao evento.",
    icon: ShieldCheck,
  },
];

export function ManageEventMemberDialog(props: ManageEventMemberDialogProps): JSX.Element {
  const [email, setEmail] = createSignal("");
  const [role, setRole] = createSignal<EventMemberRole>("staff");
  const [submitting, setSubmitting] = createSignal(false);
  const [error, setError] = createSignal<string | null>(null);

  const reset = () => {
    setEmail("");
    setRole("staff");
    setError(null);
  };

  const submit = async (e: SubmitEvent) => {
    e.preventDefault();
    const cleanEmail = email().trim().toLowerCase();

    if (!cleanEmail || !cleanEmail.includes("@")) {
      setError("Informe um e-mail válido.");
      return;
    }

    setError(null);
    setSubmitting(true);
    try {
      const ok = await props.onSubmit({
        email: cleanEmail,
        role: role(),
      });
      if (ok) {
        reset();
        props.onOpenChange(false);
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={props.open}
      onOpenChange={(open) => {
        if (!open) reset();
        props.onOpenChange(open);
      }}
      title={
        <span class="flex items-center gap-2">
          <UserPlus class="size-5 text-primary" />
          <span>Adicionar membro</span>
        </span>
      }
      description="Convide uma pessoa para colaborar na gestão deste evento."
      footer={
        <div class="flex items-center justify-end gap-2 pt-2">
          <Button
            type="button"
            variant="ghost"
            onClick={() => {
              reset();
              props.onOpenChange(false);
            }}
            disabled={submitting()}
          >
            Cancelar
          </Button>
          <Button
            type="submit"
            form="manage-member-form"
            disabled={submitting() || !email().trim()}
            class="gap-2"
          >
            <Show when={submitting()} fallback={<UserPlus class="size-4" />}>
              <LoaderCircle class="size-4 animate-spin" />
            </Show>
            Adicionar membro
          </Button>
        </div>
      }
    >
      <form id="manage-member-form" onSubmit={submit} class="space-y-4 py-2">
        <Field
          label="E-mail do membro"
          hint="O usuário precisa possuir uma conta vinculada a este e-mail."
          error={error() ?? undefined}
        >
          {({ id, describedBy }) => (
            <Input
              id={id}
              aria-describedby={describedBy}
              type="email"
              value={email()}
              onInput={(e) => {
                setEmail(e.currentTarget.value);
                if (error()) setError(null);
              }}
              placeholder="membro@exemplo.com"
              autocomplete="email"
              required
            />
          )}
        </Field>

        <div class="space-y-2">
          <label class="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
            Função no evento
          </label>
          <div class="grid gap-2 sm:grid-cols-3">
            <For each={ROLE_OPTIONS}>
              {(option) => {
                const selected = () => role() === option.value;
                return (
                  <button
                    type="button"
                    aria-pressed={selected() ? "true" : "false"}
                    onClick={() => setRole(option.value)}
                    class={cn(
                      "flex min-h-24 flex-col items-start gap-1.5 rounded-xl border p-3 text-left transition-all",
                      selected()
                        ? "border-primary bg-primary/10 text-primary ring-2 ring-primary/20 shadow-xs"
                        : "border-border bg-card hover:border-primary/40 hover:bg-muted/40 text-foreground",
                    )}
                  >
                    <div class="flex items-center gap-1.5 font-medium text-sm">
                      {option.icon({ class: "size-4 shrink-0" })}
                      <span>{option.label}</span>
                    </div>
                    <span
                      class={cn(
                        "text-xs leading-relaxed",
                        selected() ? "text-primary/80" : "text-muted-foreground",
                      )}
                    >
                      {option.description}
                    </span>
                  </button>
                );
              }}
            </For>
          </div>
        </div>
      </form>
    </Dialog>
  );
}
