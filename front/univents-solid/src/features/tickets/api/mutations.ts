import { useQueryClient } from "@trieoh/front-core-solid";
import { withSpan } from "@trieoh/front-core/tracing/browser";
import type {
  CreateTicketTypeRequest,
  PatchTicketTypeRequest,
  TicketType,
} from "@trieoh/univents-api/schemas";
import { Effect } from "effect";
import { apiEffect, useEffectMutation } from "@/shared/lib/effect-query";
import { createTicketFn, patchTicketFn } from "./index";
import { ticketKeys } from "./query-keys";

export interface CreateTicketInput {
  editionId: string;
  data: CreateTicketTypeRequest;
}

export interface PatchTicketInput {
  ticketId: string;
  editionId: string;
  data: PatchTicketTypeRequest;
}

export const createTicketEffect = (
  editionId: string,
  data: CreateTicketTypeRequest,
) =>
  apiEffect(() =>
    withSpan("action:ticket-create", () =>
      createTicketFn(data, editionId),
    ),
  );

export const patchTicketEffect = (
  ticketId: string,
  data: PatchTicketTypeRequest,
) =>
  apiEffect(() =>
    withSpan("action:ticket-patch", () =>
      patchTicketFn(data, ticketId),
    ),
  );

export const useCreateTicketMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ editionId, data }: CreateTicketInput) =>
      createTicketEffect(editionId, data).pipe(
        Effect.tap((_ticket: TicketType) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: ticketKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: ["events", "catalog"],
            });
          }),
        ),
      ),
  });
};

export const usePatchTicketMutation = () => {
  const queryClient = useQueryClient();
  return useEffectMutation({
    retryTransient: true,
    mutationEffect: ({ ticketId, editionId, data }: PatchTicketInput) =>
      patchTicketEffect(ticketId, data).pipe(
        Effect.tap((_ticket: TicketType) =>
          Effect.sync(() => {
            void queryClient.invalidateQueries({
              queryKey: ticketKeys.byEdition(editionId),
            });
            void queryClient.invalidateQueries({
              queryKey: ["events", "catalog"],
            });
          }),
        ),
      ),
  });
};
