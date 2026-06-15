import { useState } from 'react';
import { toast } from 'sonner';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { ArrowRight, Loader2 } from 'lucide-react';
import { useAuth } from "../../AuthProvider";
import FormInput from './Shared/FormInput';
import FormError from './Shared/FormError';
import { Button } from './Shared/Button';

const loginSchema = z.object({
  email: z.email("E-mail inválido"),
  password: z.string().min(1, "A senha é obrigatória"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export interface ModernSignInProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signUpRedirect?: () => void;
  forgotPasswordRedirect?: () => void;
}

export function ModernSignIn({
  onSuccess,
  onFailed,
  signUpRedirect,
  forgotPasswordRedirect,
}: ModernSignInProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  });

  const onSubmit = async (data: LoginFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.login(data.email, data.password);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else toast.success("Login realizado com sucesso!");
      } else {
        if (onFailed) await onFailed(response.message, response.trace);
        else toast.error(response.message || 'Erro ao entrar');
      }
    } catch {
      toast.error('Ocorreu um erro inesperado');
    } finally {
      setIsLoading(false);
    }
  };

  return (
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
          autoComplete="current-password"
          error={!!errors.password}
          {...register('password')}
        />
        <FormError message={errors.password?.message} />
      </div>

      {forgotPasswordRedirect && (
        <div className="flex justify-end">
          <button
            type="button"
            onClick={forgotPasswordRedirect}
            className="text-xs text-primary hover:underline font-medium"
          >
            Esqueceu a senha?
          </button>
        </div>
      )}

      <Button
        type="submit"
        disabled={isLoading}
        className="w-full flex items-center justify-center gap-2"
      >
        {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
          <>
            Entrar
            <ArrowRight className="w-4 h-4" />
          </>
        )}
      </Button>
    </form>
  );
}