import { type MouseEvent, useRef, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";
import {
  evaluateRules,
  type Rule,
} from "../../../utils/field-validator";

export interface ResetPasswordProps {
  token: string;
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  passwordRules?: Rule[];
}

export function ResetPassword({
  token,
  onSuccess,
  onFailed,
  loginRedirect,
  passwordRules,
}: ResetPasswordProps) {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const confirmPasswordRef = useRef<HTMLInputElement | null>(null);
  const { auth } = useAuth();

  const rules: Record<string, Rule[]> = {
    password: passwordRules || [
      { message: "Mínimo de 8 caracteres.", test: v => v.length >= 8 },
      { message: "Deve conter uma letra maiúscula.", test: v => /[A-Z]/.test(v) },
      {
        message: "Inclua pelo menos um caractere especial, ex: ! @ # $ % & * . ,",
        test: v => /[!@#$%^&*(),.?":{}|<>_\-+=~`;/\\[\]]/.test(v)
      },
      { message: "Deve conter um número.", test: v => /\d/.test(v) },
    ],
    confirmPassword: [
      { message: "Senhas não conferem.", test: v => v === password },
    ],
  };

  const passwordValidation = evaluateRules(rules.password, password);
  const confirmPasswordValidation = evaluateRules(rules.confirmPassword, confirmPassword);

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setSubmitted(true);

    const passwordInvalid = passwordValidation.some(r => !r.passed);
    const confirmPasswordInvalid = confirmPasswordValidation.some(r => !r.passed);

    if (passwordInvalid) {
      passwordRef.current?.focus();
      return;
    }
    if (confirmPasswordInvalid) {
      confirmPasswordRef.current?.focus();
      return;
    }

    setLoadingSubmit(true);

    const res = await auth.resetPassword(token, password);
    if (res.success) {
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) await onFailed(res.message, res.trace);

    setLoadingSubmit(false);
  }

  return (
    <form className="font-sans flex flex-col w-full max-w-120 min-w-60 max-h-[max(75dvh,32rem)] overflow-hidden gap-5 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] transition-transform duration-150 ease-in-out rounded-lg">
      <CardAvatar mainText="Redefinir Senha" subText="Crie uma nova senha para sua conta" />
      <div className="w-full flex flex-col gap-2.5 flex-[1_1_auto] overflow-y-auto mb-2">
        <BasicInputField
          label="Nova Senha"
          name="password"
          placeholder="**********"
          autoComplete="new-password"
          type="password"
          value={password}
          onValueChange={setPassword}
          inputRef={passwordRef}
          rulesStatus={passwordValidation}
          submitted={submitted}
        />
        <BasicInputField
          label="Confirme a Nova Senha"
          name="confirm-password"
          placeholder="**********"
          autoComplete="new-password"
          type="password"
          value={confirmPassword}
          onValueChange={setConfirmPassword}
          inputRef={confirmPasswordRef}
          rulesStatus={confirmPasswordValidation}
          submitted={submitted}
        />
      </div>
      <BasicSubmitButton label="Redefinir Senha" onSubmit={handleSubmit} loading={loadingSubmit} />
      {loginRedirect && (
        <span className="text-sm font-semibold">
          {"Lembrou-se da sua senha? "}
          <span className="cursor-pointer text-primary transition-colors duration-200 hover:opacity-80" onClick={loginRedirect}>
            Login
          </span>
        </span>
      )}
    </form>
  );
}
