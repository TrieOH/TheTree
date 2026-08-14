import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Save, UserPen } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { profileKeys } from "@/features/profile/api/query-keys";
import {
  isInitialProfileComplete,
  withProfileTimestamps,
} from "@/features/profile/model/profile-data";
import { ProfileImageInput } from "@/features/profile/ui/profile-image-input";
import { preprocessImageUpload } from "@/features/storage/api";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";

export const Route = createFileRoute("/profile/setup")({
  beforeLoad: requireAuth,
  component: InitialProfileSetup,
});

function InitialProfileSetup() {
  const { auth } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [photo, setPhoto] = useState<File>();
  const [handle, setHandle] = useState("");
  const [legalName, setLegalName] = useState("");
  const [preferredName, setPreferredName] = useState("");
  const [saving, setSaving] = useState(false);
  const complete = isInitialProfileComplete({
    handle,
    legalName,
  });

  const save = async (personalizeMore: boolean) => {
    const actorId = auth.profile()?.id;
    const submittedHandle = handle.trim().replace(/^@/, "");
    if (!actorId || !complete) return;
    if (
      /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(submittedHandle)
    ) {
      toast.error("Escolha um nome de usuário que não seja um UUID.");
      return;
    }

    setSaving(true);
    try {
      const pfpUrl = photo
        ? await preprocessImageUpload(
            photo,
            "profiles/images",
            crypto.randomUUID(),
          )
        : undefined;
      const profile = withProfileTimestamps({}, {
        legalName: legalName.trim(),
        preferredName: preferredName.trim() || legalName.trim(),
        ...(pfpUrl ? { pfpUrl } : {}),
      } as ProfileData);
      const response = await auth.upsertActorProfile(actorId, {
        handle: submittedHandle,
        profile,
      });
      if (!response.success) {
        throw new Error(
          response.message || "Não foi possível salvar o perfil.",
        );
      }
      await queryClient.invalidateQueries({
        queryKey: profileKeys.all,
        refetchType: "all",
      });
      toast.success("Perfil criado");
      await navigate({ to: personalizeMore ? "/profile/edit" : "/profile" });
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
    <main className="mx-auto flex min-h-dvh w-full max-w-xl flex-col justify-center px-5 py-12 pb-32 sm:px-8">
      <div className="mb-8 text-center">
        <h1 className="text-3xl font-bold">Configure seu perfil</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Preencha o básico para começar a usar o Univents.
        </p>
      </div>

      <form
        className="flex flex-col gap-5"
        onSubmit={(event) => {
          event.preventDefault();
          void save(false);
        }}
      >
        <div className="flex justify-center py-1">
          <ProfileImageInput
            label="Foto do perfil"
            file={photo}
            onSelect={setPhoto}
            variant="avatar"
            className="rounded-full ring-2 ring-border ring-offset-4 ring-offset-background"
          />
        </div>

        <SetupField label="Nome de usuário" htmlFor="setup-handle">
          <Input
            id="setup-handle"
            required
            value={handle}
            pattern={"[^\\s/]+"}
            autoCapitalize="none"
            autoComplete="username"
            placeholder="seu username"
            onChange={(event) => setHandle(event.target.value)}
          />
        </SetupField>
        <SetupField label="Nome civil" htmlFor="setup-legal-name">
          <Input
            id="setup-legal-name"
            required
            value={legalName}
            autoComplete="name"
            placeholder="Seu nome completo"
            onChange={(event) => setLegalName(event.target.value)}
          />
        </SetupField>
        <SetupField
          label="Nome social (opcional)"
          htmlFor="setup-preferred-name"
        >
          <Input
            id="setup-preferred-name"
            value={preferredName}
            placeholder="Como prefere ser chamado"
            onChange={(event) => setPreferredName(event.target.value)}
          />
        </SetupField>

        <div className="mt-3 flex flex-col gap-3">
          <Button
            type="submit"
            className="h-12 text-base"
            disabled={!complete || saving}
          >
            <Save className="size-4" />
            {saving ? "Salvando…" : "Salvar"}
          </Button>
          {complete && (
            <Button
              type="button"
              variant="outline"
              className="h-12 text-base"
              disabled={saving}
              onClick={() => void save(true)}
            >
              <UserPen className="size-5" />
              Personalizar mais
            </Button>
          )}
        </div>
      </form>
    </main>
  );
}

function SetupField({
  label,
  htmlFor,
  children,
}: {
  label: string;
  htmlFor: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-2">
      <Label htmlFor={htmlFor}>{label}</Label>
      {children}
    </div>
  );
}
