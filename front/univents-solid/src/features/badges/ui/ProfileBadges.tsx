import type { JSX } from "@solidjs/web";
import { Link } from "@tanstack/solid-router";
import { type Accessor, For, Show, createMemo } from "solid-js";
import type { BadgeProfileBadge } from "../model";
import { badgeDesignSchema } from "../model";
import { BadgePreview } from "./BadgePreview";

export interface ProfileBadgesProps {
  badges: BadgeProfileBadge[];
  profileIdentifier: string;
  participantName: string;
  editionLocations?: Map<string, string>;
}

function ProfileBadgeItem(props: {
  badge: BadgeProfileBadge;
  profileIdentifier: Accessor<string>;
  participantName: Accessor<string>;
  location: Accessor<string>;
}): JSX.Element {
  const design = createMemo(() => {
    const parsed = badgeDesignSchema.safeParse(props.badge.design_data);
    return parsed.success ? parsed.data : null;
  });
  const canvas = createMemo(() => design()?.canvas ?? { width: 321, height: 204 });
  const width = createMemo(() => (160 * canvas().width) / canvas().height);

  return (
    <Link
      to="/profile/$actorId/badges/$badgeId"
      params={{
        actorId: props.profileIdentifier(),
        badgeId: props.badge.emission_id,
      }}
      class="flex w-fit max-w-full items-start"
      aria-label={`Abrir badge ${props.badge.template_name ?? props.badge.edition_name}`}
    >
      <BadgePreview
        badge={props.badge}
        framed={false}
        ticketName={props.badge.ticket_name ?? undefined}
        participantName={props.participantName()}
        location={props.location()}
        class="relative h-auto max-w-full"
        style={{
          width: `${width()}px`,
          height: "auto",
          "max-width": "100%",
        }}
      />
    </Link>
  );
}

export function ProfileBadges(props: ProfileBadgesProps): JSX.Element {
  return (
    <Show when={props.badges.length > 0}>
      <div class="flex flex-wrap items-start justify-center gap-2 sm:justify-start">
        <For each={props.badges}>
          {(badge) => (
            <ProfileBadgeItem
              badge={badge}
              profileIdentifier={() => props.profileIdentifier}
              participantName={() => props.participantName}
              location={() => props.editionLocations?.get(badge.edition_id) ?? ""}
            />
          )}
        </For>
      </div>
    </Show>
  );
}
