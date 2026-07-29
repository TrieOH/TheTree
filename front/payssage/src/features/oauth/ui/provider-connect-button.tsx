import { useMutation } from "@tanstack/react-query";
import { PlugZap } from "lucide-react";
import { toast } from "sonner";
import { env } from "#/env";
import { Button } from "#/shared/ui/shadcn/button";
import { connectProviderFn, getProviderCallbackUrl } from "../api";

const providerDetails: Partial<
  Record<string, { label: string; logo?: string }>
> = {
  mercadopago: {
    label: "Mercado Pago",
    logo: "/external-logos/MP_RGB_HANDSHAKE_color_vertical.svg",
  },
};

const providerLabel = (provider: string) =>
  providerDetails[provider]?.label ??
  provider
    .replaceAll("_", " ")
    .replace(/\b\w/g, (letter) => letter.toUpperCase());

interface ButtonProps {
  provider: string;
  flow: "collector" | "seller";
  organizationId?: string;
  walletId?: string;
}

function ProviderConnectButton({
  provider,
  flow,
  organizationId,
  walletId,
}: ButtonProps) {
  const details = providerDetails[provider];
  const label = providerLabel(provider);
  // biome-ignore format: stay
  const { mutate, isPending } = useMutation({
    mutationFn: () => {
      if (flow === "seller" && !walletId) {
        throw new Error("Wallet ID is required for seller flow");
      }
      return connectProviderFn(
        provider,
        flow === "seller"
          ? {
            flow,
            wallet_id: walletId as string,
            provider_redirect_url: getProviderCallbackUrl(provider),
            final_redirect_url: window.location.href,
          }
          : {
            flow,
            organization_id: organizationId,
            provider_redirect_url: getProviderCallbackUrl(provider),
            final_redirect_url: window.location.href,
          },
      );
    },
    onSuccess: (response) => {
      if (response.success) window.location.assign(response.data);
      else toast.error(response.message || `Failed to connect ${label}`);
    },
    onError: () => toast.error(`Failed to connect ${label}`),
  });

  return (
    <Button
      variant="outline"
      className="h-12 justify-start gap-3 px-4"
      disabled={isPending}
      onClick={() => mutate()}
    >
      {details?.logo ? (
        <img src={details.logo} alt="" className="size-8 object-contain" />
      ) : (
        <span className="flex size-8 items-center justify-center rounded-sm bg-primary/10 text-primary">
          <PlugZap className="size-4" />
        </span>
      )}
      <span className="flex flex-col items-start leading-tight">
        <span>{isPending ? "Connecting..." : label}</span>
        <span className="text-[10px] font-normal text-muted-foreground">
          Connect as {flow}
        </span>
      </span>
    </Button>
  );
}

interface SectionProps {
  flow: "collector" | "seller";
  organizationId?: string;
  walletId?: string;
}

export function ProviderConnectSection(props: SectionProps) {
  return (
    <section className="space-y-3 rounded-sm border border-dashed bg-muted/20 p-4">
      <div>
        <h2 className="text-sm font-semibold">Connect a provider</h2>
        <p className="text-xs text-muted-foreground">
          Choose the payment provider account you want to connect.
        </p>
      </div>
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {env.VITE_SUPPORTED_PROVIDERS.map((provider) => (
          <ProviderConnectButton
            key={provider}
            provider={provider}
            {...props}
          />
        ))}
      </div>
    </section>
  );
}
