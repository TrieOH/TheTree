import { type MouseEvent, useRef, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";
import {
  evaluateRules,
  type Rule,
} from "../../../utils/field-validator";

export interface SignInProps {
  onSuccess?: (message?: string) => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  signUpRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  forgotPasswordRedirect?: (e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
}

export function SignIn({
  onSuccess,
  onFailed,
  signUpRedirect,
  forgotPasswordRedirect,
  emailRules,
}: SignInProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const emailRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const { auth } = useAuth();

  const rules: Record<string, Rule[]> = {
    email: emailRules || [
      { message: "Digite um e-mail válido, ex: exemplo@dominio.com", test: v => /\S+@\S+\.\S+/.test(v) },
    ],
    password: [],
  };

  const emailValidation = evaluateRules(rules.email, email);
  const passwordValidation = evaluateRules(rules.password, password);

  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    setSubmitted(true);

    const emailInvalid = emailValidation.some(r => !r.passed);
    const passwordInvalid = passwordValidation.some(r => !r.passed);

    if (emailInvalid) {
      emailRef.current?.focus();
      return;
    }
    if (passwordInvalid) {
      passwordRef.current?.focus();
      return;
    }

    setLoadingSubmit(true);

    const res = await auth.login(email, password);
    if (res.success) {
      if (onSuccess) await onSuccess(res.message);
    } else if (onFailed) await onFailed(res.message, res.trace);

    setLoadingSubmit(false);
  }
  return (
    <form className="font-sans flex flex-col w-full max-w-120 min-w-60 max-h-[max(75dvh,32rem)] overflow-hidden gap-5 items-center bg-background text-foreground p-[1.25rem_1.5rem] shadow-[0_0.25rem_0.25rem_0_rgba(0,0,0,0.25)] transition-transform duration-150 ease-in-out rounded-lg">
      <CardAvatar mainText="Faça login na sua conta" subText="Insira seus dados para fazer login" />
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
          autoComplete="current-password"
          type="password"
          value={password}
          onValueChange={setPassword}
          inputRef={passwordRef}
          rulesStatus={passwordValidation}
          submitted={submitted}
        />
      </div>
      <BasicSubmitButton label="Entrar" onSubmit={handleSubmit} loading={loadingSubmit} />
      {forgotPasswordRedirect && (
        <span className="text-sm font-semibold text-center cursor-pointer">
          <span className="cursor-pointer text-primary transition-colors duration-200 hover:opacity-80" onClick={forgotPasswordRedirect}>
            Esqueceu sua senha?
          </span>
        </span>
      )}
      {signUpRedirect && <>
        <div className="flex items-center gap-2.5 w-full text-base font-semibold opacity-60">
          <hr className="flex-1" />
          OU
          <hr className="flex-1" />
        </div>
        <span className="text-sm font-semibold">
          {"Ainda não possui uma conta? "}
          <span className="cursor-pointer text-primary transition-colors duration-200 hover:opacity-80" onClick={signUpRedirect}>Cadastra-se</span>
        </span>
      </>}
    </form>
  );
}
