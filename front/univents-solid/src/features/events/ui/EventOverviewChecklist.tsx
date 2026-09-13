import type { JSX } from "@solidjs/web";
import { StepChecklist, type ChecklistItem } from "@/widgets/ui/StepChecklist";

export function EventOverviewChecklist(props: {
  editionCount: number;
  hasLogo: boolean;
  hasBanner: boolean;
  hasDescription: boolean;
  paymentConnected: boolean;
  logoUploading: boolean;
  bannerUploading: boolean;
  onAddLogo: () => void;
  onAddBanner: () => void;
  onEdit: () => void;
}): JSX.Element {
  const items = (): ChecklistItem[] => [
    {
      id: "edition",
      title: "Edição criada",
      description:
        "Crie uma edição para publicar datas, catálogo e programação.",
      completed: props.editionCount > 0,
    },
    {
      id: "logo",
      title: "Logo cadastrado",
      description: "Identifica o evento nos cards e páginas públicas.",
      completed: props.hasLogo,
      action: props.hasLogo
        ? undefined
        : {
            label: "Adicionar",
            disabled: props.logoUploading,
            onClick: props.onAddLogo,
          },
    },
    {
      id: "banner",
      title: "Banner cadastrado",
      description: "Imagem principal exibida no topo do evento.",
      completed: props.hasBanner,
      action: props.hasBanner
        ? undefined
        : {
            label: "Adicionar",
            disabled: props.bannerUploading,
            onClick: props.onAddBanner,
          },
    },
    {
      id: "description",
      title: "Descrição preenchida",
      description: "Apresente o evento para quem ainda não o conhece.",
      completed: props.hasDescription,
      action: { label: "Editar", onClick: props.onEdit },
    },
    ...(props.editionCount > 0
      ? [
          {
            id: "payment",
            title: "Pagamento conectado",
            description: "Necessário para vender ingressos ou produtos.",
            completed: props.paymentConnected,
          },
        ]
      : []),
  ];

  return (
    <StepChecklist
      title="Checklist do evento"
      items={items()}
      class="order-8 w-full sm:fixed sm:right-4 sm:top-24 sm:z-40 sm:w-auto"
      mobileInline
    />
  );
}
