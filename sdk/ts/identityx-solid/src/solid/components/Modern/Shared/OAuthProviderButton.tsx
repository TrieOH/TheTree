import type { OAuthProviderI } from "@trieoh/identityx-sdk-ts";

export interface OAuthProviderButtonProps {
  provider: OAuthProviderI;
  onClick: () => void;
  isLoading?: boolean;
}

export function OAuthProviderButton(props: OAuthProviderButtonProps) {
  const label = () => (props.provider === "github" ? "GitHub" : "Google");
  return (
    <button
      type="button"
      onClick={props.onClick}
      disabled={props.isLoading}
      class="w-full flex items-center justify-center gap-2.5 h-12 px-4 rounded-sm text-sm font-semibold border border-input bg-background hover:bg-muted/60 transition-all disabled:opacity-50"
    >
      {props.isLoading ? "Conectando..." : `Continuar com ${label()}`}
    </button>
  );
}
