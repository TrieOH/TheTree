import type { JSX } from "@solidjs/web";
import { StepChecklist, type ChecklistItem } from "@/widgets/ui/StepChecklist";

export function EditionOverviewChecklist(props: {
  hasLogo: boolean;
  hasBanner: boolean;
  hasDescription: boolean;
  hasTagline: boolean;
  hasLocation: boolean;
  logoUploading: boolean;
  bannerUploading: boolean;
  onAddLogo: () => void;
  onAddBanner: () => void;
  onEdit: () => void;
}): JSX.Element {
  const items = (): ChecklistItem[] => [
    {
      id: "banner",
      title: "Banner cadastrado",
      description: "Imagem principal exibida no topo da edição.",
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
      id: "logo",
      title: "Logo cadastrado",
      description: "Identifica a edição nos cards e páginas públicas.",
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
      id: "description",
      title: "Descrição preenchida",
      description: "Apresente a edição para quem ainda não a conhece.",
      completed: props.hasDescription,
      action: { label: "Editar", onClick: props.onEdit },
    },
    {
      id: "tagline",
      title: "Tagline definida",
      description: "Resumo curto exibido junto ao nome da edição.",
      completed: props.hasTagline,
      action: { label: "Editar", onClick: props.onEdit },
    },
    {
      id: "location",
      title: "Local definido",
      description: "Local onde a edição será realizada.",
      completed: props.hasLocation,
      action: { label: "Editar", onClick: props.onEdit },
    },
  ];

  return (
    <StepChecklist
      title="Checklist da edição"
      items={items()}
    />
  );
}
