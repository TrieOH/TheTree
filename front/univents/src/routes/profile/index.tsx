import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileView } from "@/features/profile/ui/profile-view";

export const Route = createFileRoute("/profile/")({
  beforeLoad: requireAuth,
  component: MyProfilePage,
});

function MyProfilePage() {
  const { auth } = useAuth();
  const navigate = useNavigate();
  const actorId = auth.profile()?.id;
  const loadProfile = useCallback(
    async (id: string) => {
      const response = await auth.getActorProfile(id);
      if (response.success && response.data) {
        await navigate({
          to: "/profile/$actorId",
          params: { actorId: response.data.handle || id },
          replace: true,
        });
      }
      return response;
    },
    [auth, navigate],
  );

  return <ProfileView actorId={actorId} loadProfile={loadProfile} ownProfile />;
}
