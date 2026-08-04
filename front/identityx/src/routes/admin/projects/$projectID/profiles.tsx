import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useLayoutHeader } from "@trieoh/ui-base";
import { Braces, Save } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import {
  projectProfileSchemaQueryOptions,
  upsertProjectProfileSchemaFn,
} from "@/features/profiles/api";
import { DEFAULT_PROFILE_SCHEMA } from "@/features/profiles/model";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

export const Route = createFileRoute("/admin/projects/$projectID/profiles")({
  component: RouteComponent,
});

function RouteComponent() {
  const { projectID } = Route.useParams();
  const queryClient = useQueryClient();
  const query = projectProfileSchemaQueryOptions(projectID);
  const { data, isLoading } = useQuery(query);
  const [schemaText, setSchemaText] = useState(() =>
    JSON.stringify(DEFAULT_PROFILE_SCHEMA, null, 2),
  );
  const [active, setActive] = useState(true);

  useEffect(() => {
    if (!data) return;
    setSchemaText(JSON.stringify(data.schema, null, 2));
    setActive(data.active);
  }, [data]);

  const save = useMutation({
    mutationFn: async () => {
      const schema: unknown = JSON.parse(schemaText);
      if (!schema || typeof schema !== "object" || Array.isArray(schema)) {
        throw new Error("The profile schema must be a JSON object");
      }
      return upsertProjectProfileSchemaFn(projectID, {
        schema: schema as Record<string, unknown>,
        active,
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: query.queryKey });
      toast.success("Profile schema saved");
    },
    onError: (error: Error) =>
      toast.error(
        error instanceof SyntaxError ? "Invalid JSON schema" : error.message,
      ),
  });

  const header = useMemo(
    () => (
      <div>
        <h1 className="text-lg font-semibold">Profiles</h1>
        <p className="text-sm text-muted-foreground">
          Define the JSON Schema used to validate and render actor profiles.
        </p>
      </div>
    ),
    [],
  );
  useLayoutHeader(header);

  return (
    <div className="mx-auto max-w-4xl space-y-4">
      <div className="rounded-md border border-border bg-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
          <div className="flex items-center gap-2">
            <Braces className="size-4 text-primary" />
            <div>
              <p className="text-sm font-semibold">Profile JSON Schema</p>
              <p className="text-xs text-muted-foreground">
                Draft 2020-12 · version {data?.version ?? "new"}
              </p>
            </div>
          </div>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="checkbox"
              checked={active}
              onChange={(event) => setActive(event.target.checked)}
              className="size-4 accent-primary"
            />
            Active
          </label>
        </div>
        <textarea
          aria-label="Profile JSON Schema"
          value={schemaText}
          onChange={(event) => setSchemaText(event.target.value)}
          disabled={isLoading}
          spellCheck={false}
          className="min-h-120 w-full resize-y bg-transparent p-4 font-mono text-xs leading-5 outline-none"
        />
      </div>
      <div className="flex justify-end">
        <ShadowButton
          onClick={() => save.mutate()}
          disabled={save.isPending || isLoading}
          variant="accent-solid"
          className="h-9 rounded-sm px-4"
          leftIcon={<Save size={16} />}
          value={save.isPending ? "Saving..." : "Save schema"}
        />
      </div>
    </div>
  );
}
