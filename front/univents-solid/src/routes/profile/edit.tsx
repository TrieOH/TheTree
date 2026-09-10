import { createFileRoute, useNavigate } from "@tanstack/solid-router";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileEditor } from "@/features/profile/ui/ProfileEditor";

export const Route = createFileRoute("/profile/edit")({
  beforeLoad: requireAuth,
  component: EditProfilePage,
});

function EditProfilePage() {
  const { auth } = useAuth();
  const navigate = useNavigate();
  const actorId = auth.profile()?.id;
  const finish = () =>
    void navigate({ to: "/profile", search: { tab: "about" } });
  return (
    <ProfileEditor
      load={() =>
        actorId
          ? auth.getActorProfile(actorId).then((profile) => ({ profile }))
          : Promise.resolve({
            profile: { success: false, message: "Usuário não autenticado" },
          })
      }
      save={(profile, handle) =>
        actorId
          ? auth.upsertActorProfile(actorId, {
            handle,
            profile: profile as never,
          })
          : Promise.resolve({
            success: false,
            message: "Usuário não autenticado",
          })
      }
      onCancel={finish}
      onSaved={finish}
    />
  );
}
