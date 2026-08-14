import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { userBadgesQueryOptions } from "@/features/badges/api";
import type { BadgeProfileBadge } from "@/features/badges/model";
import { BadgePreview } from "@/features/badges/ui/badge-preview";
import { profileKeys } from "@/features/profile/api/query-keys";

export const Route = createFileRoute("/profile/$actorId/badges/$badgeId")({
  component: ProfileBadgePage,
});

function ProfileBadgePage() {
  const { actorId, badgeId } = Route.useParams();
  const { auth } = useAuth();
  const profile = useQuery({
    queryKey: profileKeys.detail(actorId),
    queryFn: () =>
      /^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(actorId)
        ? auth.getActorProfile(actorId)
        : auth.getProfileByHandle(actorId),
  });
  const actorProfileId =
    profile.data?.success ? profile.data.data?.actor_id : undefined;
  const badges = useQuery({
    ...userBadgesQueryOptions(actorProfileId ?? ""),
    enabled: Boolean(actorProfileId),
  });
  const badge = [...(badges.data?.attendant.current ?? []), ...(badges.data?.staff.current ?? [])].find(
    (item) => item.emission_id === badgeId,
  );

  if (profile.isLoading || badges.isLoading) {
    return <main className="grid min-h-dvh place-items-center">Carregando…</main>;
  }
  if (!badge) {
    return <main className="grid min-h-dvh place-items-center">Badge não encontrada.</main>;
  }

  return (
    <main className="flex min-h-dvh flex-col items-center justify-center gap-6 bg-background p-4">
      <Link to="/profile/$actorId" params={{ actorId }} className="self-start inline-flex items-center gap-2 text-sm text-muted-foreground">
        <ArrowLeft className="size-4" /> Voltar ao perfil
      </Link>
      <BadgePreview
        badge={badge as BadgeProfileBadge}
        className="relative w-full max-w-5xl"
        contain
        framed={false}
      />
    </main>
  );
}
