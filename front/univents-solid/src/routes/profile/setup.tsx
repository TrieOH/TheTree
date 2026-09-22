import type { JSX } from "@solidjs/web";
import { createFileRoute } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { Show, createMemo, createSignal } from "solid-js";
import SaveIcon from "~icons/lucide/save";
import UserPenIcon from "~icons/lucide/user-pen";
import z from "zod";

import { requireAuth } from "@/features/auths/lib/route-guard";
import { useProfileSetupMutation } from "@/features/profile/api/mutations";
import { ProfileImageInput } from "@/features/profile/ui/ProfileImageInput";
import { toast } from "@/shared/ui/toast";

const Save = SaveIcon as unknown as (props: { class?: string }) => JSX.Element;
const UserPen = UserPenIcon as unknown as (props: {
  class?: string;
}) => JSX.Element;

export const Route = createFileRoute("/profile/setup")({
  beforeLoad: requireAuth,
  validateSearch: z.object({ returnTo: z.string().optional() }),
  component: ProfileSetup,
});

function ProfileSetup(): JSX.Element {
  const { auth } = useAuth();
  const navigate = Route.useNavigate();
  const search = Route.useSearch();
  const [photo, setPhoto] = createSignal<File>();
  const [handle, setHandle] = createSignal("");
  const [legalName, setLegalName] = createSignal("");
  const [preferredName, setPreferredName] = createSignal("");
  const [saving, setSaving] = createSignal(false);
  const setupMutation = useProfileSetupMutation();
  const complete = createMemo(() =>
    Boolean(handle().trim() && legalName().trim()),
  );

  const save = async (personalizeMore: boolean) => {
    const actorId = auth.profile()?.id;
    const submittedHandle = handle().trim().replace(/^@/, "");
    if (!actorId || !complete()) return;
    if (
      /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(submittedHandle)
    ) {
      toast.error("Escolha um nome de usuário que não seja um UUID.");
      return;
    }
    setSaving(true);
    try {
      await setupMutation.mutateAsync({
        actorId,
        auth,
        handle: submittedHandle,
        legalName: legalName().trim(),
        preferredName: preferredName().trim(),
        photo: photo(),
      });
      toast.success("Perfil criado");
      if (search().returnTo && !personalizeMore) {
        window.location.assign(search().returnTo!);
      } else {
        void navigate(
          personalizeMore
            ? { to: "/profile/edit", search: { returnTo: search().returnTo } }
            : { to: "/profile", search: { tab: "about" } },
        );
      }
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : "Não foi possível salvar o perfil.",
      );
    } finally {
      setSaving(false);
    }
  };

  return (
    <main class="mx-auto flex min-h-dvh w-full max-w-xl flex-col justify-center px-5 py-12 pb-32 sm:px-8">
      <div class="mb-8 text-center">
        <h1 class="text-3xl font-bold">Configure seu perfil</h1>
        <p class="mt-2 text-sm text-muted-foreground">
          Preencha o básico para começar a usar o Univents.
        </p>
      </div>
      <form
        class="flex flex-col gap-5"
        onSubmit={(event) => {
          event.preventDefault();
          void save(false);
        }}
      >
        <div class="flex justify-center py-1">
          <ProfileImageInput
            label="Foto do perfil"
            variant="avatar"
            onSelect={setPhoto}
          />
        </div>
        <SetupField label="Nome de usuário" htmlFor="setup-handle">
          <input
            id="setup-handle"
            required
            pattern="[^\\s/]+"
            autocomplete="username"
            placeholder="seu username"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
            value={handle()}
            onInput={(event) => setHandle(event.currentTarget.value)}
          />
        </SetupField>
        <SetupField label="Nome civil" htmlFor="setup-legal-name">
          <input
            id="setup-legal-name"
            required
            autocomplete="name"
            placeholder="Seu nome completo"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
            value={legalName()}
            onInput={(event) => setLegalName(event.currentTarget.value)}
          />
        </SetupField>
        <SetupField
          label="Nome social (opcional)"
          htmlFor="setup-preferred-name"
        >
          <input
            id="setup-preferred-name"
            placeholder="Como prefere ser chamado"
            class="w-full rounded-md border border-input bg-background px-3 py-2"
            value={preferredName()}
            onInput={(event) => setPreferredName(event.currentTarget.value)}
          />
        </SetupField>
        <div class="mt-3 flex flex-col gap-3">
          <button
            type="submit"
            disabled={!complete() || saving()}
            class="flex h-12 items-center justify-center gap-2 rounded-md bg-primary text-base text-primary-foreground disabled:opacity-50"
          >
            <Save class="size-4" />
            {saving() ? "Salvando…" : "Salvar"}
          </button>
          <Show when={complete()}>
            <button
              type="button"
              disabled={saving()}
              class="flex h-12 items-center justify-center gap-2 rounded-md border border-input text-base disabled:opacity-50"
              onClick={() => void save(true)}
            >
              <UserPen class="size-5" />
              Personalizar mais
            </button>
          </Show>
        </div>
      </form>
    </main>
  );
}

function SetupField(props: {
  label: string;
  htmlFor: string;
  children: JSX.Element;
}): JSX.Element {
  return (
    <label class="space-y-2 text-sm font-medium" for={props.htmlFor}>
      {props.label}
      {props.children}
    </label>
  );
}
