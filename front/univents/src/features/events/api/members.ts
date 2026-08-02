import { queryOptions } from "@tanstack/react-query";
import { createClientOnlyFn } from "@tanstack/react-start";
import { orvalData } from "@trieoh/api-client";
import {
  addEventMember,
  listEventMembers,
  removeEventMember,
} from "@trieoh/univents-api";
import type { EventMember } from "@trieoh/univents-api/schemas";
import type { EventMemberRole } from "../model/member";
import { eventKeys } from "./query-keys";

export type { EventMemberRole } from "../model/member";

export interface EventMemberWithEmailI extends EventMember {
  email?: string;
}

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

export const getEventMembersFn = createClientOnlyFn((eventId: string) => {
  return listEventMembers(eventId).then(orvalData<EventMemberWithEmailI[]>);
});

export const allEventMembersQueryOptions = (eventId: string) =>
  queryOptions({
    queryKey: eventKeys.members(eventId),
    queryFn: () => getEventMembersFn(eventId),
  });

export const addEventMemberFn = createClientOnlyFn(
  ({ eventId, email, role }: AddEventMemberInput) => {
    return addEventMember(eventId, { email, role }).then(
      orvalData<EventMemberWithEmailI>,
    );
  },
);

export const removeEventMemberFn = createClientOnlyFn(
  ({ eventId, userId, email }: RemoveEventMemberInput) => {
    return removeEventMember(eventId, userId, { email }).then(orvalData<null>);
  },
);
