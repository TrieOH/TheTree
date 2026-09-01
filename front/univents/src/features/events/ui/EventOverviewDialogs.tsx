import { AlertModal } from "@/widgets/ui/alert-modal";
import type { EventI } from "../model";
import {
  ManageEventModal,
  type ManageEventModalProps,
} from "./ManageEventModal";

export function EventOverviewDialogs({
  event,
  editOpen,
  publishOpen,
  discontinueOpen,
  disconnectOpen,
  publishing,
  discontinuing,
  disconnecting,
  onEditOpenChange,
  onPublishOpenChange,
  onDiscontinueOpenChange,
  onDisconnectOpenChange,
  onEdit,
  onPublish,
  onDiscontinue,
  onDisconnect,
}: {
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
  onEdit: ManageEventModalProps["onCreate"];
  onPublish: () => void;
  onDiscontinue: () => void;
  onDisconnect: () => void;
}) {
  return (
    <>
      <ManageEventModal
        key={event?.id ?? "event"}
        open={editOpen}
        onOpenChange={onEditOpenChange}
        event={event}
        onCreate={onEdit}
      />
      <AlertModal
        open={disconnectOpen}
        onOpenChange={onDisconnectOpenChange}
        title="Desconectar Mercado Pago?"
        description="Este evento deixará de receber novos pagamentos até uma conta ser conectada novamente."
        confirmLabel="Desconectar"
        variant="destructive"
        loading={disconnecting}
        onConfirm={onDisconnect}
      />
      <AlertModal
        open={publishOpen}
        onOpenChange={onPublishOpenChange}
        title="Publicar evento?"
        description="Depois de publicar, o painel público ficará disponível para o evento."
        confirmLabel="Publicar evento"
        variant="default"
        loading={publishing}
        onConfirm={onPublish}
      />
      <AlertModal
        open={discontinueOpen}
        onOpenChange={onDiscontinueOpenChange}
        title="Descontinuar evento?"
        description="O evento deixará de ser ativo e a data de atualização será atualizada."
        confirmLabel="Descontinuar evento"
        variant="destructive"
        loading={discontinuing}
        onConfirm={onDiscontinue}
      />
    </>
  );
}
