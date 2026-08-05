import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { EmailTemplateKind } from "@trieoh/identityx-api/schemas";
import { useLayoutHeader } from "@trieoh/ui-base";
import {
  Code2,
  KeyRound,
  MailCheck,
  RotateCcw,
  Save,
  Undo2,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { toast } from "sonner";
import {
  emailTemplatesQueryOptions,
  restoreEmailTemplate,
  saveEmailTemplate,
} from "@/features/email-templates/api";
import { cn } from "@/shared/lib/utils";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

export const Route = createFileRoute("/admin/projects/$projectID/emails")({
  component: EmailTemplatesPage,
});

const templateKinds = [
  {
    kind: "verify" as const,
    title: "Verificação de e-mail",
    description: "Enviado após a criação da conta.",
    icon: MailCheck,
  },
  {
    kind: "reset" as const,
    title: "Redefinição de senha",
    description: "Enviado quando o usuário esquece a senha.",
    icon: KeyRound,
  },
];

const variables = [
  ["{{.ActionURL}}", "Link da ação"],
  ["{{.ProjectName}}", "Nome do projeto"],
  ["{{.Email}}", "E-mail do usuário"],
  ["{{.Expiry}}", "Expiração em minutos"],
  ["{{.ProjectDomain}}", "Domínio do projeto"],
] as const;

function EmailTemplatesPage() {
  const { projectID } = Route.useParams();
  const queryClient = useQueryClient();
  const query = emailTemplatesQueryOptions(projectID);
  const { data: templates = [], isLoading } = useQuery(query);
  const [kind, setKind] = useState<EmailTemplateKind>("verify");
  const selected = templates.find((template) => template.kind === kind);
  const [subject, setSubject] = useState("");
  const [body, setBody] = useState("");
  const editorRef = useRef<HTMLTextAreaElement>(null);

  const resetDraft = () => {
    setSubject(selected?.subject ?? "");
    setBody(selected?.body ?? "");
  };

  useEffect(resetDraft, [selected]);

  const dirty =
    subject !== (selected?.subject ?? "") || body !== (selected?.body ?? "");
  const hasActionUrl = body.includes("{{.ActionURL}}");
  const canSave = Boolean(subject.trim() && hasActionUrl && dirty);

  const invalidate = () =>
    queryClient.invalidateQueries({ queryKey: query.queryKey });
  const save = useMutation({
    mutationFn: () => saveEmailTemplate(projectID, kind, { subject, body }),
    onSuccess: async () => {
      await invalidate();
      toast.success("Template salvo");
    },
    onError: (error: Error) => toast.error(error.message),
  });
  const restore = useMutation({
    mutationFn: () => restoreEmailTemplate(projectID, kind),
    onSuccess: async () => {
      await invalidate();
      toast.success("Template padrão restaurado");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  const insertVariable = (variable: string) => {
    const editor = editorRef.current;
    if (!editor) return;
    const start = editor.selectionStart;
    const next = `${body.slice(0, start)}${variable}${body.slice(editor.selectionEnd)}`;
    setBody(next);
    requestAnimationFrame(() => {
      editor.focus();
      editor.setSelectionRange(
        start + variable.length,
        start + variable.length,
      );
    });
  };

  const header = useMemo(
    () => (
      <div>
        <h1 className="text-lg font-semibold">E-mails de autenticação</h1>
        <p className="text-sm text-muted-foreground">
          Personalize as mensagens enviadas aos usuários deste projeto.
        </p>
      </div>
    ),
    [],
  );
  useLayoutHeader(header);

  if (isLoading)
    return <p className="text-sm text-muted-foreground">Carregando…</p>;

  return (
    <div className="space-y-5">
      <div className="grid gap-3 md:grid-cols-2">
        {templateKinds.map((item) => (
          <button
            key={item.kind}
            type="button"
            onClick={() => setKind(item.kind)}
            className={cn(
              "flex items-center gap-3 rounded-lg border p-4 text-left transition-colors",
              kind === item.kind
                ? "border-primary bg-primary/5 ring-1 ring-primary"
                : "border-border bg-card hover:bg-muted/40",
            )}
          >
            <span
              className={cn(
                "flex size-10 shrink-0 items-center justify-center rounded-md",
                kind === item.kind
                  ? "bg-primary text-primary-foreground"
                  : "bg-muted text-muted-foreground",
              )}
            >
              <item.icon className="size-5" />
            </span>
            <span>
              <span className="block text-sm font-semibold">{item.title}</span>
              <span className="text-xs text-muted-foreground">
                {item.description}
              </span>
            </span>
          </button>
        ))}
      </div>

      <div className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_minmax(420px,0.9fr)]">
        <section className="overflow-hidden rounded-lg border border-border bg-card">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
            <div className="flex items-center gap-2">
              <Code2 className="size-4 text-muted-foreground" />
              <h2 className="font-semibold">Editar mensagem</h2>
              <span
                className={cn(
                  "rounded-full px-2 py-0.5 text-[11px] font-medium",
                  selected?.source === "override"
                    ? "bg-primary/10 text-primary"
                    : "bg-muted text-muted-foreground",
                )}
              >
                {selected?.source === "override" ? "Personalizado" : "Padrão"}
              </span>
              {dirty && (
                <span className="text-xs font-medium text-amber-600">
                  Alterações não salvas
                </span>
              )}
            </div>
            <div className="flex gap-2">
              <ShadowButton
                variant="ghost"
                leftIcon={<Undo2 className="size-4" />}
                value="Descartar"
                disabled={!dirty}
                onClick={resetDraft}
              />
              <ShadowButton
                variant="outline"
                leftIcon={<RotateCcw className="size-4" />}
                value="Usar padrão"
                disabled={selected?.source !== "override" || restore.isPending}
                onClick={() => restore.mutate()}
              />
              <ShadowButton
                variant="solid"
                leftIcon={<Save className="size-4" />}
                value={save.isPending ? "Salvando…" : "Salvar"}
                disabled={!canSave || save.isPending}
                onClick={() => save.mutate()}
              />
            </div>
          </div>

          <div className="space-y-5 p-5">
            <label className="block space-y-1.5">
              <span className="flex items-center justify-between text-sm font-medium">
                Assunto
                <span className="font-normal text-muted-foreground">
                  {subject.length}/200
                </span>
              </span>
              <input
                value={subject}
                onChange={(event) => setSubject(event.target.value)}
                maxLength={200}
                placeholder="Ex.: Confirme seu e-mail"
                className="h-11 w-full rounded-md border border-input bg-background px-3 text-sm outline-none focus:ring-2 focus:ring-ring"
              />
            </label>

            <div className="space-y-2">
              <div>
                <p className="text-sm font-medium">Variáveis</p>
                <p className="text-xs text-muted-foreground">
                  Clique para inserir na posição do cursor.
                </p>
              </div>
              <div className="flex flex-wrap gap-2">
                {variables.map(([variable, label]) => (
                  <button
                    key={variable}
                    type="button"
                    title={label}
                    onClick={() => insertVariable(variable)}
                    className="rounded-md border border-border bg-muted/50 px-2.5 py-1.5 font-mono text-xs text-foreground hover:border-primary hover:bg-primary/5"
                  >
                    {variable}
                  </button>
                ))}
              </div>
            </div>

            <label className="block space-y-1.5">
              <span className="flex items-center justify-between text-sm font-medium">
                Corpo do e-mail
                <span className="font-normal text-muted-foreground">HTML</span>
              </span>
              <textarea
                ref={editorRef}
                value={body}
                onChange={(event) => setBody(event.target.value)}
                spellCheck={false}
                className={cn(
                  "min-h-130 w-full resize-y rounded-md border bg-background p-4 font-mono text-xs leading-5 outline-none focus:ring-2 focus:ring-ring",
                  hasActionUrl ? "border-input" : "border-destructive",
                )}
              />
              {!hasActionUrl && (
                <span className="block text-xs text-destructive">
                  Inclua {"{{.ActionURL}}"} para que o usuário consiga concluir
                  a ação.
                </span>
              )}
            </label>
          </div>
        </section>

        <section className="overflow-hidden rounded-lg border border-border bg-card xl:sticky xl:top-4 xl:self-start">
          <div className="border-b border-border px-5 py-4">
            <p className="truncate font-medium">{subject || "Sem assunto"}</p>
            <p className="mt-0.5 text-xs text-muted-foreground">
              Para: pessoa@example.com
            </p>
          </div>
          <iframe
            title="Prévia do e-mail"
            sandbox=""
            srcDoc={previewHtml(body)}
            className="h-180 w-full bg-white"
          />
        </section>
      </div>
    </div>
  );
}

function previewHtml(template: string) {
  const values: Record<string, string> = {
    ActionURL: "https://example.com/auth/action?token=preview",
    ProjectName: "Meu projeto",
    Expiry: "30",
    ProjectDomain: "example.com",
    Email: "pessoa@example.com",
  };
  const rendered = template.replace(
    /\{\{if \.([^}]+)\}\}([\s\S]*?)(?:\{\{else\}\}([\s\S]*?))?\{\{end\}\}/g,
    (_, name: string, truthy: string, fallback = "") =>
      values[name] ? truthy : fallback,
  );
  return Object.entries(values).reduce(
    (html, [name, value]) => html.replaceAll(`{{.${name}}}`, value),
    rendered,
  );
}
