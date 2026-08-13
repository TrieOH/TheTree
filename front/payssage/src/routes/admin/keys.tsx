import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useServerFn } from "@tanstack/react-start";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { Check, Copy, Key, Plus } from "lucide-react";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import {
  createApiKeyServerFn,
  listCapabilitiesServerFn,
} from "#/features/keys/api";
import { cn } from "#/shared/lib/utils";
import { Badge } from "#/shared/ui/shadcn/badge";
import { Button } from "#/shared/ui/shadcn/button";
import { Card, CardContent } from "#/shared/ui/shadcn/card";
import { Input } from "#/shared/ui/shadcn/input";

export const Route = createFileRoute("/admin/keys")({
  component: RouteComponent,
});

type CapabilityLike = {
  id?: string;
  resource?: string;
  action?: string;
  name?: string;
};

function RouteComponent() {
  const queryClient = useQueryClient();
  const { auth } = useAuth();
  const projectId = auth.profile()?.project_id ?? undefined;
  const subjectId = auth.profile()?.id ?? undefined;
  const [name, setName] = useState("");
  const [selectedCapabilities, setSelectedCapabilities] = useState<string[]>(
    [],
  );
  const [copiedId, setCopiedId] = useState<string | null>(null);
  const listCapabilities = useServerFn(listCapabilitiesServerFn);
  const createApiKey = useServerFn(createApiKeyServerFn);

  const { data: capabilities = [] } = useQuery({
    queryKey: ["identityx", projectId, "capabilities"],
    queryFn: () => listCapabilities({ data: { projectId: projectId ?? "" } }),
    enabled: !!projectId,
  });
  const [capabilitySearch, setCapabilitySearch] = useState("");

  const createMutation = useMutation({
    mutationFn: async () => {
      if (!projectId || !subjectId) throw new Error("Missing actor context");
      return createApiKey({
        data: {
          projectId,
          subjectId,
          payload: {
            name,
            capabilities: selectedCapabilities,
          },
        },
      });
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({
        queryKey: ["identityx", projectId, "api_keys"],
      });
      toast.success("API key created");
      setName("");
      setSelectedCapabilities([]);
      navigator.clipboard.writeText(data.raw_key);
      setCopiedId(data.key?.id ?? "new-key");
    },
    onError: (error: Error) => toast.error(error.message),
  });

  const apiKey = createMutation.data?.key;

  const normalizedCapabilities = useMemo(() => {
    return capabilities
      .map((capability) => {
        const item = capability as CapabilityLike;
        const id = String(
          item.id ?? item.name ?? item.resource ?? item.action ?? capability,
        );
        const label =
          item.name ??
          ([item.resource, item.action].filter(Boolean).join(":") || id);

        return {
          id,
          label,
        };
      })
      .filter(
        (capability, index, self) =>
          self.findIndex((item) => item.id === capability.id) === index,
      );
  }, [capabilities]);

  const visibleCapabilities = useMemo(() => {
    const search = capabilitySearch.trim().toLowerCase();
    if (!search) return normalizedCapabilities;
    return normalizedCapabilities.filter(
      (capability) =>
        capability.label.toLowerCase().includes(search) ||
        capability.id.toLowerCase().includes(search),
    );
  }, [capabilitySearch, normalizedCapabilities]);

  const selectedCount = selectedCapabilities.length;
  const allVisibleSelected =
    visibleCapabilities.length > 0 &&
    visibleCapabilities.every((capability) =>
      selectedCapabilities.includes(capability.id),
    );

  const toggleCapability = (capabilityId: string) => {
    setSelectedCapabilities((current) =>
      current.includes(capabilityId)
        ? current.filter((item) => item !== capabilityId)
        : [...current, capabilityId],
    );
  };

  return (
    <div className="py-4 px-6 space-y-6">
      <div className="flex items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold tracking-tight">API Keys</h1>
          <p className="text-sm text-muted-foreground">Your keys</p>
        </div>
      </div>

      <Card>
        <CardContent className="p-6 space-y-4">
          <div className="space-y-2">
            <span className="text-[10px] font-black uppercase tracking-[0.2em]">
              Key name
            </span>
            <Input
              value={name}
              onChange={(event) => setName(event.target.value)}
              placeholder="e.g. Checkout prod key"
              className="rounded-none"
            />
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between gap-3">
              <div>
                <p className="text-[10px] font-black uppercase tracking-[0.2em]">
                  Capabilities
                </p>
                <p className="text-xs text-muted-foreground">
                  Select the capabilities this key should carry.
                </p>
              </div>
              <Badge
                variant="outline"
                className="rounded-none uppercase tracking-widest"
              >
                {selectedCount} selected
              </Badge>
            </div>

            <div className="flex flex-col gap-3 rounded-none border border-border bg-muted/20 p-3">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Input
                  value={capabilitySearch}
                  onChange={(event) => setCapabilitySearch(event.target.value)}
                  placeholder="Search capabilities..."
                  className="rounded-none sm:max-w-sm"
                />
                <div className="flex items-center gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="rounded-none"
                    onClick={() => {
                      setSelectedCapabilities((current) => {
                        const merged = new Set(current);
                        for (const capability of visibleCapabilities) {
                          merged.add(capability.id);
                        }
                        return [...merged];
                      });
                    }}
                    disabled={
                      visibleCapabilities.length === 0 || allVisibleSelected
                    }
                  >
                    Select visible
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="rounded-none"
                    onClick={() => setSelectedCapabilities([])}
                    disabled={selectedCapabilities.length === 0}
                  >
                    Clear
                  </Button>
                </div>
              </div>

              <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                {visibleCapabilities.length === 0 ? (
                  <div className="col-span-full rounded-none border border-dashed border-border bg-background px-4 py-8 text-center text-xs uppercase tracking-widest text-muted-foreground">
                    No capabilities found.
                  </div>
                ) : (
                  visibleCapabilities.map((capability) => {
                    const active = selectedCapabilities.includes(capability.id);

                    return (
                      <button
                        key={capability.id}
                        type="button"
                        onClick={() => toggleCapability(capability.id)}
                        className={cn(
                          "group flex items-center justify-between gap-3 border px-3 py-3 text-left transition-all",
                          active
                            ? "border-primary bg-primary/5 text-primary shadow-[0_0_0_1px_rgba(0,0,0,0.02)]"
                            : "border-border bg-background text-foreground hover:border-primary/50 hover:bg-primary/5/20",
                        )}
                      >
                        <div className="min-w-0 flex items-center gap-2">
                          <span
                            className={cn(
                              "inline-flex size-4 items-center justify-center border text-[10px] font-black shrink-0",
                              active
                                ? "border-primary bg-primary text-primary-foreground"
                                : "border-border text-muted-foreground",
                            )}
                          >
                            {active ? "✓" : ""}
                          </span>
                          <span className="truncate text-sm font-semibold">
                            {capability.label}
                          </span>
                        </div>
                      </button>
                    );
                  })
                )}
              </div>
            </div>
          </div>

          <div className="flex justify-end">
            <Button
              type="button"
              onClick={() => createMutation.mutate()}
              disabled={createMutation.isPending || !name.trim()}
              className="rounded-none gap-2 uppercase tracking-widest font-black"
            >
              <Plus className="size-4" />
              Create Key
            </Button>
          </div>
        </CardContent>
      </Card>

      {apiKey ? (
        <Card>
          <CardContent className="p-6 space-y-4">
            <div className="flex items-center gap-2">
              <Key className="size-4 text-primary" />
              <h2 className="text-sm font-semibold uppercase tracking-[0.2em]">
                New API Key
              </h2>
            </div>
            <div className="rounded-none border border-border bg-muted/40 p-4 font-mono text-sm break-all">
              {createMutation.data?.raw_key}
            </div>
            <div className="flex flex-wrap items-center gap-2 text-xs">
              <Badge variant="outline" className="rounded-none">
                {apiKey.display_prefix}
              </Badge>
              <span className="text-muted-foreground">
                {copiedId === apiKey.id
                  ? "Copied"
                  : "Copy immediately, this value is shown once."}
              </span>
              <Button
                type="button"
                size="sm"
                variant="outline"
                className="rounded-none ml-auto gap-2"
                onClick={() => {
                  const rawKey = createMutation.data?.raw_key ?? null;
                  if (!rawKey) return;
                  navigator.clipboard.writeText(rawKey);
                  setCopiedId(apiKey.id);
                }}
              >
                {copiedId === apiKey.id ? (
                  <Check className="size-4" />
                ) : (
                  <Copy className="size-4" />
                )}
                Copy
              </Button>
            </div>
          </CardContent>
        </Card>
      ) : null}
    </div>
  );
}
