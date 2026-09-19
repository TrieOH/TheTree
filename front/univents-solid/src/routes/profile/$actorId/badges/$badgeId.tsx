import type { JSX } from "@solidjs/web";
import { Link, createFileRoute } from "@tanstack/solid-router";
import { useQuery, useQueryClient } from "@trieoh/front-core-solid";
import { useAuth } from "@trieoh/identityx-sdk-ts-solid";
import { Show, createMemo } from "solid-js";
import ArrowLeftIcon from "~icons/lucide/arrow-left";
import Loader2Icon from "~icons/lucide/loader-2";

import { badgeKeys, userBadgesQueryOptions } from "@/features/badges/api";
import { allProfileBadges, badgeDesignSchema } from "@/features/badges/model";
import type { BadgeProfileGroups } from "@/features/badges/model";
import { BadgePreview } from "@/features/badges/ui/BadgePreview";
import { editionLocationQueryOptions } from "@/features/editions/api";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { profileDetailQueryOptions, profileKeys } from "@/features/profile/api";
import { asUniventsProfile } from "@/features/profile/model/profile-data";

const ArrowLeft = ArrowLeftIcon as unknown as (props: { class?: string }) => JSX.Element;
const Loader2 = Loader2Icon as unknown as (props: { class?: string }) => JSX.Element;

export const Route = createFileRoute("/profile/$actorId/badges/$badgeId")({
  component: ProfileBadgePage,
});

function ProfileBadgePage(): JSX.Element {
  const params = Route.useParams();
  const { auth } = useAuth();
  const queryClient = useQueryClient();

  const profileQuery = useQuery(() =>
    profileDetailQueryOptions(params().actorId, auth),
  );

  const profileRecord = createMemo(() => {
    const raw = profileQuery().data as any;
    if (!raw) return null;
    return raw.actor_id ? raw : (raw.data ?? null);
  });

  const actorProfileId = createMemo(() => {
    const fromData = profileRecord()?.actor_id;
    if (fromData) return fromData;
    const actorParam = params().actorId;
    if (/^[0-9a-f]{8}(?:-[0-9a-f]{4}){3}-[0-9a-f]{12}$/i.test(actorParam)) {
      return actorParam;
    }
    return undefined;
  });

  const profileData = createMemo(() =>
    asUniventsProfile(profileRecord()?.profile ?? {}),
  );

  const badgesQuery = useQuery(() => ({
    ...userBadgesQueryOptions(actorProfileId() ?? ""),
    enabled: Boolean(actorProfileId()),
  }));

  const badgeGroups = createMemo(() => {
    if (badgesQuery().data) return badgesQuery().data;
    const actorId = actorProfileId();
    if (!actorId) return undefined;
    const cachedDirect = queryClient.getQueryData<BadgeProfileGroups>(
      badgeKeys.user(actorId),
    );
    if (cachedDirect) return cachedDirect;
    const cachedOwn = queryClient.getQueryData<BadgeProfileGroups>(
      profileKeys.tab("badges", actorId, true),
    );
    if (cachedOwn) return cachedOwn;
    const cachedPublic = queryClient.getQueryData<BadgeProfileGroups>(
      profileKeys.tab("badges", actorId, false),
    );
    if (cachedPublic) return cachedPublic;
    return undefined;
  });

  const badge = createMemo(() => {
    const groups = badgeGroups();
    if (!groups) return undefined;
    return allProfileBadges(groups).find(
      (item) => item.emission_id === params().badgeId,
    );
  });

  const eventsQuery = useQuery(() => allPublicEventsQueryOptions());
  const events = createMemo(() => (eventsQuery().data ?? []) as EventI[]);

  const editionId = createMemo(() => badge()?.edition_id);

  const editionQuery = useQuery(() =>
    editionLocationQueryOptions(editionId(), events()),
  );

  const isLoading = createMemo(() => {
    if (badge()) return false;
    if (profileQuery().isLoading) return true;
    if (actorProfileId() && badgesQuery().isLoading) return true;
    if (!profileRecord() && profileQuery().fetchStatus === "fetching") return true;
    if (actorProfileId() && !badgeGroups() && badgesQuery().fetchStatus === "fetching") return true;
    return false;
  });

  const isError = createMemo(() => {
    if (isLoading()) return false;
    if (badge()) return false;
    if (profileQuery().isError) return true;
    if (actorProfileId() && badgesQuery().isError) return true;
    return false;
  });

  const design = createMemo(() => {
    const b = badge();
    if (!b) return null;
    const parsed = badgeDesignSchema.safeParse(b.design_data);
    return parsed.success ? parsed.data : null;
  });

  const canvas = createMemo(() => design()?.canvas ?? { width: 321, height: 204 });
  const aspectRatio = createMemo(() => canvas().width / canvas().height);

  return (
    <Show
      when={!isLoading()}
      fallback={
        <main class="grid min-h-dvh place-items-center bg-background p-6">
          <div class="flex flex-col items-center gap-2 text-muted-foreground">
            <Loader2 class="size-6 animate-spin text-primary" />
            <p class="text-sm">Carregando crachá…</p>
          </div>
        </main>
      }
    >
      <Show
        when={!isError()}
        fallback={
          <main class="grid min-h-dvh place-items-center bg-background p-6 text-center">
            <div>
              <p class="font-medium text-foreground">
                Não foi possível carregar este crachá.
              </p>
              <Link
                to="/profile/$actorId"
                params={{ actorId: params().actorId }}
                search={{ tab: "badges" }}
                class="mt-3 inline-flex items-center gap-1.5 text-sm text-primary hover:underline"
              >
                <ArrowLeft class="size-4" />
                Voltar ao perfil
              </Link>
            </div>
          </main>
        }
      >
        <Show
          when={badge()}
          fallback={
            <main class="grid min-h-dvh place-items-center bg-background p-6 text-center">
              <div>
                <p class="font-medium text-foreground">Crachá não encontrado.</p>
                <Link
                  to="/profile/$actorId"
                  params={{ actorId: params().actorId }}
                  search={{ tab: "badges" }}
                  class="mt-3 inline-flex items-center gap-1.5 text-sm text-primary hover:underline"
                >
                  <ArrowLeft class="size-4" />
                  Voltar ao perfil
                </Link>
              </div>
            </main>
          }
        >
          {(currentBadge) => (
            <main class="relative flex min-h-dvh items-center justify-center overflow-hidden bg-background">
              <Link
                to="/profile/$actorId"
                params={{ actorId: params().actorId }}
                search={{ tab: "badges" }}
                class="absolute top-4 left-4 z-10 inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3 py-2 text-sm text-muted-foreground shadow-xs backdrop-blur-sm transition-colors hover:text-foreground sm:top-6 sm:left-6"
              >
                <ArrowLeft class="size-4" />
                <span class="hidden sm:inline">Voltar ao perfil</span>
                <span class="sm:hidden">Voltar</span>
              </Link>

              <div class="flex h-dvh w-dvw max-h-full max-w-full items-center justify-center">
                <BadgePreview
                  badge={currentBadge()}
                  class="relative"
                  framed={false}
                  style={{
                    width: `min(100dvw, calc(100dvh * ${aspectRatio()}))`,
                    height: `min(100dvh, calc(100dvw / ${aspectRatio()}))`,
                  }}
                  ticketName={(currentBadge() as any).ticket_name ?? undefined}
                  participantName={
                    profileData().preferredName || profileData().legalName || ""
                  }
                  location={editionQuery().data?.location_name ?? ""}
                />
              </div>
            </main>
          )}
        </Show>
      </Show>
    </Show>
  );
}
