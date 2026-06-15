import { useState } from 'react';
import { AnimatePresence, motion } from "motion/react";
import { ModernSignIn } from './ModernSignIn';
import { ModernSignUp } from './ModernSignUp';
import { ModernForgotPassword } from './ModernForgotPassword';
import { AuthLayout } from './Shared/AuthLayout';

export type AuthView = "signin" | "signup" | "forgot-password";

export interface ModernAuthProps {
  initialView?: AuthView;
  onLoginSuccess?: (message?: string) => Promise<void>;
  onSignUpSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
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

export function ModernAuth({
  initialView = "signin",
  onLoginSuccess,
  onSignUpSuccess,
  onFailed,
}: ModernAuthProps) {
  const [view, setView] = useState<AuthView>(initialView);

  const handleSignUpSuccess = async (message?: string) => {
    setView("signin");
    if (onSignUpSuccess) await onSignUpSuccess(message);
  };

  const config = viewConfig[view];

  return (
    <AuthLayout>
      <div className="w-full max-w-md z-10 flex flex-col">
        <div className="text-center mb-6">
          <motion.h1
            key={view + "-h1"}
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            className="font-heading text-3xl font-bold tracking-tight mb-2"
          >
            {config.title}
          </motion.h1>
          <p className="text-muted-foreground text-sm mb-4">
            {config.subtitle}
          </p>

          <div className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-muted/50 border border-border/50 text-xs font-medium text-muted-foreground">
            {config.toggleLabel}
            <button
              type="button"
              onClick={() => setView(config.toggleTo)}
              className="text-primary font-semibold hover:underline"
            >
              {config.toggleAction}
            </button>
          </div>
        </div>

        {/* Animated form area with fixed min-height to prevent layout shift */}
        <div className="relative min-h-80">
          <AnimatePresence mode="wait">
            {view === "signin" && (
              <motion.div
                key="signin"
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: 10 }}
                transition={{ duration: 0.2 }}
                className="absolute inset-0"
              >
                <ModernSignIn
                  onSuccess={onLoginSuccess}
                  onFailed={onFailed}
                  signUpRedirect={() => setView("signup")}
                  forgotPasswordRedirect={() => setView("forgot-password")}
                />
              </motion.div>
            )}

            {view === "signup" && (
              <motion.div
                key="signup"
                initial={{ opacity: 0, x: 10 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -10 }}
                transition={{ duration: 0.2 }}
                className="absolute inset-0"
              >
                <ModernSignUp
                  onSuccess={handleSignUpSuccess}
                  onFailed={onFailed}
                  signInRedirect={() => setView("signin")}
                />
              </motion.div>
            )}

            {view === "forgot-password" && (
              <motion.div
                key="forgot-password"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.2 }}
                className="absolute inset-0"
              >
                <ModernForgotPassword
                  onFailed={onFailed}
                  signInRedirect={() => setView("signin")}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </AuthLayout>
  );
}