import { useState } from 'react';
import { motion } from "motion/react";
import { toast } from 'sonner';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { ArrowRight, Loader2 } from 'lucide-react';
import { useAuth } from "../../AuthProvider";
import FormInput from './Shared/FormInput';
import FormError from './Shared/FormError';
import { Button } from './Shared/Button';

const signUpSchema = z.object({
  email: z.string().email("E-mail inválido"),
  password: z.string().min(6, "A senha deve ter pelo menos 6 caracteres"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "As senhas não coincidem",
  path: ["confirmPassword"],
});

type SignUpFormValues = z.infer<typeof signUpSchema>;

export interface ModernSignUpProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernSignUp({
  onSuccess,
  onFailed,
  signInRedirect,
}: ModernSignUpProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpFormValues>({
    resolver: zodResolver(signUpSchema),
    defaultValues: { email: '', password: '', confirmPassword: '' },
  });

  const onSubmit = async (data: SignUpFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.register(data.email, data.password);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else {
          toast.success("Conta criada com sucesso!");
          if (signInRedirect) signInRedirect();
        }
      } else {
        if (onFailed) await onFailed(response.message, response.trace);
        else toast.error(response.message || 'Erro ao criar conta');
      }
    } catch {
      toast.error('Ocorreu um erro inesperado');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md z-10 flex flex-col antialiased selection:bg-primary/10 selection:text-primary">
      <div className="text-center mb-6">
        <motion.h1
          initial={{ opacity: 0, x: 10 }}
          animate={{ opacity: 1, x: 0 }}
          className="font-heading text-3xl font-bold tracking-tight mb-2"
        >
          Criar sua conta
        </motion.h1>
        <p className="text-muted-foreground text-sm mb-4">
          Preencha os dados abaixo para se cadastrar
        </p>

        {signInRedirect && (
          <div className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-muted/50 border border-border/50 text-xs font-medium text-muted-foreground">
            Já possui uma conta?
            <button
              type="button"
              onClick={signInRedirect}
              className="text-primary font-semibold hover:underline"
            >
              Faça login
            </button>
          </div>
        )}
      </div>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <div className="space-y-1">
          <FormInput
            label="E-mail"
            type="email"
            autoComplete="email"
            error={!!errors.email}
            {...register('email')}
          />
          <FormError message={errors.email?.message} />
        </div>

        <div className="space-y-1">
          <FormInput
            label="Senha"
            type="password"
            autoComplete="new-password"
            error={!!errors.password}
            {...register('password')}
          />
          <FormError message={errors.password?.message} />
        </div>

        <div className="space-y-1">
          <FormInput
            label="Confirmar Senha"
            type="password"
            autoComplete="new-password"
            error={!!errors.confirmPassword}
            {...register('confirmPassword')}
          />
          <FormError message={errors.confirmPassword?.message} />
        </div>

        <Button
          type="submit"
          disabled={isLoading}
          className="w-full flex items-center justify-center gap-2"
        >
          {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
            <>
              Criar Conta
              <ArrowRight className="w-4 h-4" />
            </>
          )}
        </Button>
      </form>
    </div>
  );
}
