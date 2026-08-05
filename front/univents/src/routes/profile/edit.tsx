import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { profileKeys } from "@/features/profile/api/query-keys";
import type { ProfileSchemaNode } from "@/features/profile/model/profile-data";
import { ProfileEditor } from "@/features/profile/ui/profile-editor";

export const Route = createFileRoute("/profile/edit")({
  beforeLoad: requireAuth,
  component: EditProfilePage,
});

function EditProfilePage() {
  const { auth } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const actorId = auth.profile()?.id;
  const load = useCallback(async () => {
    const [schema, profile] = await Promise.all([
      auth.getProfileSchema(),
      actorId
        ? auth.getActorProfile(actorId)
        : Promise.resolve({
            success: false as const,
            message: "Usuário não autenticado",
          }),
    ]);
    return {
      schema: schema as typeof schema & {
        data?: { schema: ProfileSchemaNode };
      },
      profile,
    };
  }, [actorId, auth]);
  const save = useCallback(
    (profile: ProfileData, handle?: string) =>
      actorId
        ? auth.upsertActorProfile(actorId, { handle, profile })
        : Promise.resolve({
            success: false as const,
            message: "Usuário não autenticado",
          }),
    [actorId, auth],
  );
  return (
    <ProfileEditor
      load={load}
      save={save}
      onCancel={() => navigate({ to: "/profile" })}
      onSaved={() => {
        void queryClient
          .invalidateQueries({ queryKey: profileKeys.all, refetchType: "all" })
          .then(() => navigate({ to: "/profile" }));
      }}
    />
  );
}
