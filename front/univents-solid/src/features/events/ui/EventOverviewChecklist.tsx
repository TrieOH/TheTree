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
      description: props.logoUploading
        ? "Upload em andamento. Aguarde a conclusão."
        : "Identifica o evento nos cards e páginas públicas.",
      completed: props.hasLogo,
      action: props.hasLogo
        ? undefined
        : {
          label: props.logoUploading ? "Enviando..." : "Adicionar",
          disabled: props.logoUploading,
          onClick: props.onAddLogo,
        },
    },
    {
      id: "banner",
      title: "Banner cadastrado",
      description: props.bannerUploading
        ? "Upload em andamento. Aguarde a conclusão."
        : "Imagem principal exibida no topo do evento.",
      completed: props.hasBanner,
      action: props.hasBanner
        ? undefined
        : {
          label: props.bannerUploading ? "Enviando..." : "Adicionar",
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
          action: props.paymentConnected
            ? undefined
            : {
              label: "Conectar",
              onClick: () => {
                document
                  .getElementById("event-payment-panel")
                  ?.scrollIntoView({ behavior: "smooth" });
              },
            },
        },
      ]
      : []),
  ];

  return (
    <StepChecklist
      title="Checklist do evento"
      items={items()}
    />
  );
}
