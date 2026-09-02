import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { productKeys } from "@/features/products/api/query-keys";
import { getErrorMessage } from "@/shared/lib/errors";
import type { TicketCreateOutputI } from "../model";
import { syncTicketCaches } from "./cache";
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

export function useCreateTicketMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ editionId, data }: CreateTicketInput) =>
      createTicketFn(data, editionId),
    onSuccess: (ticket) => {
      syncTicketCaches(queryClient, ticket);
      void queryClient.invalidateQueries({
        queryKey: productKeys.storeStock(ticket.edition_id),
      });
      toast.success("Ticket criado com sucesso!");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível criar o ticket")),
  });
}

export function useUpdateTicketMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ ticketId, data }: UpdateTicketInput) =>
      patchTicketFn(data, ticketId),
    onSuccess: (ticket) => {
      syncTicketCaches(queryClient, ticket);
      void queryClient.invalidateQueries({
        queryKey: ticketKeys.myTicket(ticket.edition_id),
      });
      void queryClient.invalidateQueries({
        queryKey: productKeys.storeStock(ticket.edition_id),
      });
      toast.success("Ticket atualizado com sucesso!");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível atualizar o ticket"),
      ),
  });
}
