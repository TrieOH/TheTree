import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { ArrowLeft, Globe } from "lucide-react";
import { cn } from "@/shared/lib/utils";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import type { ProfileSchemaNode } from "../model/profile-data";
import { asUniventsProfile } from "../model/profile-data";
import type { PendingImages, ProfileImageField } from "./profile-editor";
import { EditorCard, EditorField, EditorObject } from "./profile-editor-fields";
import { ProfileImageInput } from "./profile-image-input";

export function InlineProfileEditor({
  schema,
  values,
  handle,
  onHandleChange,
  onChange,
  onCancel,
  pendingImages,
  onImageSelect,
}: {
  schema: ProfileSchemaNode;
  values: ProfileData;
  handle: string;
  onHandleChange: (handle: string) => void;
  onChange: (values: ProfileData) => void;
  onCancel: () => void;
  pendingImages: PendingImages;
  onImageSelect: (field: ProfileImageField, file: File) => void;
}) {
  const profile = asUniventsProfile(values);
  return (
    <>
      <section className="relative border-b border-border bg-card shadow-md">
        <ProfileImageInput
          label="Banner"
          currentUrl={profile.bannerUrl}
          file={pendingImages.bannerUrl}
          onSelect={(file) => onImageSelect("bannerUrl", file)}
          variant="banner"
          className={cn(
            "relative h-40 bg-background bg-linear-to-br sm:h-48 md:h-56",
            "from-primary/40 via-primary/15 to-muted",
          )}
        />
        <div className="pointer-events-none absolute inset-x-4 top-4 z-20 mx-auto max-w-7xl">
          <div className="pointer-events-auto w-fit">
            <Button
              type="button"
              variant="secondary"
              size="icon"
              className="rounded-full shadow-lg"
              onClick={onCancel}
              aria-label="Voltar ao perfil"
            >
              <ArrowLeft className="size-4" />
            </Button>
          </div>
        </div>
        <div className="mx-auto max-w-7xl px-4">
          <div className="relative pb-5 md:pb-6">
            <ProfileImageInput
              label="Foto do perfil"
              currentUrl={profile.pfpUrl}
              file={pendingImages.pfpUrl}
              onSelect={(file) => onImageSelect("pfpUrl", file)}
              variant="avatar"
              className="-mt-12 md:-mt-16"
            />
            <div className="mt-3">
              <div className="grid items-end gap-3 sm:grid-cols-2 lg:grid-cols-3">
                <div className="min-w-0 space-y-1.5">
                  <Label htmlFor="profile-handle">Nome de usuário</Label>
                  <Input
                    id="profile-handle"
                    value={handle}
                    placeholder="seu nickname"
                    pattern={"[^\\s/]+"}
                    title="Use um nickname sem espaços ou barras"
                    autoCapitalize="none"
                    autoComplete="username"
                    spellCheck={false}
                    className="h-10 bg-background"
                    onChange={(event) => onHandleChange(event.target.value)}
                  />
                </div>
                {[
                  "preferredName",
                  "legalName",
                  "role",
                  "organization",
                  "pronouns",
                  "tagline",
                ].map((name) => (
                  <EditorField
                    key={name}
                    schema={schema}
                    values={values}
                    onChange={onChange}
                    path={[name]}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      <div className="mx-auto mt-4 grid max-w-7xl gap-4 px-4 md:mt-5 md:grid-cols-[minmax(0,1fr)_360px] md:gap-5">
        <div className="space-y-5">
          <EditorCard title="Sobre mim">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={["aboutMe"]}
              compact
              hideLabel
            />
          </EditorCard>
          <EditorCard title="Idiomas">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={[
                schema.properties?.languages ? "languages" : "specializations",
              ]}
              compact
              hideLabel
            />
          </EditorCard>
          <EditorCard title="Contato" icon={<Globe className="size-4" />}>
            <div className="grid gap-4 sm:grid-cols-2">
              <EditorField
                schema={schema}
                values={values}
                onChange={onChange}
                path={["website"]}
              />
              <EditorField
                schema={schema}
                values={values}
                onChange={onChange}
                path={["contactEmail"]}
              />
            </div>
            <EditorObject
              schema={schema}
              values={values}
              onChange={onChange}
              name="socials"
            />
          </EditorCard>
        </div>

        <aside className="space-y-5">
          <EditorCard title="Privacidade">
            <EditorObject
              schema={schema}
              values={values}
              onChange={onChange}
              name="visibility"
            />
          </EditorCard>
          <EditorCard title="Outros detalhes">
            <EditorField
              schema={schema}
              values={values}
              onChange={onChange}
              path={["timezone"]}
              compact
            />
          </EditorCard>
        </aside>
      </div>
    </>
  );
}
