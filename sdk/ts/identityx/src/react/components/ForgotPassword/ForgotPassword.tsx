import { type MouseEvent, useRef, useState } from "react";
import { evaluateRules, type Rule } from "../../../utils/field-validator";
import CardAvatar from "../Form/CardAvatar";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";

export interface ForgotPasswordProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
}

export default function ForgotPassword({
  onSuccess,
  onFailed,
  loginRedirect,
  emailRules,
}: ForgotPasswordProps) {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const emailRef = useRef<HTMLInputElement | null>(null);

  const { auth } = useAuth();

  const rules: Record<string, Rule[]> = {
    email: emailRules || [
      { message: "Digite um e-mail válido, ex: exemplo@dominio.com", test: v => /\S+@\S+\.\S+/.test(v) },
    ],
  };

  const emailValidation = evaluateRules(rules.email, email);

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setSubmitted(true);

    const emailInvalid = emailValidation.some(r => !r.passed);

    if (emailInvalid) {
      emailRef.current?.focus();
      return;
    }

    setLoadingSubmit(true);

    const res = await auth.sendForgotPassword(email);

    if (res.success) {
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) await onFailed(res.message, res.trace);

    setLoadingSubmit(false);
  }

  return (
    <form className="font-sans flex flex-col w-full max-w-120 min-w-60 max-h-[max(75dvh,32rem)] overflow-hidden gap-5 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] transition-transform duration-150 ease-in-out rounded-lg">
      <CardAvatar
        mainText="Esqueci a senha"
        subText="Insira seu e-mail para receber um link de redefinição."
      />

      <div className="w-full flex flex-col gap-2.5 flex-[1_1_auto] overflow-y-auto mb-2">
        <BasicInputField
          label="Email"
          name="email"
          placeholder="teste@gmail.com"
          autoComplete="email"
          type="email"
          value={email}
          onValueChange={setEmail}
          inputRef={emailRef}
          rulesStatus={emailValidation}
          submitted={submitted}
        />
      </div>

      <BasicSubmitButton
        label={loadingSubmit ? "Enviando..." : "Enviar link de redefinição"}
        onSubmit={handleSubmit} loading={loadingSubmit}
      />
      {loginRedirect && <>
        <div className="flex items-center gap-2.5 w-full text-base font-semibold opacity-60">
          <hr className="flex-1" />
          OU
          <hr className="flex-1" />
        </div>

        <span className="text-sm font-semibold">
          {"Lembrou-se da sua senha? "}
          <span className="cursor-pointer text-primary transition-colors duration-200 hover:opacity-80" onClick={loginRedirect}>
            Login
          </span>
        </span>
      </>}
    </form>
  );
}