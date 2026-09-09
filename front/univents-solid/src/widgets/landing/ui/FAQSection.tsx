import { Accordion } from "@/shared/ui/Accordion";
export interface FAQItem {
  question: string;
  answer: string;
}

export function FAQSection(props: { items: FAQItem[] }) {
  return (
    <Accordion
      items={props.items.map((item) => ({
        value: item.question,
        title: item.question,
        content: item.answer,
      }))}
    />
  );
}
