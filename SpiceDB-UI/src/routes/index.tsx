import { readSchema, writeSchema } from "#/features/schema/api";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { toast } from "sonner";
import SchemaEditor from "#/features/schema/ui/schema-editor";
import { Database, ShieldCheck, Users } from "lucide-react";

export const Route = createFileRoute("/")({
  loader: async () => await readSchema(),
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
  const { schemaText } = Route.useLoaderData();
  const [value, setValue] = useState(schemaText);
  const [isDirty, setIsDirty] = useState(false);
  const [isPublishing, setIsPublishing] = useState(false);

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
    setValue(newValue ?? "");
    setIsDirty(newValue !== schemaText);
  }

  async function handlePublish() {
    try {
      setIsPublishing(true);
      const res = await writeSchema({ data: { schema: value } });
      if (res.success) toast.success("Schema publicado com sucesso");
      else toast.warning(res.message)
      setIsDirty(false);
    } catch {
      toast.error("Erro ao publicar o schema");
    } finally {
      setIsPublishing(false);
    }
  }

  return (
    <div className="flex flex-col h-screen">
      <div className="flex items-center gap-2.5 px-4 py-3 border-b">
        {stats.map(({ label, value: v, icon: Icon }) => (
          <div key={label} className="bg-card border border-border rounded-lg px-4 py-2.5 flex items-center gap-3">
            <div className="p-1.5 bg-muted rounded-md">
              <Icon className="w-4 h-4 text-muted-foreground" />
            </div>
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-medium">{label}</span>
              <span className="text-xl font-semibold tracking-tight">{v}</span>
            </div>
          </div>
        ))}
      </div>
      <header className="flex items-center justify-between px-4 py-2 border-b">
        <span className="text-sm text-muted-foreground">
          {isDirty ? "alterações não publicadas" : "sem alterações"}
        </span>
        <button
          onClick={handlePublish}
          disabled={!isDirty || isPublishing}
        >
          {isPublishing ? "Publicando..." : "Publicar"}
        </button>
      </header>

      <SchemaEditor value={value} onChange={handleChange} />
    </div>
  );
}
