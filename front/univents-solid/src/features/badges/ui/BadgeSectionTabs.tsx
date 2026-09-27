import type { JSX } from "@solidjs/web";
import BadgeCheckIcon from "~icons/lucide/badge-check";
import FileTextIcon from "~icons/lucide/file-text";

import { SectionTabs, type SectionTabItem } from "@/shared/ui/SectionTabs";

const BadgeCheck = BadgeCheckIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileText = FileTextIcon as unknown as (props: { class?: string }) => JSX.Element;

export type BadgeSection = "templates" | "emissions";

export interface BadgeSectionTabsProps {
  active: BadgeSection;
  onChange: (section: BadgeSection) => void;
  templatesCount?: number;
  emissionsCount?: number;
  class?: string;
}

export function BadgeSectionTabs(props: BadgeSectionTabsProps): JSX.Element {
  const items = (): SectionTabItem<BadgeSection>[] => [
    {
      id: "templates",
      label: "Templates",
      icon: FileText,
      count: props.templatesCount,
    },
    {
      id: "emissions",
      label: "Crachás emitidos",
      icon: BadgeCheck,
      count: props.emissionsCount,
    },
  ];

  return (
    <SectionTabs<BadgeSection>
      items={items()}
      active={props.active}
      onChange={props.onChange}
      ariaLabel="Crachás da edição"
      class={props.class}
    />
  );
}
