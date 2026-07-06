import { type MouseEvent, useRef, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";
import {
  evaluateRules,
  type Rule,
} from "../../../utils/field-validator";

export interface SignUpProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  loginRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
  passwordRules?: Rule[];
}

export function SignUp({
  onSuccess,
  onFailed,
  loginRedirect,
  emailRules,
  passwordRules,
}: SignUpProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const emailRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const confirmPasswordRef = useRef<HTMLInputElement | null>(null);
  const { auth } = useAuth();

  const rules: Record<string, Rule[]> = {
    email: emailRules || [
      { message: "Digite um e-mail válido, ex: exemplo@dominio.com", test: v => /\S+@\S+\.\S+/.test(v) },
    ],
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

  const emailValidation = evaluateRules(rules.email, email);
  const passwordValidation = evaluateRules(rules.password, password);
  const confirmPasswordValidation = evaluateRules(rules.confirmPassword, confirmPassword);

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setSubmitted(true);

    const emailInvalid = emailValidation.some(r => !r.passed);
    const passwordInvalid = passwordValidation.some(r => !r.passed);
    const confirmPasswordInvalid = confirmPasswordValidation.some(r => !r.passed);

    if (emailInvalid) {
      emailRef.current?.focus();
      return;
    }
    if (passwordInvalid) {
      passwordRef.current?.focus();
      return;
    }
    if (confirmPasswordInvalid) {
      confirmPasswordRef.current?.focus();
      return;
    }

    setLoadingSubmit(true);

    const res = await auth.register(email, password);
    if (res.success) {
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) await onFailed(res.message, res.trace);

    setLoadingSubmit(false);
  }
  return (
    <form className="font-sans flex flex-col w-full max-w-120 min-w-60 max-h-[max(75dvh,32rem)] overflow-hidden gap-5 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] transition-transform duration-150 ease-in-out rounded-lg">
      <CardAvatar mainText="Crie sua conta" subText="Insira seus dados para começar" />
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
        <BasicInputField
          label="Senha"
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
          label="Confirme a Senha"
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
      <BasicSubmitButton label="Criar Conta" onSubmit={handleSubmit} loading={loadingSubmit} />
      {loginRedirect && <>
        <div className="flex items-center gap-2.5 w-full text-base font-semibold opacity-60">
          <hr className="flex-1" />
          OU
          <hr className="flex-1" />
        </div>
        <span className="text-sm font-semibold">
          {"Já possui uma conta? "}
          <span className="cursor-pointer text-primary transition-colors duration-200 hover:opacity-80" onClick={loginRedirect}>Entre</span>
        </span>
      </>}
    </form>
  );
}
