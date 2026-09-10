import { createFileRoute } from "@tanstack/solid-router";
import type { JSX } from "@solidjs/web";
import FileText from "~icons/lucide/file-text";
import { LegalPage } from "@/widgets/legal/LegalPage";

const FileTextIcon = FileText as unknown as () => JSX.Element;

const text =
  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.";
export const Route = createFileRoute("/terms")({ component: TermsPage });
function TermsPage() {
  return (
    <LegalPage
      icon={() => <FileTextIcon />}
      label="Legal"
      title="Termos de uso"
      description="Estes termos descrevem as regras de uso da plataforma, responsabilidades do usuário e condições gerais de acesso."
      question="Dúvidas sobre os termos de uso?"
      sections={[
        "Regras de uso da plataforma",
        "Responsabilidades do usuário",
        "Condições gerais de acesso",
        "Condutas permitidas e proibidas",
        "Possíveis alterações nos termos",
      ].map((title, i) => ({ title: `${i + 1}. ${title}`, content: text }))}
    />
  );
}
