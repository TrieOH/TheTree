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
  backLink?: unknown;
  providers?: OAuthProviderI[];
}
export function ModernAuth(props: ModernAuthProps) {
  const [view, setView] = createSignal<AuthView>(props.initialView ?? "signin");
  const signupSuccess = async (message?: string) => {
    setView("signin");
    await props.onSignUpSuccess?.(message);
  };
  return (
    <AuthLayout backLink={props.backLink as never}>
      <div class="w-full max-w-md z-10 flex flex-col gap-6">
        <div class="text-center">
          <h1 class="font-heading text-3xl font-bold">
            {view() === "signin"
              ? "Bem-vindo de volta"
              : view() === "signup"
                ? "Criar sua conta"
                : "Recuperar senha"}
          </h1>
          <p class="text-muted-foreground text-sm">
            {view() === "signin"
              ? "Entre com suas credenciais para acessar a plataforma"
              : view() === "signup"
                ? "Preencha os dados abaixo para se cadastrar"
                : "Digite seu e-mail para receber instruções"}
          </p>
          <button
            type="button"
            class="mt-3 text-sm text-primary hover:underline"
            onClick={() => setView(view() === "signin" ? "signup" : "signin")}
          >
            {view() === "signin"
              ? "Ainda não tem uma conta? Cadastre-se"
              : "Voltar ao login"}
          </button>
        </div>
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
    </AuthLayout>
  );
}
