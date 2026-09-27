import { Show, For, createEffect, createMemo, createSignal } from "solid-js";
import { useQuery } from "@trieoh/front-core-solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import type { BadgeProfileGroups } from "@/features/badges/model";
import { allProfileBadges } from "@/features/badges/model/profile-badges";
import { ProfileBadges } from "@/features/badges/ui/ProfileBadges";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import type { CertificationI } from "@/features/certifications/model";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";
import { allPublicEventsQueryOptions } from "@/features/events/api";
import type { EventI } from "@/features/events/model";
import { getPublicEditionsFn } from "@/features/editions/api";
import { PurchasesContent } from "@/features/purchases/ui/PurchasesContent";
import { ProfileCollectionItem } from "./ProfileCollectionItem";
import { ProfileCard } from "./ProfileCard";
import { ProfileEmptyState } from "./ProfileEmptyState";
import type { ProfileTab } from "./profile-tabs";

export function ProfileTabContent(props: {
  tab: Exclude<ProfileTab, "about">;
  actorId: string;
  publicIdentifier: string;
  participantName: string;
  ownProfile: boolean;
}) {
  const badgesQuery = useQuery(() => ({
    ...userBadgesQueryOptions(props.actorId),
    enabled: props.tab === "badges" && Boolean(props.actorId),
  }));

  const certsQuery = useQuery(() => ({
    ...myCertificationsQueryOptions(),
    enabled: props.tab === "certificates" && props.ownProfile,
  }));

  const data = createMemo(() => {
    if (props.tab === "badges") return badgesQuery().data;
    if (props.tab === "certificates") return certsQuery().data;
    return undefined;
  });

  const isLoading = createMemo(() => {
    if (props.tab === "badges") return badgesQuery().isLoading;
    if (props.tab === "certificates") return certsQuery().isLoading;
    return false;
  });

  const eventsQuery = useQuery(() => ({
    ...allPublicEventsQueryOptions(),
    enabled: props.tab === "badges",
  }));

  const [editionLocations, setEditionLocations] = createSignal<
    Map<string, string>
  >(new Map());

  createEffect(
    () => [props.tab, eventsQuery().data] as const,
    ([tab, evData]) => {
      if (tab !== "badges") return;
      const events = (evData ?? []) as EventI[];
      if (events.length === 0) return;
      void (async () => {
        const map = new Map<string, string>();
        await Promise.all(
          events.map(async (ev) => {
            try {
              const editions = await getPublicEditionsFn(ev.id);
              for (const ed of editions) {
                map.set(ed.id, ed.location_name ?? "");
              }
            } catch {
              // ignore
            }
          }),
        );
        setEditionLocations(new Map(map));
      })();
    },
  );

  return (
    <Show
      when={props.tab === "purchases"}
      fallback={
        <Show
          when={!isLoading()}
          fallback={<CollectionSkeleton tab={props.tab} />}
        >
          <Show
            when={props.tab === "certificates"}
            fallback={
              <Show when={data()}>
                <Show
                  when={props.tab === "badges"}
                  fallback={
                    <div class="mx-auto mt-4 max-w-7xl px-4">
                      <ProfileCard title={tabTitle(props.tab)}>
                        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                          <For each={collectItems(data())}>
                            {(item) => (
                              <ProfileCollectionItem
                                tab={props.tab}
                                item={item}
                              />
                            )}
                          </For>
                        </div>
                      </ProfileCard>
                    </div>
                  }
                >
                  <div class="mx-auto mt-5 max-w-7xl px-4">
                    <ProfileBadges
                      badges={allProfileBadges(data() as BadgeProfileGroups)}
                      profileIdentifier={props.publicIdentifier}
                      participantName={props.participantName}
                      editionLocations={editionLocations()}
                    />
                  </div>
                </Show>
              </Show>
            }
          >
            <div class="mx-auto mt-5 max-w-7xl px-4">
              <Show
                when={props.ownProfile}
                fallback={
                  <ProfileEmptyState message="Os certificados deste usuário são visíveis apenas para ele." />
                }
              >
                <UserCertificationsSection
                  certifications={(certsQuery().data ?? []) as CertificationI[]}
                  participantName={props.participantName}
                />
              </Show>
            </div>
          </Show>
        </Show>
      }
    >
      <div class="mx-auto mt-4 max-w-7xl px-4">
        <PurchasesContent />
      </div>
    </Show>
  );
}

const tabTitle = (tab: Exclude<ProfileTab, "about">) =>
  tab === "badges"
    ? "Crachás"
    : tab === "certificates"
      ? "Certificados"
      : "Compras";

function collectItems(data: unknown): unknown[] {
  if (Array.isArray(data)) return data;
  if (data && typeof data === "object") {
    return Object.values(data as Record<string, unknown>).flatMap(collectItems);
  }
  return [];
}

function CollectionSkeleton(props: { tab: ProfileTab }) {
  return (
    <div class="mx-auto mt-4 max-w-7xl px-4">
      <ProfileCard title={tabTitle(props.tab as Exclude<ProfileTab, "about">)}>
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          <div class="h-44 rounded-md bg-muted animate-pulse" />
          <div class="h-44 rounded-md bg-muted animate-pulse" />
          <div class="h-44 rounded-md bg-muted animate-pulse" />
        </div>
      </ProfileCard>
    </div>
  );
}
