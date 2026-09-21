import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { untrack } from "solid-js";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileEditor } from "@/features/profile/ui/ProfileEditor";
import {
  cacheProfile,
  getCachedOwnProfile,
  profileKeys,
  syncActorProfileCache,
} from "@/features/profile/api";

export const Route = createFileRoute("/profile/edit")({
  beforeLoad: requireAuth,
  head: () => ({ meta: [{ title: "Editar perfil - Univents" }] }),
  component: EditProfilePage,
});

function EditProfilePage() {
  const { auth } = useAuth();
  const navigate = Route.useNavigate();
  const queryClient = useQueryClient();
  const actorId = () => auth.profile()?.id;
  const finish = () => {
    void queryClient.invalidateQueries({ queryKey: profileKeys.details() });
    void navigate({ to: "/profile", search: { tab: "about" } });
  };

  const handleSave = (profile: Record<string, unknown>, handle?: string) => {
    const id = actorId();
    const { pfpUrl, ...data } = profile;
    if (!id) {
      return Promise.resolve({
        success: false,
        message: "Usuário não autenticado",
      });
    }
    return auth
      .upsertActorProfile(id, {
        handle,
        pfp_url: typeof pfpUrl === "string" ? pfpUrl : null,
        profile: data as never,
      })
      .then((response) => {
        if (response.success) {
          syncActorProfileCache(queryClient, id, profile, handle);
        }
        return response;
      });
  };

  return (
    <ProfileEditor
      initialProfile={untrack(() => getCachedOwnProfile(queryClient, auth.profile()?.id))}
      load={() => {
        const id = actorId();
        return id
          ? Promise.all([auth.getProfileSchema(), auth.getActorProfile(id)]).then(([, profile]) => {
            cacheProfile(queryClient, profile, id);
            return { profile };
          })
          : Promise.resolve({
            profile: { success: false, message: "Usuário não autenticado" },
          });
      }}
      save={handleSave}
      onCancel={finish}
      onSaved={finish}
    />
  );
}
