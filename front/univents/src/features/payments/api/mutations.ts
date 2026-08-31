import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { syncEventCaches } from "@/features/events/api/cache";
import { eventKeys } from "@/features/events/api/query-keys";
import { getErrorMessage } from "@/shared/lib/errors";
import type { PaymentProviderI } from "../model";
import {
  completeEventSellerFn,
  connectEventSellerFn,
  disconnectEventSellerFn,
} from ".";

export function useConnectEventSellerMutation() {
  return useMutation({
    mutationFn: ({
      eventId,
      provider,
    }: {
      eventId: string;
      provider: PaymentProviderI;
    }) => connectEventSellerFn(eventId, provider),
    onSuccess: ({ auth_url }) => window.location.assign(auth_url),
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao conectar Mercado Pago")),
  });
}

export function useCompleteEventSellerMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      eventId,
      sellerId,
      publicKey,
    }: {
      eventId: string;
      sellerId: string;
      publicKey: string;
    }) =>
      completeEventSellerFn(eventId, {
        seller_id: sellerId,
        public_key: publicKey,
      }),
    onSuccess: (event) => syncEventCaches(queryClient, event),
  });
}

export function useDisconnectEventSellerMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: disconnectEventSellerFn,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: eventKeys.all });
      toast.success("Mercado Pago desconectado");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Erro ao desconectar Mercado Pago")),
  });
}
