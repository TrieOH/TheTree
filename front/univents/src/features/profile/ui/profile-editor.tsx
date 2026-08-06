import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { ArrowLeft, Save } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { preprocessImageUpload } from "@/features/storage/api";
import { Button } from "@/shared/ui/shadcn/button";
import { Skeleton } from "@/shared/ui/shadcn/skeleton";
import type { ProfileSchemaNode } from "../model/profile-data";
import {
  applyProfileSchemaDefaults,
  withProfileTimestamps,
} from "../model/profile-data";
import { ProfileWarning } from "./profile-editor-fields";
import { InlineProfileEditor } from "./profile-editor-layout";
import {
  normalizeProfileLinks,
  projectProfileValues,
} from "./profile-editor-utils";

export type ProfileImageField = "pfpUrl" | "bannerUrl";
export type PendingImages = Partial<Record<ProfileImageField, File>>;

export interface ProfileEditorProps {
  load: () => Promise<{
    schema: {
      success: boolean;
      data?: { schema: ProfileSchemaNode };
      message?: string;
    };
    profile: {
      success: boolean;
      data?: { handle?: string | null; profile: ProfileData };
      message?: string;
    };
  }>;
  save: (
    profile: ProfileData,
    handle?: string,
  ) => Promise<{
    success: boolean;
    data?: { profile: ProfileData };
    message?: string;
  }>;
  onCancel: () => void;
  onSaved: (profile: ProfileData) => void;
}

export function ProfileEditor({
  load,
  save,
  onCancel,
  onSaved,
}: ProfileEditorProps) {
  const [schema, setSchema] = useState<ProfileSchemaNode>();
  const [original, setOriginal] = useState<ProfileData>({});
  const [values, setValues] = useState<ProfileData>({});
  const [handle, setHandle] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string>();
  const [pendingImages, setPendingImages] = useState<PendingImages>({});

  useEffect(() => {
    let active = true;
    setError(undefined);
    load()
      .then(({ schema: schemaResult, profile }) => {
        if (!active) return;
        if (!schemaResult.success || !schemaResult.data) {
          setError(
            schemaResult.message ||
              "Não foi possível carregar o schema do perfil.",
          );
        } else {
          setSchema(schemaResult.data.schema);
        }
        if (profile.success && profile.data) {
          setHandle(profile.data.handle ?? "");
          setOriginal(profile.data.profile);
          setValues(
            schemaResult.data
              ? applyProfileSchemaDefaults(
                  schemaResult.data.schema,
                  profile.data.profile,
                )
              : profile.data.profile,
          );
        } else if (schemaResult.data) {
          setValues(applyProfileSchemaDefaults(schemaResult.data.schema, {}));
        }
      })
      .catch((cause: unknown) => {
        if (!active) return;
        setError(
          cause instanceof Error
            ? cause.message
            : "Não foi possível carregar o perfil.",
        );
      })
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [load]);

  if (loading)
    return (
      <main className="min-h-dvh bg-background pb-32">
        <Skeleton className="h-40 w-full sm:h-48 md:h-56" />
        <div className="mx-auto max-w-7xl px-4">
          <Skeleton className="-mt-12 size-24 rounded-full border-4 border-background md:-mt-16 md:size-32" />
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            <Skeleton className="h-28 rounded-md" />
            <Skeleton className="h-28 rounded-md" />
          </div>
        </div>
        <div className="mx-auto mt-5 grid max-w-7xl gap-5 px-4 md:grid-cols-[minmax(0,1fr)_360px]">
          <Skeleton className="h-96 rounded-md" />
          <Skeleton className="h-64 rounded-md" />
        </div>
      </main>
    );
  if (error || !schema)
    return (
      <main className="mx-auto max-w-3xl px-4 py-16 md:px-6">
        <ProfileWarning message={error || "Schema de perfil indisponível."} />
        <div className="mt-4 flex justify-center">
          <Button type="button" variant="outline" onClick={onCancel}>
            <ArrowLeft className="size-4" />
            Voltar ao perfil
          </Button>
        </div>
      </main>
    );

  return (
    <main className="min-h-dvh bg-background pb-32">
      <form
        onSubmit={async (event) => {
          event.preventDefault();
          const submittedHandle = handle.trim().replace(/^@/, "");
          if (
            /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(
              submittedHandle,
            )
          ) {
            toast.error("Escolha um handle que não seja um UUID.");
            return;
          }
          setSaving(true);
          try {
            const imageEntries = Object.entries(pendingImages) as [
              ProfileImageField,
              File,
            ][];
            const uploadResults = await Promise.allSettled(
              imageEntries.map(async ([field, file]) => ({
                field,
                url: await preprocessImageUpload(
                  file,
                  "profiles/images",
                  crypto.randomUUID(),
                ),
              })),
            );
            const nextValues = normalizeProfileLinks(values);
            const failedImages: string[] = [];
            uploadResults.forEach((result, index) => {
              const [field] = imageEntries[index] ?? [];
              if (!field) return;
              if (result.status === "fulfilled") {
                nextValues[result.value.field] = result.value.url;
              } else {
                const label = field === "pfpUrl" ? "foto do perfil" : "banner";
                const reason =
                  result.reason instanceof Error
                    ? result.reason.message
                    : "falha no upload ou na moderação";
                failedImages.push(`${label} (${reason})`);
              }
            });
            const response = await save(
              projectProfileValues(
                schema,
                withProfileTimestamps(original, nextValues),
              ),
              handle.trim() || undefined,
            );
            if (!response.success)
              throw new Error(
                response.message || "Não foi possível salvar o perfil.",
              );
            toast.success("Perfil atualizado");
            if (failedImages.length > 0) {
              toast.warning(
                `As outras informações foram salvas, mas houve erro em: ${failedImages.join("; ")}.`,
              );
            }
            onSaved(nextValues);
          } catch (cause) {
            toast.error(
              cause instanceof Error
                ? cause.message
                : "Não foi possível salvar o perfil.",
            );
          } finally {
            setSaving(false);
          }
        }}
      >
        <InlineProfileEditor
          schema={schema}
          values={values}
          handle={handle}
          onHandleChange={setHandle}
          onChange={setValues}
          onCancel={onCancel}
          pendingImages={pendingImages}
          onImageSelect={(field, file) =>
            setPendingImages((current) => ({ ...current, [field]: file }))
          }
        />
        <div className="fixed inset-x-0 bottom-0 z-60 border-t border-border bg-background/95 p-3 shadow-[0_-8px_30px_rgb(0_0_0/0.08)] backdrop-blur">
          <div className="mx-auto flex max-w-7xl items-center justify-between gap-3 px-1 md:px-4">
            <Button
              type="button"
              variant="ghost"
              className="h-8 rounded-sm"
              disabled={saving}
              onClick={onCancel}
            >
              <ArrowLeft className="size-4" />
              Cancelar
            </Button>
            <p className="hidden text-sm text-muted-foreground sm:block">
              Revise a prévia antes de salvar.
            </p>
            <Button type="submit" disabled={saving} className="h-8 rounded-sm">
              <Save className="size-4" />
              {saving ? "Salvando…" : "Salvar alterações"}
            </Button>
          </div>
        </div>
      </form>
    </main>
  );
}
