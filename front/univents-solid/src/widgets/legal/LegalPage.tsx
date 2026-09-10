import type { JSX } from "@solidjs/web";
import { For } from "solid-js";
import { Reveal } from "@/shared/ui/Reveal";

export interface LegalSection {
  title: string;
  content: string;
}

export function LegalPage(props: {
  icon: () => JSX.Element;
  label: string;
  title: string;
  description: string;
  sections: LegalSection[];
  question: string;
}) {
  return (
    <main class="min-h-screen bg-background pb-28">
      <section class="border-b border-border/40 bg-card/30">
        <div class="mx-auto max-w-4xl px-4 py-6 md:px-6 md:py-8">
          <Reveal direction="left">
            <div class="mb-2 flex items-center gap-3 text-primary">
              <div class="flex size-8 items-center justify-center rounded-lg bg-primary/10">
                <span class="size-4">{props.icon()}</span>
              </div>
              <span class="text-xs font-semibold uppercase tracking-widest">
                {props.label}
              </span>
            </div>
          </Reveal>
          <Reveal delay={0.1}>
            <h1 class="text-2xl font-semibold tracking-tight md:text-3xl">
              {props.title}
            </h1>
            <p class="mt-1.5 max-w-2xl text-sm leading-relaxed text-muted-foreground md:text-base">
              {props.description}
            </p>
          </Reveal>
        </div>
      </section>
      <section class="mx-auto max-w-4xl px-4 py-12 md:px-6 md:py-16">
        <div class="space-y-12">
          <For each={props.sections}>
            {(section) => (
              <Reveal>
                <article class="space-y-4">
                  <h2 class="text-xl font-semibold tracking-tight md:text-2xl">
                    {section.title}
                  </h2>
                  <p class="text-justify text-sm leading-relaxed text-muted-foreground md:text-base">
                    {section.content}
                  </p>
                </article>
              </Reveal>
            )}
          </For>
        </div>
        <div class="mt-16 border-t border-border/60 pt-8">
          <h2 class="mb-4 text-xl font-semibold tracking-tight md:text-2xl">
            {props.question}
          </h2>
          <p class="text-justify text-sm leading-relaxed text-muted-foreground md:text-base">
            Caso você tenha qualquer dúvida,{" "}
            <a
              class="text-primary underline-offset-4 hover:underline"
              href="/contact"
            >
              acesse nossa página de contato
            </a>
            .
          </p>
        </div>
      </section>
    </main>
  );
}
