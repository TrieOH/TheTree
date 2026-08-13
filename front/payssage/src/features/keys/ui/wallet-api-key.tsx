import { useMutation } from "@tanstack/react-query";
import { useServerFn } from "@tanstack/react-start";
import { useAuth } from "@trieoh/identityx-sdk-ts/react";
import { KeyRound } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import type { WalletI } from "#/features/wallets/model";
import { Button } from "#/shared/ui/shadcn/button";
import { createApiKeyServerFn } from "../api";
import type { ApiKeyCreateResponseI } from "../model";
import { ApiKeyCreatedModal } from "./api-key-created-modal";

export function WalletApiKey({ wallet }: { wallet: WalletI }) {
  const { auth } = useAuth();
  const createApiKey = useServerFn(createApiKeyServerFn);
  const [created, setCreated] = useState<ApiKeyCreateResponseI | null>(null);
  const mutation = useMutation({
    mutationFn: () =>
      createApiKey({
        data: {
          projectId: auth.profile()?.project_id ?? "",
          subjectId: wallet.owner_id,
          payload: { name: `${wallet.name} API key`, capabilities: [] },
        },
      }),
    onSuccess: (result) => {
      setCreated({
        id: result.key?.id ?? "new-key",
        name: result.key?.name ?? `${wallet.name} API key`,
        prefix: result.key?.display_prefix ?? "",
        created_at: result.key?.created_at ?? new Date().toISOString(),
        revoked_at: result.key?.revoked_at ?? null,
        key: result.raw_key,
      });
      toast.success("API key criada");
    },
    onError: () => toast.error("Não foi possível criar a API key"),
  });

  return (
    <>
      <Button
        type="button"
        variant="outline"
        className="rounded-none"
        disabled={
          mutation.isPending || !wallet.owner_id || !auth.profile()?.project_id
        }
        onClick={() => mutation.mutate()}
      >
        <KeyRound />
        {mutation.isPending ? "Creating…" : "Create API key"}
      </Button>
      <ApiKeyCreatedModal
        apiKey={created}
        isOpen={created !== null}
        onClose={() => setCreated(null)}
      />
    </>
  );
}
