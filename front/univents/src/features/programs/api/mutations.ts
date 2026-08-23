import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { productKeys } from "@/features/products/api/query-keys";
import { getErrorMessage } from "@/shared/lib/errors";
import type {
  OccurrenceCreateOutput,
  ProgramCreateInput,
  ProgramCreateOutput,
} from "../model";
import {
  createOccurrenceFn,
  createProgramFn,
  deleteOccurrenceFn,
  deleteProgramFn,
  deregisterOccurrenceFn,
  patchOccurrenceFn,
  patchProgramFn,
  registerOccurrenceFn,
} from ".";
import { programKeys } from "./query-keys";

export function useProgramMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      data,
    }: {
      id?: string;
      data: ProgramCreateInput | ProgramCreateOutput;
    }) => (id ? patchProgramFn(id, data) : createProgramFn(editionId, data)),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: programKeys.byEdition(editionId) });
      toast.success("Programa salvo");
    },
    onError: (error) =>
      toast.error(getErrorMessage(error, "Não foi possível salvar o programa")),
  });
}

export function useDeleteProgramMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: deleteProgramFn,
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: programKeys.byEdition(editionId) });
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Programa excluído");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir o programa"),
      ),
  });
}

export function useOccurrenceMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      programId,
      data,
    }: {
      id?: string;
      programId?: string;
      data: OccurrenceCreateOutput;
    }) => {
      if (id) return patchOccurrenceFn(id, data);
      if (!programId)
        throw new Error("programId é obrigatório para criar uma ocorrência");
      return createOccurrenceFn(programId, data);
    },
    onSuccess: () => {
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Ocorrência salva");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível salvar a ocorrência"),
      ),
  });
}

export function useDeleteOccurrenceMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: deleteOccurrenceFn,
    onSuccess: () => {
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Ocorrência excluída");
    },
    onError: (error) =>
      toast.error(
        getErrorMessage(error, "Não foi possível excluir a ocorrência"),
      ),
  });
}

export function useOccurrenceRegistrationMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({
      occurrenceId,
      registered,
    }: {
      occurrenceId: string;
      registered: boolean;
    }) =>
      registered
        ? deregisterOccurrenceFn(occurrenceId)
        : registerOccurrenceFn(occurrenceId),
    onSuccess: (_, { registered }) => {
      void qc.invalidateQueries({
        queryKey: programKeys.myParticipations(editionId),
      });
      void qc.invalidateQueries({
        queryKey: productKeys.storeStock(editionId),
      });
      toast.success(registered ? "Inscrição cancelada" : "Inscrição realizada");
    },
    onError: (error, { registered }) =>
      toast.error(
        getErrorMessage(
          error,
          registered
            ? "Não foi possível cancelar a inscrição"
            : "Não foi possível realizar a inscrição",
        ),
      ),
  });
}
