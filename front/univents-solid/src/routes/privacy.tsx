import { createFileRoute } from "@tanstack/solid-router";
import type { JSX } from "@solidjs/web";
import Shield from "~icons/lucide/shield";
import { LegalPage } from "@/widgets/legal/LegalPage";

const ShieldIcon = Shield as unknown as () => JSX.Element;

const text =
  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.";
export const Route = createFileRoute("/privacy")({ component: PrivacyPage });
function PrivacyPage() {
  return (
    <LegalPage
      icon={() => <ShieldIcon />}
      label="Legal"
      title="Política de Privacidade"
      description="Explicamos quais dados coletamos, como usamos as informações e quais controles você tem sobre sua conta."
      question="Dúvidas sobre o tratamento de dados?"
      sections={[
        "Quais dados podem ser coletados",
        "Como os dados podem ser utilizados",
        "Como os dados podem ser armazenados",
        "Com quem os dados podem ser compartilhados",
        "Direitos do usuário",
        "Forma de solicitar esclarecimentos ou exclusão de dados",
      ].map((title, i) => ({ title: `${i + 1}. ${title}`, content: text }))}
    />
  );
}
