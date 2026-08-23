import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import {
  allJoinedEventsQueryOptions,
  allOwnEventsQueryOptions,
} from "@/features/events/api";
import { myTicketQueryOptions } from "@/features/tickets/api";
import { myParticipationsQueryOptions } from "../api";
import { useOccurrenceRegistrationMutation } from "../api/mutations";

export function useProgramRegistration(
  editionId: string,
  eventId: string,
  eventSlug: string,
  redirectPath: string,
) {
  const navigate = useNavigate();
  const { auth, isAuthenticated } = useAuth();
  const { data: heldTicket, isPending: isTicketPending } = useQuery(
    myTicketQueryOptions(editionId, isAuthenticated),
  );
  const { data: participations, isPending: areParticipationsPending } =
    useQuery(myParticipationsQueryOptions(editionId, isAuthenticated));
  const mutation = useOccurrenceRegistrationMutation(editionId);
  const { data: ownEvents = [] } = useQuery({
    ...allOwnEventsQueryOptions(),
    enabled: isAuthenticated,
  });
  const { data: joinedEvents = [] } = useQuery({
    ...allJoinedEventsQueryOptions(),
    enabled: isAuthenticated,
  });
  const participationByOccurrence = new Map(
    participations?.map((participation) => [
      participation.occurrence_id,
      participation.status,
    ]) ?? [],
  );

  return {
    isAuthenticated,
    hasTicket: isAuthenticated
      ? isTicketPending || areParticipationsPending
        ? undefined
        : Boolean(heldTicket)
      : null,
    ticketStatus: heldTicket?.status,
    accessLevel: heldTicket?.ticket_type.access_level,
    isStaff:
      ownEvents.some((event) => event.id === eventId) ||
      joinedEvents.some((event) => event.id === eventId) ||
      auth.profile()?.id ===
        ownEvents.find((event) => event.id === eventId)?.owner_id,
    participationByOccurrence,
    pendingOccurrenceId: mutation.isPending
      ? mutation.variables?.occurrenceId
      : undefined,
    onToggle: (occurrenceId: string, registered: boolean) => {
      if (!isAuthenticated) {
        void navigate({
          to: "/auth",
          search: { redirect: redirectPath },
        });
      } else if (!heldTicket) {
        void navigate({
          to: "/events/$slug/store",
          params: { slug: eventSlug },
          search: { tab: "tickets" },
        });
      } else {
        mutation.mutate({ occurrenceId, registered });
      }
    },
  };
}
