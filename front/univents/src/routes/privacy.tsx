import { createFileRoute } from "@tanstack/react-router";
import { Shield } from "lucide-react";
import { motion } from "motion/react";

export const Route = createFileRoute("/privacy")({
  component: PrivacyPage,
});

const privacySections = [
  {
    id: "dados-coletados",
    title: "1. Quais dados podem ser coletados",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
  {
    id: "uso-dados",
    title: "2. Como os dados podem ser utilizados",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
  {
    id: "armazenamento",
    title: "3. Como os dados podem ser armazenados",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
  {
    id: "compartilhamento",
    title: "4. Com quem os dados podem ser compartilhados",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
  {
    id: "direitos",
    title: "5. Direitos do usuário",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
  {
    id: "esclarecimentos-exclusao",
    title: "6. Forma de solicitar esclarecimentos ou exclusão de dados",
    content:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
  },
];

function PrivacyPage() {
  return (
    <main className="min-h-screen bg-background pb-28">
      <section className="border-b border-border/40 bg-card/30">
        <div className="mx-auto max-w-4xl px-4 py-6 md:px-6 md:py-8">
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4 }}
            className="flex items-center gap-3 text-primary mb-2"
          >
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Shield className="h-4 w-4" />
            </div>
            <span className="text-xs font-semibold uppercase tracking-widest">
              Legal
            </span>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.1 }}
          >
            <h1 className="text-2xl md:text-3xl font-semibold tracking-tight text-foreground">
              Política de Privacidade
            </h1>
            <p className="mt-1.5 max-w-2xl text-sm md:text-base leading-relaxed text-muted-foreground">
              Explicamos quais dados coletamos, como usamos as informações e
              quais controles você tem sobre sua conta.
            </p>
          </motion.div>
        </div>
      </section>

      <section className="mx-auto max-w-4xl px-4 py-12 md:px-6 md:py-16">
        <div className="space-y-12">
          {privacySections.map((section, index) => (
            <motion.div
              key={section.id}
              initial={{ opacity: 0, y: 15 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-10%" }}
              transition={{ duration: 0.4, delay: index * 0.1 }}
              className="space-y-4"
            >
              <h2 className="text-xl md:text-2xl font-semibold tracking-tight">
                {section.title}
              </h2>
              <p className="text-muted-foreground leading-relaxed text-sm md:text-base text-justify">
                {section.content}
              </p>
            </motion.div>
          ))}
        </div>

        <motion.div
          initial={{ opacity: 0, y: 15 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: "-10%" }}
          transition={{ duration: 0.4, delay: 0.2 }}
          className="mt-16 pt-8 border-t border-border/60"
        >
          <h2 className="text-xl md:text-2xl font-semibold tracking-tight mb-4">
            Dúvidas sobre o tratamento de dados?
          </h2>
          <p className="text-muted-foreground leading-relaxed text-sm md:text-base text-justify">
            Caso você tenha qualquer dúvida sobre as informações acima, ou
            deseje entrar em contato com a equipe responsável,{" "}
            <a
              href="/contact"
              className="text-primary hover:underline underline-offset-4"
            >
              acesse nossa página de contato
            </a>
            .
          </p>
        </motion.div>
      </section>
    </main>
  );
}
