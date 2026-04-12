import { schemaQueryOptions, writeSchema } from "#/features/schema/api";
import { createFileRoute } from "@tanstack/react-router";
import { useState, useEffect, useRef, useCallback } from "react";
import { toast } from "sonner";
import SchemaEditor from "#/features/schema/ui/schema-editor";
import { Database, ShieldCheck, Users, RefreshCw, AlertTriangle } from "lucide-react";
import { useMutation, useQueryClient, useQuery } from "@tanstack/react-query";
import { useTheme } from "next-themes";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "#/shared/ui/shadcn/alert-dialog";
import type { editor } from "monaco-editor";
import { cn } from "#/shared/lib/class-utils";
import { Button } from "#/shared/ui/shadcn/button";

export const Route = createFileRoute("/")({
  loader: async ({ context }) => {
    return context.queryClient.fetchQuery(schemaQueryOptions);
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
  const loaderData = Route.useLoaderData();
  const queryClient = useQueryClient();
  const { resolvedTheme } = useTheme();

  const { data } = useQuery({
    ...schemaQueryOptions,
    initialData: loaderData,
    refetchOnMount: false,
  });

  const schemaText = data.schemaText;

  const [value, setValue] = useState(schemaText);
  const [isDirty, setIsDirty] = useState(false);
  const [showConfirmRefresh, setShowConfirmRefresh] = useState(false);
  const [cursorInfo, setCursorInfo] = useState({ line: 1, column: 1 });

  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
  const serverSchemaRef = useRef(schemaText);

  useEffect(() => {
    if (schemaText === serverSchemaRef.current) return;
    serverSchemaRef.current = schemaText;

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
    onSuccess: (_, newSchema) => {
      queryClient.setQueryData(schemaQueryOptions.queryKey, (old) => {
        if (!old) return old
        return {
          ...old,
          schemaText: newSchema,
        }
      })
      serverSchemaRef.current = newSchema;
      setValue(newSchema);
      setIsDirty(false);

      toast.success("Schema publicado com sucesso");
      queryClient.invalidateQueries(schemaQueryOptions);
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : "Erro ao publicar o schema");
    },
  });

  const count = countTokens(value);

  const stats = [
    { label: "Defs", count: count.defs, icon: Database },
    { label: "Rels", count: count.rels, icon: Users },
    { label: "Perms", count: count.perms, icon: ShieldCheck },
  ];

  const handleChange = useCallback((newValue: string | undefined) => {
    const val = newValue ?? "";
    setValue(val);
    setIsDirty(val !== serverSchemaRef.current);
  }, []);

  const handleEditorMount = useCallback((ed: editor.IStandaloneCodeEditor) => {
    editorRef.current = ed;
    ed.onDidChangeCursorPosition((e) => {
      setCursorInfo({ line: e.position.lineNumber, column: e.position.column });
    });
  }, []);

  function handlePublish() {
    mutation.mutate(value);
  }

  async function performRefresh() {
    try {
      const freshData = await queryClient.fetchQuery({
        ...schemaQueryOptions,
        staleTime: 0,
      });

      serverSchemaRef.current = freshData.schemaText;
      setValue(freshData.schemaText);

      setIsDirty(false);
      setShowConfirmRefresh(false);
      toast.info("Schema atualizado");
    } catch (err) {
      toast.error("Erro ao atualizar o schema");
    }
  }

  function handleRefresh() {
    if (isDirty) setShowConfirmRefresh(true);
    else performRefresh();
  }

  const monacoTheme = resolvedTheme === "dark" ? "zed-dark" : "zed-light";

  return (
    <main className="flex flex-col h-screen min-w-75 border-l">
      {/* ... (stats and refresh button code) */}
      <div className="flex items-center gap-px border-b">
        {stats.map(({ label, count: c, icon: Icon }, i) => (
          <div
            key={label}
            className={[
              "flex items-center gap-2.5 px-4 py-3 flex-1 min-w-0",
              i < stats.length - 1 ? "border-r border-border" : "",
            ].join(" ")}
          >
            <Icon className="w-4 h-4 text-muted-foreground shrink-0" />
            <div className="flex items-baseline gap-1.5 min-w-0">
              <span className="text-xl font-semibold tabular-nums leading-none">{c}</span>
              <span className="text-xs text-muted-foreground truncate">{label}</span>
            </div>
          </div>
        ))}
        <button
          onClick={handleRefresh}
          className="px-4 py-3 border-l border-border hover:bg-muted transition-colors shrink-0 cursor-pointer"
          title="Atualizar schema"
        >
          <RefreshCw className="w-4 h-4 text-muted-foreground" />
        </button>
      </div>

      <div className="flex items-center justify-between px-3 py-1.5 border-b gap-3">
        {isDirty ? (
          <span className="inline-flex items-center gap-1.5 text-[11px] font-medium text-accent-foreground/80 bg-accent/20 border border-accent/30 rounded px-2 py-0.5 leading-5">
            <span className="w-1.5 h-1.5 rounded-full bg-accent shrink-0" />
            não publicado
          </span>
        ) : (
          <span className="text-[11px] text-muted-foreground leading-5">publicado</span>
        )}
        <Button
          onClick={handlePublish}
          disabled={!isDirty || mutation.isPending}
          className={cn(
            "bg-primary text-primary-foreground px-3 py-1 rounded",
            "text-xs font-medium disabled:opacity-40 transition-opacity shrink-0",
            "cursor-pointer active:scale-95",
            "duration-200 transition-transform"
          )}
        >
          {mutation.isPending ? "Publicando…" : "Publicar"}
        </Button>
      </div>

      <div className="flex-1 min-h-0">
        <SchemaEditor
          value={value}
          onChange={handleChange}
          onMount={handleEditorMount}
          theme={monacoTheme}
        />
      </div>

      <div className="flex items-center justify-between px-3 py-1 border-t text-[11px] text-muted-foreground">
        <span className="tabular-nums">Ln {cursorInfo.line}, Col {cursorInfo.column}</span>
        <span>.zed</span>
      </div>

      <AlertDialog open={showConfirmRefresh} onOpenChange={setShowConfirmRefresh}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <div className="flex items-center gap-2 mb-2">
              <AlertTriangle className="w-5 h-5 text-accent" />
              <AlertDialogTitle>Descartar alterações?</AlertDialogTitle>
            </div>
            <AlertDialogDescription>
              Você tem alterações não publicadas. Deseja realmente descartá-las e atualizar o schema original?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancelar</AlertDialogCancel>
            <AlertDialogAction onClick={performRefresh}>
              Descartar e Atualizar
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </main>
  );
}