import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { ProfileView } from "@/features/profile/ui/profile-view";

export const Route = createFileRoute("/profile/$actorId")({
  component: PublicProfilePage,
});

function PublicProfilePage() {
  const { actorId } = Route.useParams();
  const { auth } = useAuth();
  const viewerActorId = auth.profile()?.id;
  const navigate = useNavigate();
  const loadProfile = useCallback(
    async (identifier: string) => {
      const isActorId = /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(
        identifier,
      );
      const response = isActorId
        ? await auth.getActorProfile(identifier)
        : await auth.getProfileByHandle(identifier);
      if (!response.success) return response;
      const handle = response.data?.handle;

      if (isActorId && handle && handle !== identifier) {
        await navigate({
          to: "/profile/$actorId",
          params: { actorId: handle },
          replace: true,
        });
      }

      return response;
    },
    [auth, navigate],
  );
  return (
    <ProfileView
      actorId={actorId}
      loadProfile={loadProfile}
      viewerActorId={viewerActorId}
    />
  );
}
