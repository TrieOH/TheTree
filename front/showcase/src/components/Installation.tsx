import { CodeBlock } from "@trieoh/ui-solid";

export function Installation() {
  return (
    <section class="mt-10 min-w-0">
      <h2 class="mb-4 text-lg font-semibold tracking-[-0.4px]">Instalação</h2>
      <p class="my-3.5 text-xs leading-[1.7] text-muted-foreground">
        Adicione ao aplicativo no monorepo.
      </p>
      <CodeBlock
        filename="Terminal"
        format={false}
        code={'pnpm add "@trieoh/ui-solid@workspace:*"'}
      />
      <p class="my-3.5 text-xs leading-[1.7] text-muted-foreground">
        Importe a integração com o Tailwind na folha de estilos principal.
      </p>
      <CodeBlock
        filename="src/App.css"
        format={false}
        code={
          '@import "tailwindcss";\n@import "@trieoh/ui-solid/tailwind.css";'
        }
      />
      <p class="my-3.5 text-xs leading-[1.7] text-muted-foreground">
        Para usar o tema padrão, importe também{" "}
        <code class="font-mono">@trieoh/ui-solid/theme.css</code>.
      </p>
    </section>
  );
}
