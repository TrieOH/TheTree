import { createFileRoute } from "@tanstack/react-router";
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
  const actorId = auth.profile()?.id;
  const loadProfile = useCallback(
    (id: string) => auth.getActorProfile(id),
    [auth],
  );

  return <ProfileView actorId={actorId} loadProfile={loadProfile} ownProfile />;
}
