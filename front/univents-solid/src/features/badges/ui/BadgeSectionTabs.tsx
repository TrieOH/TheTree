import type { JSX } from "@solidjs/web";
import BadgeCheckIcon from "~icons/lucide/badge-check";
import FileTextIcon from "~icons/lucide/file-text";

import { Tabs, type TabsItem } from "@trieoh/ui-solid";

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
  const items = (): TabsItem<BadgeSection>[] => [
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
    <Tabs<BadgeSection>
      items={items()}
      value={props.active}
      onChange={props.onChange}
      ariaLabel="Crachás da edição"
      class={props.class}
    />
  );
}
