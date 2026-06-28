import { AuthProvider } from "../../react";
import BasicSubmitButton from "../../react/components/Form/BasicSubmitButton";
import CardAvatar from "../../react/components/Form/CardAvatar";
import { Check, Info } from "lucide-react";

interface VerifyEmailMockProps {
  status: "verifying" | "success" | "error" | "already_verified";
  message: string;
}

/**
 * A visual-only version of VerifyEmail for Storybook.
 * It doesn't trigger any real API calls.
 */
function VerifyEmailVisualMock({ status, message }: VerifyEmailMockProps) {
  const mainText = status === "verifying" ? "Verificando..." :
    status === "error" ? "Ops! Falha na verificação" : "Tudo pronto!";

  return (
    <div className="font-sans flex flex-col w-full h-full min-h-screen items-center justify-center bg-background text-foreground p-6">
      <div className="w-full max-w-120 flex flex-col gap-8 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] rounded-lg">
        <CardAvatar
          mainText={mainText}
          subText={message}
        />

        {status === "verifying" && (
          <div className="w-full flex justify-center py-4">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
          </div>
        )}

        {status === "error" && (
          <div className="w-full flex flex-col gap-6 items-center">
            <div className="w-full mt-2">
              <BasicSubmitButton
                label="Tentar novamente"
                onSubmit={() => { }}
                loading={false}
              />
            </div>
          </div>
        )}

        {(status === "success" || status === "already_verified") && (
          <div className="w-full flex flex-col gap-4 items-center">
            <div className="w-16 h-16 flex items-center justify-center rounded-full bg-primary/10 text-primary">
              {status === "success" ? <Check size={40} /> : <Info size={40} />}
            </div>
            <p className="text-sm font-semibold opacity-50 text-center">
              Você já pode fechar esta janela ou continuar navegando em sua conta.
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

interface VerifyEmailWithProviderProps {
  token: string;
  isProjectMode?: boolean;
  mockState?: "verifying" | "success" | "error" | "already_verified";
}

export default function VerifyEmailWithProvider({
  token,
  isProjectMode = true,
  mockState
}: VerifyEmailWithProviderProps) {
  if (mockState) {
    const messages = {
      verifying: "Verificando seu e-mail...",
      success: "E-mail verificado com sucesso!",
      error: "Falha na verificação do e-mail.",
      already_verified: "Seu e-mail já está verificado."
    };
    return <VerifyEmailVisualMock status={mockState} message={messages[mockState]} />;
  }

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <AuthProvider baseURL="http://localhost:8080" isProjectMode={isProjectMode}>
        <VerifyEmailWithRealLogic token={token} />
      </AuthProvider>
    </div>
  )
}

import { VerifyEmail as VerifyEmailWithRealLogic } from "../../react";
