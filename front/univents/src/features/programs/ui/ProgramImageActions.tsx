import { ImagePlus, Trash2 } from "lucide-react";
import { useEffect, useRef } from "react";
import { toast } from "sonner";
import { useUploadQueue } from "@/features/upload-queue";
import { Button } from "@/shared/ui/shadcn/button";
import { patchProgramFn } from "../api";
import type { ProgramI } from "../model";

export function ProgramImageActions({
  program,
  editionId,
}: {
  program: ProgramI;
  editionId: string;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const { enqueue, tasks } = useUploadQueue();
  const active = tasks.some(
    (task) =>
      task.owner.type === "program" &&
      task.owner.id === program.id &&
      task.association?.handlerKey === "program-image" &&
      !["completed", "failed", "rejected"].includes(task.status),
  );
  useEffect(() => {
    const id = `program-image-upload-${program.id}`;
    if (active)
      toast.loading("Enviando banner do programa…", { id, duration: Infinity });
    else toast.dismiss(id);
    return () => {
      toast.dismiss(id);
    };
  }, [active, program.id]);
  const remove = async () => {
    await patchProgramFn(program.id, {
      kind: program.kind,
      name: program.name,
      description: program.description,
      min_access_level: program.min_access_level,
      staff_only: program.staff_only,
      price: program.price,
      banner_url: null,
    });
  };
  return (
    <div
      className="flex items-center gap-1"
      onClick={(event) => event.stopPropagation()}
    >
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={(event) => {
          const file = event.target.files?.[0];
          event.target.value = "";
          if (file)
            void enqueue({
              file,
              owner: { type: "program", id: program.id, label: program.name },
              mediaType: "banner",
              storagePath: `editions/${editionId}/programs/${program.id}/banner`,
              association: {
                handlerKey: "program-image",
                input: { editionId },
              },
            });
        }}
      />
      <Button
        type="button"
        size="sm"
        variant="ghost"
        title="Adicionar ou trocar banner"
        disabled={active}
        onClick={() => inputRef.current?.click()}
      >
        <ImagePlus className="size-4" />
      </Button>
      {program.banner_url ? (
        <Button
          type="button"
          size="sm"
          variant="ghost"
          title="Remover banner"
          disabled={active}
          onClick={() => void remove()}
        >
          <Trash2 className="size-4" />
        </Button>
      ) : null}
    </div>
  );
}
