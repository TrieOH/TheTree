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

const resetPasswordSchema = z.object({
  password: z.string().min(6, "A senha deve ter pelo menos 6 caracteres"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "As senhas não coincidem",
  path: ["confirmPassword"],
});

type ResetPasswordFormValues = z.infer<typeof resetPasswordSchema>;

export interface ModernResetPasswordProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signInRedirect?: () => void;
}

export function ModernResetPassword({
  token,
  onSuccess,
  onFailed,
  signInRedirect,
}: ModernResetPasswordProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetPasswordFormValues>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: { password: '', confirmPassword: '' },
  });

  const onSubmit = async (data: ResetPasswordFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.resetPassword(token, data.password);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else toast.success("Senha redefinida com sucesso!");
      } else {
        if (onFailed) await onFailed(response.message, response.trace);
        else toast.error(response.message || 'Erro ao redefinir senha');
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
          Redefinir Senha
        </motion.h1>
        <p className="text-muted-foreground text-sm mb-4">
          Crie uma nova senha para sua conta
        </p>
      </div>

      <form
        onSubmit={handleSubmit(onSubmit)}
        className="space-y-4"
      >
        <div className="space-y-1">
          <FormInput
            label="Nova Senha"
            type="password"
            autoComplete="new-password"
            error={!!errors.password}
            {...register('password')}
          />
          <FormError message={errors.password?.message} />
        </div>

        <div className="space-y-1">
          <FormInput
            label="Confirmar Nova Senha"
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
              Redefinir Senha
              <ArrowRight className="w-4 h-4" />
            </>
          )}
        </Button>

        {signInRedirect && (
          <div className="text-center mt-4">
            <button
              type="button"
              onClick={signInRedirect}
              className="text-xs text-primary font-semibold hover:underline"
            >
              Voltar para o login
            </button>
          </div>
        )}
      </form>
    </div>
  );
}
