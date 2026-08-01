import { createFileRoute } from "@tanstack/react-router";
import { motion } from "motion/react";
import { useState } from "react";
import { Button } from "@/shared/ui/shadcn/button";
import { Input } from "@/shared/ui/shadcn/input";
import { Label } from "@/shared/ui/shadcn/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/shadcn/card";
import { Headset, Mail, Phone } from "lucide-react";

export const Route = createFileRoute("/contact")({
  component: ContactPage,
});

const SUPPORT_EMAIL = "suporte@trieoh.com";
const SUPPORT_PHONE = "+55 (99) 99999-9999";

function ContactPage() {
  const [subject, setSubject] = useState("");
  const [message, setMessage] = useState("");

  const handleEmailClick = () => {
    const mailtoLink = `mailto:${SUPPORT_EMAIL}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(message)}`;
    window.location.href = mailtoLink;
  };

  return (
    <main className="min-h-screen bg-background pb-28 md:pb-40">
      <section className="border-b border-border/40 bg-card/30">
        <div className="mx-auto max-w-5xl px-4 py-6 md:px-6 md:py-8">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
            className="flex items-center gap-3 text-primary mb-2"
          >
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Headset className="h-4 w-4" />
            </div>
            <span className="text-xs font-semibold uppercase tracking-widest">
              Suporte
            </span>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <h1 className="text-2xl md:text-3xl font-semibold tracking-tight text-foreground">
              Contato
            </h1>
            <p className="mt-1.5 max-w-2xl text-sm md:text-base leading-relaxed text-muted-foreground">
              Fale com nossa equipe para suporte, parcerias ou dúvidas sobre a plataforma.
            </p>
          </motion.div>
        </div>
      </section>

      <section className="mx-auto max-w-5xl px-4 py-12 md:px-6 md:py-16">
        <div className="grid gap-8 md:grid-cols-2 lg:gap-12">

          <motion.div
            initial={{ opacity: 0, x: -15 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-10%" }}
            transition={{ duration: 0.4 }}
            className="space-y-8"
          >
            <div>
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Fale Conosco
              </h2>
              <p className="text-muted-foreground text-justify">
                Tem alguma dúvida ou precisa de suporte? Preencha o formulário e nós entraremos em contato com você o mais rápido possível.
                Nossa equipe de atendimento está pronta para te ajudar com qualquer situação.
              </p>
            </div>

            <div className="space-y-6">
              <div className="flex items-center space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                  <Mail className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-sm font-medium leading-none">Email</p>
                  <p className="text-sm text-muted-foreground mt-1">{SUPPORT_EMAIL}</p>
                </div>
              </div>

              <div className="flex items-center space-x-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                  <Phone className="h-5 w-5" />
                </div>
                <div>
                  <p className="text-sm font-medium leading-none">Telefone</p>
                  <p className="text-sm text-muted-foreground mt-1">{SUPPORT_PHONE}</p>
                </div>
              </div>


            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 15 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true, margin: "-10%" }}
            transition={{ duration: 0.4, delay: 0.2 }}
          >
            <Card className="shadow-md bg-card/95 backdrop-blur-sm border-border/50">
              <CardHeader>
                <CardTitle>Envie uma mensagem</CardTitle>
                <CardDescription>
                  Preencha os campos abaixo para preparar o envio do e-mail.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="subject">Assunto</Label>
                  <Input
                    id="subject"
                    placeholder="Como podemos ajudar?"
                    value={subject}
                    onChange={(e) => setSubject(e.target.value)}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="message">Mensagem</Label>
                  <textarea
                    id="message"
                    placeholder="Descreva sua dúvida ou solicitação..."
                    className="flex min-h-[140px] w-full rounded-lg border border-input bg-background dark:bg-input/30 px-3 py-2 text-sm shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50"
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                  />
                </div>
                <Button className="w-full mt-4" onClick={handleEmailClick}>
                  <Mail className="mr-2 h-4 w-4" /> Enviar por E-mail
                </Button>
              </CardContent>
            </Card>
          </motion.div>

        </div>
      </section>
    </main>
  );
}
