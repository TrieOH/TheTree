import type { JSX } from "@solidjs/web";
import AlertTriangleIcon from "~icons/lucide/alert-triangle";
import AwardIcon from "~icons/lucide/award";
import FileTextIcon from "~icons/lucide/file-text";

import { Tabs, type TabsItem } from "@trieoh/ui-solid";

const AlertTriangle = AlertTriangleIcon as unknown as (props: { class?: string }) => JSX.Element;
const Award = AwardIcon as unknown as (props: { class?: string }) => JSX.Element;
const FileText = FileTextIcon as unknown as (props: { class?: string }) => JSX.Element;

export type CertificationSection = "templates" | "certificates" | "errors";

export interface CertificationSectionTabsProps {
  active: CertificationSection;
  onChange: (section: CertificationSection) => void;
  templatesCount?: number;
  certificatesCount?: number;
  errorsCount?: number;
  class?: string;
}

export function CertificationSectionTabs(props: CertificationSectionTabsProps): JSX.Element {
  const items = (): TabsItem<CertificationSection>[] => [
    {
      id: "templates",
      label: "Templates",
      icon: FileText,
      count: props.templatesCount,
    },
    {
      id: "certificates",
      label: "Certificados",
      icon: Award,
      count: props.certificatesCount,
    },
    {
      id: "errors",
      label: "Erros de emissão",
      icon: AlertTriangle,
      count: props.errorsCount,
    },
  ];

  return (
    <Tabs<CertificationSection>
      items={items()}
      value={props.active}
      onChange={props.onChange}
      ariaLabel="Certificações da edição"
      class={props.class}
    />
  );
}
