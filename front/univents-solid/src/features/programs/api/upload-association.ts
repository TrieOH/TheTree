import { orvalData } from "@trieoh/api-client";
import { listEditionPrograms, patchProgram } from "@trieoh/univents-api";
import type { Program } from "@trieoh/univents-api/schemas";
import {
  UploadAssociationError,
  registerUploadAssociationHandler,
} from "@/features/upload-queue";
import { appQueryClient } from "@/shared/lib/query-client";
import { programKeys } from "./query-keys";

async function associateProgramImage(
  task: {
    owner: { id: string };
    association?: { input?: Record<string, unknown> };
  },
  uploadedUrl: string,
) {
  const editionId = task.association?.input?.editionId;
  if (typeof editionId !== "string") {
    throw new UploadAssociationError("Edição da programação não encontrada.", {
      status: 400,
    });
  }

  const programs = await listEditionPrograms(editionId, { public: true })
    .then(orvalData<Program[]>)
    .catch(() => []);

  const program = (programs ?? []).find((item) => item.id === task.owner.id);
  if (!program) {
    throw new UploadAssociationError("Programação não encontrada.", { status: 404 });
  }

  try {
    const response = await patchProgram(program.id, {
      kind: program.kind,
      name: program.name,
      description: program.description ?? undefined,
      min_access_level: program.min_access_level ?? 0,
      staff_only: Boolean(program.staff_only),
      banner_url: uploadedUrl,
      price: program.price ?? undefined,
    });

    if (!response) {
      throw new UploadAssociationError("Não foi possível associar a imagem à programação.", {
        status: 500,
      });
    }

    void appQueryClient.invalidateQueries({
      queryKey: programKeys.byEdition(editionId),
    });
  } catch (err: unknown) {
    if (err instanceof UploadAssociationError) {
      throw err;
    }
    const error = err instanceof Error ? err : undefined;
    const details =
      err && typeof err === "object"
        ? (err as { status?: number; message?: string })
        : undefined;
    throw new UploadAssociationError(
      error?.message ||
      details?.message ||
      "Não foi possível associar a imagem à programação.",
      { status: details?.status || 500 },
    );
  }
}

registerUploadAssociationHandler("program-image", associateProgramImage);
