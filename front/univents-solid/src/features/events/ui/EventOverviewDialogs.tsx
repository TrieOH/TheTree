import type { JSX } from "@solidjs/web";
import { AlertModal } from "@/widgets/ui/AlertModal";
import type { EventI } from "../model";
import {
  ManageEventDialog,
  type ManageEventValues,
} from "./ManageEventDialog";

export interface EventOverviewDialogsProps {
  event: EventI | null;
  editOpen: boolean;
  publishOpen: boolean;
  discontinueOpen: boolean;
  disconnectOpen: boolean;
  publishing: boolean;
  discontinuing: boolean;
  disconnecting: boolean;
  onEditOpenChange: (open: boolean) => void;
  onPublishOpenChange: (open: boolean) => void;
  onDiscontinueOpenChange: (open: boolean) => void;
  onDisconnectOpenChange: (open: boolean) => void;
  onEdit: (values: ManageEventValues) => Promise<boolean>;
  onPublish: () => void;
  onDiscontinue: () => void;
  onDisconnect: () => void;
}

export function EventOverviewDialogs(
  props: EventOverviewDialogsProps,
): JSX.Element {
  return (
    <>
      <ManageEventDialog
        open={props.editOpen}
        onOpenChange={props.onEditOpenChange}
        event={props.event}
        onSubmit={props.onEdit}
      />

      <AlertModal
        open={props.disconnectOpen}
        onOpenChange={props.onDisconnectOpenChange}
        title="Desconectar Mercado Pago?"
        description="Este evento deixará de receber novos pagamentos até uma conta ser conectada novamente."
        confirmLabel="Desconectar"
        variant="destructive"
        loading={props.disconnecting}
        onConfirm={props.onDisconnect}
      />

      <AlertModal
        open={props.publishOpen}
        onOpenChange={props.onPublishOpenChange}
        title="Publicar evento?"
        description="Depois de publicar, o painel público ficará disponível para o evento."
        confirmLabel="Publicar evento"
        variant="default"
        loading={props.publishing}
        onConfirm={props.onPublish}
      />

      <AlertModal
        open={props.discontinueOpen}
        onOpenChange={props.onDiscontinueOpenChange}
        title="Descontinuar evento?"
        description="O evento deixará de ser ativo e a data de atualização será atualizada."
        confirmLabel="Descontinuar evento"
        variant="destructive"
        loading={props.discontinuing}
        onConfirm={props.onDiscontinue}
      />
    </>
  );
}
