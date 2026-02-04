import { CrudDialog } from "@/shared/ui/dialog/CrudDialog";
import { projectStore } from "../store";


export function ProjectDialog() {
  return (
    <CrudDialog
      store={projectStore}
      title="Project"
    >
      {/* TODO */}
    </CrudDialog>
  );
}