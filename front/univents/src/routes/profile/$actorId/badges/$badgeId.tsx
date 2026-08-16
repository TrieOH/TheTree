import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { queryError } from "@trieoh/front-core";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { ArrowLeft } from "lucide-react";
import { userBadgesQueryOptions } from "@/features/badges/api";
import type { BadgeProfileBadge } from "@/features/badges/model";
import { BadgePreview } from "@/features/badges/ui/badge-preview";
import { profileKeys } from "@/features/profile/api/query-keys";
import { asUniventsProfile } from "@/features/profile/model/profile-data";

export const Route = createFileRoute("/profile/$actorId/badges/$badgeId")({
  component: ProfileBadgePage,
});

function ProfileBadgePage() {
  const { actorId, badgeId } = Route.useParams();
  const { auth } = useAuth();
  const profile = useQuery({
    queryKey: profileKeys.detail(actorId),
    queryFn: async () => {
      const response =
        await (/^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(actorId)
          ? auth.getActorProfile(actorId)
          : auth.getProfileByHandle(actorId));
      if (!response.success || !response.data) {
        throw queryError(
          response.message || "Não foi possível carregar este perfil.",
          response.code,
        );
      }
      return response.data;
    },
  });
  const actorProfileId = profile.data?.actor_id;
  const profileData = asUniventsProfile(profile.data?.profile ?? {});
  const badges = useQuery({
    ...userBadgesQueryOptions(actorProfileId ?? ""),
    enabled: Boolean(actorProfileId),
  });
  const badge = [
    ...(badges.data?.attendant.current ?? []),
    ...(badges.data?.staff.current ?? []),
  ].find((item) => item.emission_id === badgeId);

  if (profile.isLoading || badges.isLoading) {
    return (
      <main className="grid min-h-dvh place-items-center bg-background p-6">
        <p className="text-sm text-muted-foreground">Carregando badge…</p>
      </main>
    );
  }
  if (profile.error || badges.error) {
    return (
      <main className="grid min-h-dvh place-items-center bg-background p-6 text-center">
        <div>
          <p className="font-medium">Não foi possível carregar este badge.</p>
          <Link
            to="/profile/$actorId"
            params={{ actorId }}
            className="mt-3 inline-flex text-sm text-primary hover:underline"
          >
            Voltar ao perfil
          </Link>
        </div>
      </main>
    );
  }
  if (!badge) {
    return (
      <main className="grid min-h-dvh place-items-center bg-background p-6 text-center">
        <div>
          <p className="font-medium">Badge não encontrada.</p>
          <Link
            to="/profile/$actorId"
            params={{ actorId }}
            className="mt-3 inline-flex text-sm text-primary hover:underline"
          >
            Voltar ao perfil
          </Link>
        </div>
      </main>
    );
  }

  return (
    <main className="relative flex min-h-dvh items-center justify-center overflow-hidden bg-background">
      <Link
        to="/profile/$actorId"
        params={{ actorId }}
        className="absolute top-4 left-4 z-10 inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3 py-2 text-sm text-muted-foreground shadow-sm backdrop-blur transition-colors hover:text-foreground sm:top-6 sm:left-6"
      >
        <ArrowLeft className="size-4" />
        <span className="hidden sm:inline">Voltar ao perfil</span>
        <span className="sm:hidden">Voltar</span>
      </Link>
      <div className="flex h-dvh w-dvw max-h-full max-w-full items-center justify-center">
        <BadgePreview
          badge={badge as BadgeProfileBadge}
          className="relative h-full w-full max-h-full max-w-full"
          contain
          framed={false}
          participantName={
            profileData.legalName || profileData.preferredName || ""
          }
        />
      </div>
    </main>
  );
}
