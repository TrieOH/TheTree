import {
  registerUploadAssociationHandler,
  UploadAssociationError,
  uploadAssociationErrorFromResponse,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import { authQueryFetcher } from "@/shared/lib/api/fetch";
import type { ProgramI } from "../model";
import { patchProgramFn } from "./index";
import { programKeys } from "./query-keys";

registerUploadAssociationHandler("program-image", async (task, url) => {
  const editionId = task.association?.input?.editionId;
  if (typeof editionId !== "string")
    throw new UploadAssociationError("Edição do programa não encontrada.", {
      status: 400,
    });
  const programs = await authQueryFetcher<ProgramI[]>(
    `/editions/${editionId}/programs`,
  );
  const program = programs.find((item) => item.id === task.owner.id);
  if (!program)
    throw new UploadAssociationError("Programa não encontrado.", {
      status: 404,
    });
  const response = await patchProgramFn(program.id, {
    kind: program.kind,
    name: program.name,
    description: program.description,
    min_access_level: program.min_access_level,
    staff_only: program.staff_only,
    price: program.price,
    banner_url: url,
  });
  if (!response.success)
    throw uploadAssociationErrorFromResponse(
      response,
      "Não foi possível associar a imagem.",
    );
  getContext().queryClient.setQueryData<ProgramI[]>(
    programKeys.byEdition(editionId),
    (old = []) =>
      old.map((item) => (item.id === response.data.id ? response.data : item)),
  );
});
