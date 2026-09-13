import { createFileRoute } from "@tanstack/solid-router";
import { useQueryClient } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileEditor } from "@/features/profile/ui/ProfileEditor";
import { cacheProfile, getCachedOwnProfile } from "@/features/profile/api";
import { profileKeys } from "@/features/profile/api/query-keys";

export const Route = createFileRoute("/profile/edit")({
  beforeLoad: requireAuth,
  head: () => ({ meta: [{ title: "Editar perfil - Univents" }] }),
  component: EditProfilePage,
});

function EditProfilePage() {
  const { auth } = useAuth();
  const navigate = Route.useNavigate();
  const queryClient = useQueryClient();
  const actorId = auth.profile()?.id;
  const cachedProfile = getCachedOwnProfile(queryClient, actorId);
  const finish = () => {
    void queryClient.invalidateQueries({ queryKey: profileKeys.details() });
    void navigate({ to: "/profile", search: { tab: "about" } });
  };
  return (
    <ProfileEditor
      initialProfile={cachedProfile}
      load={() =>
        actorId
          ? Promise.all([auth.getProfileSchema(), auth.getActorProfile(actorId)]).then(([, profile]) => {
              cacheProfile(queryClient, profile, actorId);
              return { profile };
            })
          : Promise.resolve({
            profile: { success: false, message: "Usuário não autenticado" },
          })
      }
      save={(profile, handle) => {
        const { pfpUrl, ...data } = profile;
        return actorId
          ? auth.upsertActorProfile(actorId, {
            handle,
            pfp_url: typeof pfpUrl === "string" ? pfpUrl : null,
            profile: data as never,
          })
          : Promise.resolve({
            success: false,
            message: "Usuário não autenticado",
          });
      }}
      onCancel={finish}
      onSaved={finish}
    />
  );
}
