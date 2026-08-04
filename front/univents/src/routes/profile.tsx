import {
  createFileRoute,
  Outlet,
  useRouterState,
} from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { useCallback } from "react";
import { requireAuth } from "@/features/auths/lib/route-guard";
import { ProfileView } from "@/features/profile/ui/profile-view";
import { cn } from "@/shared/lib/utils";

export const Route = createFileRoute("/profile")({
  beforeLoad: requireAuth,
  component: MyProfilePage,
});

function MyProfilePage() {
  const { auth } = useAuth();
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });
  const actorId = auth.profile()?.id;
  const loadProfile = useCallback(
    (id: string) => auth.getActorProfile(id),
    [auth],
  );
  const isProfileIndex = pathname.replace(/\/$/, "") === "/profile";
  return (
    <>
      <div className={cn(!isProfileIndex && "hidden")}>
        <ProfileView actorId={actorId} loadProfile={loadProfile} ownProfile />
      </div>
      {!isProfileIndex && <Outlet />}
    </>
  );
}
