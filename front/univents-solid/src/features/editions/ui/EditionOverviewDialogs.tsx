import type { JSX } from "@solidjs/web";
import { AlertModal } from "@/widgets/ui/AlertModal";
import type { EditionI } from "../model";
import {
  ManageEditionDialog,
  type ManageEditionValues,
} from "./ManageEditionDialog";

export interface EditionOverviewDialogsProps {
  edition: EditionI | null;
  editOpen: boolean;
  publishOpen: boolean;
  publishing: boolean;
  onEditOpenChange: (open: boolean) => void;
  onPublishOpenChange: (open: boolean) => void;
  onEdit: (values: ManageEditionValues) => Promise<boolean>;
  onPublish: () => void;
}

export function EditionOverviewDialogs(
  props: EditionOverviewDialogsProps,
): JSX.Element {
  return (
    <>
      <ManageEditionDialog
        open={props.editOpen}
        onOpenChange={props.onEditOpenChange}
        edition={props.edition}
        onSubmit={props.onEdit}
      />

      <AlertModal
        open={props.publishOpen}
        onOpenChange={props.onPublishOpenChange}
        title="Publicar edição?"
        description="Depois de publicar, a edição ficará visível ao público."
        confirmLabel="Publicar edição"
        variant="default"
        loading={props.publishing}
        onConfirm={props.onPublish}
      />
    </>
  );
}
