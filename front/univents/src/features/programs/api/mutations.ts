import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
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
  patchOccurrenceFn,
  patchProgramFn,
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
    onSuccess: (r) => {
      if (!r.success)
        return toast.error(r.message || "Não foi possível salvar o programa");
      void qc.invalidateQueries({ queryKey: programKeys.byEdition(editionId) });
      toast.success("Programa salvo");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
export function useDeleteProgramMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: deleteProgramFn,
    onSuccess: (r) => {
      if (!r.success)
        return toast.error(r.message || "Não foi possível excluir o programa");
      void qc.invalidateQueries({ queryKey: programKeys.byEdition(editionId) });
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Programa excluído");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
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
    onSuccess: (r) => {
      if (!r.success)
        return toast.error(r.message || "Não foi possível salvar a ocorrência");
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Ocorrência salva");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
export function useDeleteOccurrenceMutation(editionId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: deleteOccurrenceFn,
    onSuccess: (r) => {
      if (!r.success)
        return toast.error(
          r.message || "Não foi possível excluir a ocorrência",
        );
      void qc.invalidateQueries({
        queryKey: programKeys.occurrences(editionId),
      });
      toast.success("Ocorrência excluída");
    },
    onError: () => toast.error("Erro ao conectar com o servidor"),
  });
}
