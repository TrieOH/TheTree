import { orvalData } from "@trieoh/api-client";
import { listEditionPrograms } from "@trieoh/univents-api";
import {
  registerUploadAssociationHandler,
  UploadAssociationError,
} from "@/features/upload-queue";
import { getContext } from "@/integrations/tanstack-query/root-provider";
import type { ProgramI } from "../model";
import { syncProgramCache } from "./cache";
import { patchProgramFn } from "./index";

registerUploadAssociationHandler("program-image", async (task, url) => {
  const editionId = task.association?.input?.editionId;
  if (typeof editionId !== "string")
    throw new UploadAssociationError("Edição do programa não encontrada.", {
      status: 400,
    });
  const programs = await listEditionPrograms(editionId, { public: true }).then(
    orvalData<ProgramI[]>,
  );
  const program = programs.find((item) => item.id === task.owner.id);
  if (!program)
    throw new UploadAssociationError("Programa não encontrado.", {
      status: 404,
    });
  const updated = await patchProgramFn(program.id, {
    kind: program.kind,
    name: program.name,
    description: program.description,
    min_access_level: program.min_access_level,
    staff_only: program.staff_only,
    price: program.price,
    banner_url: url,
  });
  syncProgramCache(getContext().queryClient, updated);
});
