import { useState } from 'react';
import { AnimatePresence, motion, } from "motion/react";
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

  return (
    <AuthLayout className="min-h-dvh">
      <AnimatePresence mode="wait">
        {view === "signin" && (
          <motion.div
            key="signin"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 20 }}
            transition={{ duration: 0.3 }}
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
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{ duration: 0.3 }}
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
            transition={{ duration: 0.3 }}
          >
            <ModernForgotPassword
              onFailed={onFailed}
              signInRedirect={() => setView("signin")}
            />
          </motion.div>
        )}
      </AnimatePresence>
    </AuthLayout>
  );
}
