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

const forgotPasswordSchema = z.object({
  email: z.email("E-mail inválido"),
});

type ForgotPasswordFormValues = z.infer<typeof forgotPasswordSchema>;

export interface ModernForgotPasswordProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernForgotPassword({
  onSuccess,
  onFailed,
  signInRedirect,
}: ModernForgotPasswordProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordFormValues>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: { email: '' },
  });

  const onSubmit = async (data: ForgotPasswordFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.sendForgotPassword(data.email);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else toast.success("Link de redefinição enviado com sucesso!");
      } else {
        if (onFailed) await onFailed(response.message, response.trace);
        else toast.error(response.message || 'Erro ao enviar link');
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
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className="font-heading text-3xl font-bold tracking-tight mb-2"
        >
          Esqueceu a senha?
        </motion.h1>
        <p className="text-muted-foreground text-sm mb-4">
          Insira seu e-mail para receber um link de redefinição
        </p>

        {signInRedirect && (
          <div className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-muted/50 border border-border/50 text-xs font-medium text-muted-foreground">
            Lembrou sua senha?
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

        <Button
          type="submit"
          disabled={isLoading}
          className="w-full flex items-center justify-center gap-2"
        >
          {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
            <>
              Enviar link
              <ArrowRight className="w-4 h-4" />
            </>
          )}
        </Button>
      </form>
    </div>
  );
}
