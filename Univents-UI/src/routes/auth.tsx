import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@soramux/identityx-sdk-ts/react'
import { useState } from 'react';
import { motion, AnimatePresence } from "motion/react";
import { toast } from 'sonner';
import z from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowRight, Loader2 } from 'lucide-react';
import type { LoginFormValues, SignUpFormValues } from '@/features/auths/model';
import { requireGuest } from '@/features/auths/lib/route-guard';
import { useAuthActions } from '@/features/auths/hooks/use-auth-actions';
import { cn } from '@/shared/lib/utils';
import FormInput from '@/shared/ui/form/FormInput';
import { loginSchema, signUpSchema } from '@/features/auths/model';
import { Button } from '@/shared/ui/shadcn/button';
import FormError from '@/shared/ui/form/FormError';

const authSearchSchema = z.object({
  redirect: z.string().optional().catch(''),
})

export const Route = createFileRoute('/auth')({
  validateSearch: (search) => authSearchSchema.parse(search),
  beforeLoad: requireGuest,
  component: App,
})

function App() {
  const [isLogin, setIsLogin] = useState(true);
  const [isLoading, setIsLoading] = useState(false);
  const search = Route.useSearch();
  const { auth } = useAuth();
  const { handleLoginSuccess } = useAuthActions();

  const loginForm = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  });

  const signUpForm = useForm<SignUpFormValues>({
    resolver: zodResolver(signUpSchema),
    defaultValues: { email: '', password: '', confirmPassword: '' },
  });

  const onLoginSubmit = async (data: LoginFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.login(data.email, data.password);
      if (response.success) await handleLoginSuccess(search.redirect);
      else toast.error(response.message || 'Erro ao entrar');
    } catch { toast.error('Ocorreu um erro inesperado'); }
    finally { setIsLoading(false); }
  };

  const onSignUpSubmit = async (data: SignUpFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.register(
        data.email,
        data.password,
      );
      if (response.success) {
        toast.success("Conta criada com sucesso! Faça login para continuar.");
        setIsLogin(true);
      } else toast.error(response.message || 'Erro ao criar conta');
    } catch { toast.error('Ocorreu um erro inesperado'); }
    finally { setIsLoading(false); }
  };

  const hasErrors = isLogin
    ? Object.keys(loginForm.formState.errors).length > 0
    : Object.keys(signUpForm.formState.errors).length > 0;

  return (
    <main className={cn(
      "bg-background h-full text-foreground min-h-screen relative overflow-hidden",
      "flex justify-center items-center px-4",
      isLogin ? (hasErrors ? "pb-24" : "pb-8") : (hasErrors ? "pb-24" : "pb-8"),
      "antialiased selection:bg-primary/10 selection:text-primary"
    )}>
      {/* Decorative background elements */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-primary/20 rounded-full blur-3xl pointer-events-none" />

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md z-10 flex flex-col"
      >
        <div className="text-center mb-6">
          <motion.h1
            key={isLogin ? 'h1-login' : 'h1-signup'}
            initial={{ opacity: 0, x: -10 }}
            animate={{ opacity: 1, x: 0 }}
            className="font-heading text-3xl font-bold tracking-tight mb-2"
          >
            {isLogin ? 'Bem-vindo de volta' : 'Criar sua conta'}
          </motion.h1>
          <p className="text-muted-foreground text-sm mb-4">
            {isLogin
              ? 'Entre com suas credenciais para acessar a plataforma'
              : 'Preencha os dados abaixo para se cadastrar'}
          </p>

          <div className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-muted/50 border border-border/50 text-xs font-medium text-muted-foreground">
            {isLogin ? 'Ainda não tem uma conta?' : 'Já possui uma conta?'}
            <button
              type="button"
              onClick={() => { setIsLogin(!isLogin); }}
              className="text-primary font-semibold hover:underline"
            >
              {isLogin ? 'Cadastre-se' : 'Faça login'}
            </button>
          </div>
        </div>

        <div className="relative min-h-80">
          <AnimatePresence mode="wait">
            {isLogin ? (
              <motion.form
                key="login-form"
                initial={{ opacity: 0, x: -10 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: 10 }}
                transition={{ duration: 0.2 }}
                onSubmit={(e) => {
                  void loginForm.handleSubmit(onLoginSubmit)(e);
                }}
                className="absolute inset-0 space-y-4"
              >
                <div className="space-y-1">
                  <FormInput
                    label="E-mail"
                    type="email"
                    autoComplete="email"
                    error={!!loginForm.formState.errors.email?.message}
                    {...loginForm.register('email')}
                  />
                  <FormError message={loginForm.formState.errors.email?.message} />
                </div>

                <div className="space-y-1">
                  <FormInput
                    label="Senha"
                    type="password"
                    autoComplete="current-password"
                    error={!!loginForm.formState.errors.password?.message}
                    {...loginForm.register('password')}
                  />
                  <FormError message={loginForm.formState.errors.password?.message} />
                </div>

                <div className="flex justify-end">
                  <button type="button" className="text-xs text-primary hover:underline font-medium">
                    Esqueceu a senha?
                  </button>
                </div>

                <Button
                  type="submit"
                  disabled={isLoading}
                  className={cn(
                    "w-full bg-primary text-primary-foreground",
                    "h-12 rounded-sm font-semibold transition-all",
                    "flex items-center justify-center gap-2",
                    "disabled:pointer-events-none disabled:opacity-70",
                    "hover:opacity-90 active:scale-[0.98]"
                  )}
                >
                  {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
                    <>
                      Entrar
                      <ArrowRight className="w-4 h-4" />
                    </>
                  )}
                </Button>
              </motion.form>
            ) : (
              <motion.form
                key="signup-form"
                initial={{ opacity: 0, x: 10 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -10 }}
                transition={{ duration: 0.2 }}
                onSubmit={(e) => {
                  void signUpForm.handleSubmit(onSignUpSubmit)(e);
                }}
                className="absolute inset-0 space-y-4"
              >
                <div className="space-y-1">
                  <FormInput
                    label="E-mail"
                    type="email"
                    autoComplete="email"
                    error={!!signUpForm.formState.errors.email?.message}
                    {...signUpForm.register('email')}
                  />
                  <FormError message={signUpForm.formState.errors.email?.message} />
                </div>

                <div className="space-y-1">
                  <FormInput
                    label="Senha"
                    type="password"
                    autoComplete="new-password"
                    error={!!signUpForm.formState.errors.password?.message}
                    {...signUpForm.register('password')}
                  />
                  <FormError message={signUpForm.formState.errors.password?.message} />
                </div>

                <div className="space-y-1">
                  <FormInput
                    label="Confirmar Senha"
                    type="password"
                    autoComplete="new-password"
                    error={!!signUpForm.formState.errors.confirmPassword?.message}
                    {...signUpForm.register('confirmPassword')}
                  />
                  <FormError message={signUpForm.formState.errors.confirmPassword?.message} />
                </div>

                <Button
                  type="submit"
                  disabled={isLoading}
                  className={cn(
                    "w-full bg-primary text-primary-foreground",
                    "h-12 rounded-sm font-semibold transition-all",
                    "flex items-center justify-center gap-2",
                    "disabled:pointer-events-none disabled:opacity-70",
                    "hover:opacity-90 active:scale-[0.98]"
                  )}
                >
                  {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
                    <>
                      Criar Conta
                      <ArrowRight className="w-4 h-4" />
                    </>
                  )}
                </Button>
              </motion.form>
            )}
          </AnimatePresence>
        </div>
      </motion.div>
    </main>
  )
}
