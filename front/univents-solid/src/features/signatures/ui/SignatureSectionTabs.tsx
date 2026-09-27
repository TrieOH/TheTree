import type { JSX } from "@solidjs/web";
import MailIcon from "~icons/lucide/mail";
import PenLineIcon from "~icons/lucide/pen-line";

import { SectionTabs, type SectionTabItem } from "@/shared/ui/SectionTabs";

const PenLine = PenLineIcon as unknown as (props: { class?: string }) => JSX.Element;
const Mail = MailIcon as unknown as (props: { class?: string }) => JSX.Element;

export type SignatureSection = "signatures" | "invites";

export interface SignatureSectionTabsProps {
  active: SignatureSection;
  onChange: (section: SignatureSection) => void;
  signaturesCount?: number;
  invitesCount?: number;
  class?: string;
}

export function SignatureSectionTabs(props: SignatureSectionTabsProps): JSX.Element {
  const items = (): SectionTabItem<SignatureSection>[] => [
    {
      id: "signatures",
      label: "Assinaturas",
      icon: PenLine,
      count: props.signaturesCount,
    },
    {
      id: "invites",
      label: "Convites de assinatura",
      icon: Mail,
      count: props.invitesCount,
    },
  ];

  return (
    <SectionTabs<SignatureSection>
      items={items()}
      active={props.active}
      onChange={props.onChange}
      ariaLabel="Assinaturas da edição"
      class={props.class}
    />
  );
}
