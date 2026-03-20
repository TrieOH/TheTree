import { useState } from 'react'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger
} from '@/shared/ui/shadcn/collapsible'
import { motion } from 'motion/react'

interface FAQItem {
  question: string
  answer: string
}

interface FAQSectionProps {
  items: FAQItem[]
}

export function FAQSection({ items }: FAQSectionProps) {
  const [openIndex, setOpenIndex] = useState<number | null>(null)

  return (
    <div className="space-y-0">
      {items.map((item, idx) => (
        <Collapsible
          key={idx}
          open={openIndex === idx}
          onOpenChange={(open) => { setOpenIndex(open ? idx : null) }}
        >
          <div className="border-b border-border">
            <CollapsibleTrigger render={
              <button className="w-full py-4 md:py-5 flex justify-between items-center text-left group">
                <span className="text-sm md:text-base font-medium text-foreground group-hover:text-muted-foreground transition-colors pr-4">
                  {item.question}
                </span>
                <motion.span
                  animate={{ rotate: openIndex === idx ? 45 : 0 }}
                  transition={{ duration: 0.2 }}
                  className="text-muted-foreground text-lg md:text-xl shrink-0"
                >
                  +
                </motion.span>
              </button>
            } />

            <CollapsibleContent render={
              <motion.div
                initial={false}
                animate={{
                  height: openIndex === idx ? 'auto' : 0,
                  opacity: openIndex === idx ? 1 : 0
                }}
                transition={{ duration: 0.25, ease: [0.25, 0.1, 0.25, 1] }}
                className="overflow-hidden"
              >
                <div className="pb-4 md:pb-5 text-sm text-muted-foreground leading-relaxed max-w-3xl">
                  {item.answer}
                </div>
              </motion.div>
            } />
          </div>
        </Collapsible>
      ))
      }
    </div >
  )
}