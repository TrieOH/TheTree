import { schemaQueryOptions, writeSchema } from "#/features/schema/api";
import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect } from "react";
import { toast } from "sonner";
import SchemaEditor from "#/features/schema/ui/schema-editor";
import { Database, ShieldCheck, Users, RefreshCw } from "lucide-react";
import { useSuspenseQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTheme } from "next-themes";

export const Route = createFileRoute("/")({
  loader: async ({ context }) => {
    return context.queryClient.ensureQueryData(schemaQueryOptions);
  },
  component: SchemaEditorPage,
});

function countTokens(schema: string) {
  const clean = schema
    .replace(/\/\/.*/g, "")
    .replace(/\/\*[\s\S]*?\*\//g, "")
    .replace(/\r\n/g, "\n");

  return {
    defs: (clean.match(/\bdefinition\s+\w[\w/]*/g) || []).length,
    rels: (clean.match(/\brelation\s+\w[\w/]*/g) || []).length,
    perms: (clean.match(/\bpermission\s+\w[\w/]*/g) || []).length,
  };
}

function SchemaEditorPage() {
  const queryClient = useQueryClient();
  const { resolvedTheme } = useTheme();
  const { data } = useSuspenseQuery(schemaQueryOptions);
  const schemaText = data?.schemaText ?? "";

  const [value, setValue] = useState(schemaText);
  const [isDirty, setIsDirty] = useState(false);

  useEffect(() => {
    if (!isDirty) {
      setValue(schemaText);
    }
  }, [schemaText, isDirty]);

  const mutation = useMutation({
    mutationFn: async (newSchema: string) => {
      const res = await writeSchema({ data: { schema: newSchema } });
      if (!res.success) throw new Error(res.message);
      return res;
    },
    onSuccess: () => {
      toast.success("Schema publicado com sucesso");
      setIsDirty(false);
      queryClient.invalidateQueries(schemaQueryOptions);
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : "Erro ao publicar o schema");
    },
  });

  const count = countTokens(value)

  const defCount = count.defs
  const relCount = count.rels;
  const permCount = count.perms;

  const stats = [
    { label: "Definições", value: defCount, icon: Database },
    { label: "Relações", value: relCount, icon: Users },
    { label: "Permissões", value: permCount, icon: ShieldCheck },
  ];

  function handleChange(newValue: string | undefined) {
    const val = newValue ?? "";
    setValue(val);
    setIsDirty(val !== schemaText);
  }

  function handlePublish() {
    mutation.mutate(value);
  }

  function handleRefresh() {
    if (isDirty) {
      const confirmDiscard = window.confirm(
        "Você tem alterações não publicadas. Deseja realmente descartá-las e atualizar o schema?"
      );
      if (!confirmDiscard) return;
    }

    queryClient.invalidateQueries(schemaQueryOptions).then(() => {
      setIsDirty(false);
      toast.info("Schema atualizado");
    });
  }

  const monacoTheme = resolvedTheme === "dark" ? "zed-dark" : "zed-light";

  return (
    <div className="flex flex-col h-screen">
      <div className="flex items-center gap-2.5 px-4 py-3 border-b">
        {stats.map(({ label, value: v, icon: Icon }) => (
          <div key={label} className="bg-card border border-border rounded-lg px-4 py-2.5 flex items-center gap-3">
            <div className="p-1.5 bg-muted rounded-md">
              <Icon className="w-4 h-4 text-muted-foreground" />
            </div>
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">
                {label}
              </span>
              <span className="text-xl font-semibold tracking-tight">{v}</span>
            </div>
          </div>
        ))}
        <div className="ml-auto">
          <button
            onClick={handleRefresh}
            className="p-2 hover:bg-muted rounded-md transition-colors"
            title="Atualizar schema"
          >
            <RefreshCw className="w-4 h-4 text-muted-foreground" />
          </button>
        </div>
      </div>
      <header className="flex items-center justify-between px-4 py-2 border-b">
        <span className="text-sm text-muted-foreground">
          {isDirty ? "alterações não publicadas" : "sem alterações"}
        </span>
        <button
          onClick={handlePublish}
          className="bg-primary text-primary-foreground px-3 py-1.5 rounded-md text-sm font-medium disabled:opacity-50"
          disabled={!isDirty || mutation.isPending}
        >
          {mutation.isPending ? "Publicando..." : "Publicar"}
        </button>
      </header>

      <SchemaEditor
        value={value}
        onChange={handleChange}
        theme={monacoTheme}
      />
    </div>
  );
}
