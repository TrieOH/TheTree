import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";
import { useAuth } from "../../AuthProvider";
import { Button } from "./Shared/Button";
import FormError from "./Shared/FormError";
import FormInput from "./Shared/FormInput";

const schema = z.object({ email: z.email("E-mail inválido") });
type FormValues = z.infer<typeof schema>;

export interface ModernResendVerificationProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
}

export function ModernResendVerification({
  onSuccess,
  onFailed,
}: ModernResendVerificationProps) {
  const [isLoading, setIsLoading] = useState(false);
  const { auth } = useAuth();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "" },
  });

  const onSubmit = async ({ email }: FormValues) => {
    setIsLoading(true);
    try {
      const response = await auth.resendVerifyEmail(email);
      if (response.success) {
        if (onSuccess) await onSuccess(response.message);
        else toast.success("E-mail de verificação enviado!");
      } else if (onFailed) await onFailed(response.message, response.trace);
      else toast.error(response.message || "Erro ao reenviar verificação");
    } catch {
      toast.error("Ocorreu um erro inesperado");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-1">
        <FormInput
          label="E-mail"
          type="email"
          autoComplete="email"
          error={!!errors.email}
          {...register("email")}
        />
        <FormError message={errors.email?.message} />
      </div>
      <Button
        type="submit"
        disabled={isLoading}
        className="flex w-full items-center justify-center gap-2"
      >
        {isLoading ? (
          <Loader2 className="size-5 animate-spin" />
        ) : (
          <>
            Reenviar verificação
            <ArrowRight className="size-4" />
          </>
        )}
      </Button>
    </form>
  );
}
