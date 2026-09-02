import { StepChecklist } from "@/widgets/ui/step-checklist";

export function EventOverviewChecklist({
  editionCount,
  hasLogo,
  hasBanner,
  hasDescription,
  paymentConnected,
  logoUploading,
  bannerUploading,
  onAddLogo,
  onAddBanner,
  onEdit,
}: {
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
}) {
  const items = [
    {
      id: "edition",
      title: "Edição criada",
      description:
        "Crie uma edição para publicar datas, catálogo e programação.",
      completed: editionCount > 0,
    },
    {
      id: "logo",
      title: "Logo cadastrado",
      description: "Identifica o evento nos cards e páginas públicas.",
      completed: hasLogo,
      action: hasLogo
        ? undefined
        : { label: "Adicionar", disabled: logoUploading, onClick: onAddLogo },
    },
    {
      id: "banner",
      title: "Banner cadastrado",
      description: "Imagem principal exibida no topo do evento.",
      completed: hasBanner,
      action: hasBanner
        ? undefined
        : {
            label: "Adicionar",
            disabled: bannerUploading,
            onClick: onAddBanner,
          },
    },
    {
      id: "description",
      title: "Descrição preenchida",
      description: "Apresente o evento para quem ainda não o conhece.",
      completed: hasDescription,
      action: { label: "Editar", onClick: onEdit },
    },
    ...(editionCount > 0
      ? [
          {
            id: "payment",
            title: "Pagamento conectado",
            description: "Necessário para vender ingressos ou produtos.",
            completed: paymentConnected,
          },
        ]
      : []),
  ];

  return (
    <StepChecklist
      title="Event checklist"
      items={items}
      className="order-8 w-full sm:fixed sm:right-4 sm:top-24 sm:z-40 sm:w-auto!"
      mobileInline
    />
  );
}
