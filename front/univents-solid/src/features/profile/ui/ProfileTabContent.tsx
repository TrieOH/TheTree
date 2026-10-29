import { Match, Show, Switch, createMemo } from "solid-js";
import { useQuery } from "@trieoh/front-core-solid";
import { userBadgesQueryOptions } from "@/features/badges/api";
import { allProfileBadges } from "@/features/badges/model/profile-badges";
import { ProfileBadges } from "@/features/badges/ui/ProfileBadges";
import { myCertificationsQueryOptions } from "@/features/certifications/api";
import type { CertificationI } from "@/features/certifications/model";
import { UserCertificationsSection } from "@/features/certifications/ui/UserCertificationsSection";
import { badgeEditionLocationsQueryOptions } from "@/features/editions/api";
import { PurchasesContent } from "@/features/purchases/ui/PurchasesContent";
import { QueryBoundary } from "@/shared/ui/QueryBoundary";
import { ProfileCard } from "./ProfileCard";
import { ProfileEmptyState } from "./ProfileEmptyState";
import { PROFILE_TAB_LABELS, type ProfileTab } from "../model/profile-tabs";

export function ProfileTabContent(props: {
  tab: Exclude<ProfileTab, "about">;
  actorId: string;
  publicIdentifier: string;
  participantName: string;
  ownProfile: boolean;
}) {
  return (
    <div class="mx-auto mt-5 max-w-7xl px-4">
      <Switch>
        <Match when={props.tab === "badges"}>
          <BadgesTab
            actorId={props.actorId}
            publicIdentifier={props.publicIdentifier}
            participantName={props.participantName}
          />
        </Match>
        <Match when={props.tab === "certificates"}>
          <CertificatesTab
            ownProfile={props.ownProfile}
            participantName={props.participantName}
          />
        </Match>
        <Match when={props.tab === "purchases"}>
          <PurchasesContent />
        </Match>
      </Switch>
    </div>
  );
}

function BadgesTab(props: {
  actorId: string;
  publicIdentifier: string;
  participantName: string;
}) {
  const badgesQuery = useQuery(() => userBadgesQueryOptions(props.actorId));
  // Badges print their edition's location. It arrives from a second request;
  // rendering them before it lands would swap the text in after first paint, so
  // the tab waits for it to settle. A failed lookup just leaves it blank.
  const locationsQuery = useQuery(() => badgeEditionLocationsQueryOptions(props.actorId));
  const locations = createMemo(
    () => new Map(Object.entries(locationsQuery().data ?? {})),
  );

  return (
    <QueryBoundary
      query={badgesQuery}
      what="os crachás"
      skeleton={<CollectionSkeleton tab="badges" />}
    >
      {(badges) => (
        <Show
          when={!locationsQuery().isLoading}
          fallback={<CollectionSkeleton tab="badges" />}
        >
          <ProfileBadges
            badges={allProfileBadges(badges())}
            profileIdentifier={props.publicIdentifier}
            participantName={props.participantName}
            editionLocations={locations()}
          />
        </Show>
      )}
    </QueryBoundary>
  );
}

function CertificatesTab(props: { ownProfile: boolean; participantName: string }) {
  // Certificates are private: only the owner's session can list them.
  const certsQuery = useQuery(() => ({
    ...myCertificationsQueryOptions(),
    enabled: props.ownProfile,
  }));

  return (
    <Switch>
      <Match when={!props.ownProfile}>
        <ProfileEmptyState message="Os certificados deste usuário são visíveis apenas para ele." />
      </Match>
      <Match when={props.ownProfile}>
        <QueryBoundary
          query={certsQuery}
          what="seus certificados"
          skeleton={<CollectionSkeleton tab="certificates" />}
        >
          {(certifications) => (
            <UserCertificationsSection
              certifications={certifications() as CertificationI[]}
              participantName={props.participantName}
            />
          )}
        </QueryBoundary>
      </Match>
    </Switch>
  );
}

function CollectionSkeleton(props: { tab: ProfileTab }) {
  return (
    <ProfileCard title={PROFILE_TAB_LABELS[props.tab]}>
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        <div class="h-44 rounded-md bg-muted animate-pulse" />
        <div class="h-44 rounded-md bg-muted animate-pulse" />
        <div class="h-44 rounded-md bg-muted animate-pulse" />
      </div>
    </ProfileCard>
  );
}
