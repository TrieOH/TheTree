import { useState } from 'react';
import { toast } from 'sonner';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { ArrowRight, Loader2, Sparkles } from 'lucide-react';
import { useAuth } from "../../AuthProvider";
import FormInput from './Shared/FormInput';
import FormError from './Shared/FormError';
import { Button } from './Shared/Button';
import { AuthLayout } from './Shared/AuthLayout';
import { motion } from "motion/react";

const setupSchema = z.object({
  email: z.string().email("E-mail inválido"),
  password: z.string().min(8, "A senha deve ter pelo menos 8 caracteres"),
  confirmPassword: z.string().min(1, "A confirmação de senha é obrigatória"),
}).refine((data) => data.password === data.confirmPassword, {
  message: "As senhas não coincidem",
  path: ["confirmPassword"],
});

type SetupFormValues = z.infer<typeof setupSchema>;

export interface ModernSetupProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  backLink?: React.ReactNode;
}

export function ModernSetup({
  onSuccess,
  onFailed,
  backLink,
}: ModernSetupProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SetupFormValues>({
    resolver: zodResolver(setupSchema),
    defaultValues: { email: '', password: '', confirmPassword: '' },
  });

  const onSubmit = async (data: SetupFormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.setup(data.email, data.password);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else toast.success("Configuração inicial realizada com sucesso!");
      } else {
        if (onFailed) await onFailed(response.message, response.trace);
        else toast.error(response.message || 'Erro ao realizar setup');
      }
    } catch (error) {
      console.error(error);
      toast.error('Ocorreu um erro inesperado');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthLayout backLink={backLink}>
      <div className="w-full max-w-md z-10 flex flex-col">
        <div className="mb-10 text-center">
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 text-primary text-xs font-bold uppercase tracking-widest mb-4"
          >
            <Sparkles size={14} />
            Setup Inicial
          </motion.div>
          
          <h1 className="font-heading text-4xl font-bold tracking-tight mb-3">
            Bem-vindo ao IdentityX
          </h1>
          <p className="text-muted-foreground text-base max-w-[320px] mx-auto">
            Vamos começar criando sua conta de administrador mestre.
          </p>
        </div>

        <motion.div
          initial={{ opacity: 0, scale: 0.98 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.1 }}
          className="bg-card/50 backdrop-blur-sm border border-border/40 p-1 rounded-2xl"
        >
          <div className="bg-background rounded-[14px] p-6 shadow-sm">
             <form
              onSubmit={handleSubmit(onSubmit)}
              className="space-y-5"
            >
              <div className="space-y-1">
                <FormInput
                  label="E-mail do Administrador"
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
                className="w-full flex items-center justify-center gap-2 h-13 mt-4 text-base font-bold shadow-lg shadow-primary/10 group"
              >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : (
                  <>
                    Finalizar Configuração
                    <ArrowRight className="w-4 h-4 transition-transform group-hover:translate-x-1" />
                  </>
                )}
              </Button>
            </form>
          </div>
        </motion.div>

        <div className="mt-8 text-center">
           <div className="h-1 w-full bg-muted/30 rounded-full overflow-hidden">
              <motion.div 
                initial={{ width: 0 }}
                animate={{ width: "100%" }}
                transition={{ duration: 1.5, ease: "easeInOut" }}
                className="h-full bg-primary/40"
              />
           </div>
           <p className="text-[10px] text-muted-foreground mt-2 font-medium uppercase tracking-tighter">
              Passo 1 de 1 • Configurando ambiente
           </p>
        </div>
      </div>
    </AuthLayout>
  );
}