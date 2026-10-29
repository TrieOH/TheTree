import { orvalItems } from "@trieoh/api-client";
import { listEventMembers } from "@trieoh/univents-api";
import type { QueryClient } from "@trieoh/front-core-solid";

import type { EventMemberRole, EventMemberWithEmailI } from "../model/member";
import type { ActorRecord } from "./actor-emails-handler";
import { eventKeys } from "./query-keys";

export type { EventMemberRole, EventMemberWithEmailI };

export interface AddEventMemberInput {
  eventId: string;
  email: string;
  role: EventMemberRole;
}

export interface RemoveEventMemberInput {
  eventId: string;
  userId: string;
  email: string;
}

export const getEventMembersFn = async (eventId: string): Promise<EventMemberWithEmailI[]> => {
  return listEventMembers(eventId).then(orvalItems<EventMemberWithEmailI>);
};

export const allEventMembersQueryOptions = (eventId: string) => ({
  queryKey: eventKeys.members(eventId),
  queryFn: () => getEventMembersFn(eventId),
});

export const fetchActorEmails = async (actorIds: string[]): Promise<Record<string, ActorRecord>> => {
  const uniqueIds = [...new Set(actorIds.filter(Boolean))];
  if (uniqueIds.length === 0) return {};

  try {
    const res = await fetch("/api/actors/emails", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ actorIds: uniqueIds }),
    });
    if (!res.ok) return {};
    return (await res.json()) as Record<string, ActorRecord>;
  } catch (err) {
    console.warn("[actor-emails] failed to fetch actor emails", err);
    return {};
  }
};

export const actorEmailsQueryOptions = (actorIds: string[]) => ({
  queryKey: ["actor-emails", ...[...new Set(actorIds)].sort()] as const,
  queryFn: () => fetchActorEmails(actorIds),
  enabled: actorIds.length > 0,
  staleTime: 5 * 60 * 1000,
  gcTime: 15 * 60 * 1000,
});

/** Load the list and its identity details before displaying member cards. */
export async function prefetchEventMembers(client: QueryClient, eventId: string) {
  const options = allEventMembersQueryOptions(eventId);
  await client.prefetchQuery(options);
  const members = client.getQueryData<EventMemberWithEmailI[]>(options.queryKey);
  if (members?.length) {
    await client.prefetchQuery(actorEmailsQueryOptions(members.map((member) => member.user_id)));
  }
}
