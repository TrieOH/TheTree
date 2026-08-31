import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ProfileData } from "@trieoh/identityx-sdk-ts";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { toast } from "sonner";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { syncActorProfileCache } from "@/features/profile/api/cache";
import type { ProfileSchemaNode } from "@/features/profile/model/profile-data";
import { ProfileEditor } from "@/features/profile/ui/profile-editor";

export const Route = createFileRoute("/profile/edit")({
  beforeLoad: requireAuth,
  validateSearch: (search: Record<string, unknown>) => ({
    returnTo: typeof search.returnTo === "string" ? search.returnTo : undefined,
  }),
  component: EditProfilePage,
});

function EditProfilePage() {
  const { auth } = useAuth();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { returnTo } = Route.useSearch();
  const actorId = auth.profile()?.id;
  const finish = () =>
    returnTo
      ? window.location.assign(returnTo)
      : navigate({ to: "/profile", search: { tab: "about" } });
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
    if (!profile.success) {
      toast.info("Configure seu perfil para continuar", {
        id: "profile-not-created",
        description: "Preencha e salve seus dados para acessar o Univents.",
      });
    }
    return {
      schema: schema as typeof schema & {
        data?: { schema: ProfileSchemaNode };
      },
      profile,
    };
  }, [actorId, auth]);
  const save = useCallback(
    async (profile: ProfileData, handle?: string) => {
      const { pfpUrl, ...profileData } = profile;
      if (!actorId) {
        return {
          success: false as const,
          message: "Usuário não autenticado",
        };
      }
      const response = await auth.upsertActorProfile(actorId, {
        handle,
        pfp_url: typeof pfpUrl === "string" ? pfpUrl : null,
        profile: profileData,
      });
      if (response.success) {
        syncActorProfileCache(queryClient, actorId, profile, handle);
      }
      return response;
    },
    [actorId, auth, queryClient],
  );
  return (
    <ProfileEditor load={load} save={save} onCancel={finish} onSaved={finish} />
  );
}
