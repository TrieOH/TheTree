import { createSignal } from "solid-js";
import type { OAuthProviderI } from "@trieoh/identityx-sdk-ts";
import { AuthLayout } from "./Shared";
import { ModernSignIn } from "./ModernSignIn";
import { ModernSignUp } from "./ModernSignUp";
import { ModernForgotPassword } from "./ModernForgotPassword";

export type AuthView = "signin" | "signup" | "forgot-password";

export interface ModernAuthProps {
  initialView?: AuthView;
  onLoginSuccess?: (message?: string) => Promise<void>;
  onSignUpSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  backLink?: any;
  providers?: OAuthProviderI[];
}

const viewConfig: Record<AuthView, { title: string; subtitle: string; toggleLabel: string; toggleAction: string; toggleTo: AuthView }> = {
  signin: {
    title: "Bem-vindo de volta",
    subtitle: "Entre com suas credenciais para acessar a plataforma",
    toggleLabel: "Ainda não tem uma conta?",
    toggleAction: "Cadastre-se",
    toggleTo: "signup",
  },
  signup: {
    title: "Criar sua conta",
    subtitle: "Preencha os dados abaixo para se cadastrar",
    toggleLabel: "Já possui uma conta?",
    toggleAction: "Faça login",
    toggleTo: "signin",
  },
  "forgot-password": {
    title: "Recuperar senha",
    subtitle: "Digite seu e-mail para receber instruções",
    toggleLabel: "Lembrou da senha?",
    toggleAction: "Voltar ao login",
    toggleTo: "signin",
  },
};

export function ModernAuth(props: ModernAuthProps) {
  const [view, setView] = createSignal<AuthView>(props.initialView ?? "signin");

  const signupSuccess = async (message?: string) => {
    setView("signin");
    await props.onSignUpSuccess?.(message);
  };

  return (
    <AuthLayout backLink={props.backLink as never}>
      <div class="w-full max-w-md z-10 flex flex-col">
        <div class="text-center mb-6">
          <h1 class="font-heading text-3xl font-bold tracking-tight mb-2">
            {viewConfig[view()].title}
          </h1>
          <p class="text-muted-foreground text-sm mb-4">
            {viewConfig[view()].subtitle}
          </p>

          <div class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-muted/50 border border-border/50 text-xs font-medium text-muted-foreground">
            {viewConfig[view()].toggleLabel}
            <button
              type="button"
              onClick={() => setView(viewConfig[view()].toggleTo)}
              class="text-primary font-semibold hover:underline"
            >
              {viewConfig[view()].toggleAction}
            </button>
          </div>
        </div>

        <div class="relative min-h-112">
          <div class="absolute inset-0">
            {view() === "signin" && (
              <ModernSignIn
                onSuccess={props.onLoginSuccess}
                onFailed={props.onFailed}
                providers={props.providers}
                signUpRedirect={() => setView("signup")}
                forgotPasswordRedirect={() => setView("forgot-password")}
              />
            )}
            {view() === "signup" && (
              <ModernSignUp
                onSuccess={signupSuccess}
                onFailed={props.onFailed}
                signInRedirect={() => setView("signin")}
              />
            )}
            {view() === "forgot-password" && (
              <ModernForgotPassword
                onFailed={props.onFailed}
                signInRedirect={() => setView("signin")}
              />
            )}
          </div>
        </div>
      </div>
    </AuthLayout>
  );
}
