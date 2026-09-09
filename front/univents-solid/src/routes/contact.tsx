import { createFileRoute } from "@tanstack/solid-router";
import type { JSX } from "@solidjs/web";
import { createSignal } from "solid-js";
import { Reveal } from "@/shared/ui/Reveal";
import Headset from "~icons/lucide/headset";
import Mail from "~icons/lucide/mail";
import Phone from "~icons/lucide/phone";

const HeadsetIcon = Headset as unknown as () => JSX.Element;
const MailIcon = Mail as unknown as () => JSX.Element;
const PhoneIcon = Phone as unknown as () => JSX.Element;

const SUPPORT_EMAIL = "suporte@trieoh.com";
const SUPPORT_PHONE = "+55 (99) 99999-9999";
export const Route = createFileRoute("/contact")({ component: ContactPage });
function ContactPage() {
  const [subject, setSubject] = createSignal("");
  const [message, setMessage] = createSignal("");
  const send = () => {
    window.location.href = `mailto:${SUPPORT_EMAIL}?subject=${encodeURIComponent(subject())}&body=${encodeURIComponent(message())}`;
  };
  return (
    <main class="min-h-screen bg-background pb-28 md:pb-40">
      <section class="border-b border-border/40 bg-card/30">
        <div class="mx-auto max-w-5xl px-4 py-6 md:px-6 md:py-8">
          <div class="mb-2 flex items-center gap-3 text-primary">
            <div class="flex size-8 items-center justify-center rounded-lg bg-primary/10">
              <HeadsetIcon />
            </div>
            <span class="text-xs font-semibold uppercase tracking-widest">
              Suporte
            </span>
          </div>
          <h1 class="text-2xl font-semibold tracking-tight md:text-3xl">
            Contato
          </h1>
          <p class="mt-1.5 max-w-2xl text-sm leading-relaxed text-muted-foreground md:text-base">
            Fale com nossa equipe para suporte, parcerias ou dúvidas sobre a
            plataforma.
          </p>
        </div>
      </section>
      <section class="mx-auto max-w-5xl px-4 py-12 md:px-6 md:py-16">
        <div class="grid gap-8 md:grid-cols-2 lg:gap-12">
          <Reveal direction="left">
            <div class="space-y-8">
              <div>
                <h2 class="mb-2 text-2xl font-semibold tracking-tight">
                  Fale Conosco
                </h2>
                <p class="text-justify text-muted-foreground">
                  Tem alguma dúvida ou precisa de suporte? Preencha o formulário
                  e nós entraremos em contato com você o mais rápido possível.
                  Nossa equipe de atendimento está pronta para te ajudar com
                  qualquer situação.
                </p>
              </div>
              <div class="space-y-6">
                <div class="flex items-center space-x-4">
                  <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <MailIcon />
                  </div>
                  <div>
                    <p class="text-sm font-medium leading-none">Email</p>
                    <p class="mt-1 text-sm text-muted-foreground">
                      {SUPPORT_EMAIL}
                    </p>
                  </div>
                </div>
                <div class="flex items-center space-x-4">
                  <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                    <PhoneIcon />
                  </div>
                  <div>
                    <p class="text-sm font-medium leading-none">Telefone</p>
                    <p class="mt-1 text-sm text-muted-foreground">
                      {SUPPORT_PHONE}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </Reveal>
          <Reveal direction="right" delay={0.2}>
            <form
              class="rounded-xl border border-border/50 bg-card/95 shadow-md backdrop-blur-sm"
              onSubmit={(event) => {
                event.preventDefault();
                send();
              }}
            >
              <div class="p-6 pb-0">
                <h2 class="text-lg font-semibold">Envie uma mensagem</h2>
                <p class="text-sm text-muted-foreground">
                  Preencha os campos abaixo para preparar o envio do e-mail.
                </p>
              </div>
              <div class="space-y-4 p-6">
                <label class="block space-y-2 text-sm font-medium">
                  Assunto
                  <input
                    class="w-full rounded-lg border border-input bg-background px-3 py-2 font-normal"
                    value={subject()}
                    onInput={(event) =>
                      setSubject(
                        (event.currentTarget as HTMLInputElement).value,
                      )
                    }
                    placeholder="Como podemos ajudar?"
                  />
                </label>
                <label class="block space-y-2 text-sm font-medium">
                  Mensagem
                  <textarea
                    class="min-h-36 w-full rounded-lg border border-input bg-background px-3 py-2 font-normal"
                    value={message()}
                    onInput={(event) =>
                      setMessage(
                        (event.currentTarget as HTMLTextAreaElement).value,
                      )
                    }
                    placeholder="Descreva sua dúvida ou solicitação..."
                  />
                </label>
                <button
                  class="flex w-full items-center justify-center rounded-full bg-primary px-6 py-3 text-sm font-medium text-primary-foreground"
                  type="submit"
                >
                  <span class="mr-2 inline-flex size-4 shrink-0">
                    <MailIcon />
                  </span>
                  Enviar por E-mail
                </button>
              </div>
            </form>
          </Reveal>
        </div>
      </section>
    </main>
  );
}
