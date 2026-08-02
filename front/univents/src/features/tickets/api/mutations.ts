import {
  type QueryClient,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { toast } from "sonner";
import type { TicketCreateOutputI, TicketI } from "../model";
import { createTicketFn, patchTicketFn } from "./index";
import { ticketKeys } from "./query-keys";

type CreateTicketInput = {
  editionId: string;
  data: TicketCreateOutputI;
};

type UpdateTicketInput = {
  ticketId: string;
  data: TicketCreateOutputI;
};

function upsertById(tickets: TicketI[] | undefined, ticket: TicketI) {
  const list = tickets ?? [];
  const index = list.findIndex((item) => item.id === ticket.id);

  if (index === -1) return [...list, ticket];

  const next = [...list];
  next[index] = ticket;
  return next;
}

function syncTicketCaches(
  queryClient: QueryClient,
  editionId: string,
  ticket: TicketI,
) {
  queryClient.setQueryData<TicketI[]>(
    ticketKeys.listByEdition(editionId),
    (old) => upsertById(old, ticket),
  );
}

export function useCreateTicketMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateTicketInput) =>
      createTicketFn(data, editionId),
    onSuccess: (res, variables) => {
      syncTicketCaches(queryClient, variables.editionId, res);
      toast.success("Ticket criado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}

export function useUpdateTicketMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ ticketId, data }: UpdateTicketInput) =>
      patchTicketFn(data, ticketId),
    onSuccess: (ticket) => {
      const editionId = ticket.edition_id;
      syncTicketCaches(queryClient, editionId, ticket);
      toast.success("Ticket atualizado com sucesso!");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
