import type { JSX } from "@solidjs/web";

import TicketPlusIcon from "~icons/lucide/ticket-plus";

import { AdminCreateCard } from "@/widgets/ui/AdminCreateCard";

const TicketPlus = TicketPlusIcon as unknown as (props: { class?: string }) => JSX.Element;

export interface AdminCreateTicketCardProps {
  index?: number;
  animate?: boolean;
  onCreate: () => void;
}

export function AdminCreateTicketCard(
  props: AdminCreateTicketCardProps,
): JSX.Element {
  return (
    <AdminCreateCard
      title="Novo ticket"
      description="Configure ingressos pagos ou gratuitos, lotes, limite de vagas e níveis de acesso."
      icon={TicketPlus}
      onClick={props.onCreate}
      index={props.index}
      animate={props.animate}
      layout="stacked"
      minHeight="min-h-[8.5rem]"
      class="p-4"
    />
  );
}
