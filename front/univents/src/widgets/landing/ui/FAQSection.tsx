import { useState } from "react";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/shared/ui/shadcn/accordion";

interface FAQItem {
  question: string;
  answer: string;
}

interface FAQSectionProps {
  items: FAQItem[];
}

export function FAQSection({ items }: FAQSectionProps) {
  const [openValues, setOpenValues] = useState<string[]>([]);

  return (
    <Accordion
      value={openValues}
      onValueChange={(value) => {
        setOpenValues(Array.isArray(value) ? value : [value]);
      }}
      className="space-y-0"
    >
      {items.map((item) => (
        <AccordionItem
          key={item.question}
          value={String(item.question)}
          className="border-b border-border last:border-b-0"
        >
          <AccordionTrigger className="hover:no-underline">
            {item.question}
          </AccordionTrigger>

          <AccordionContent className="text-sm text-muted-foreground">
            {item.answer}
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
