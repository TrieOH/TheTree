import { type MouseEvent, useRef, useState } from "react";
import { evaluateRules, type Rule } from "../../../utils/field-validator";
import CardAvatar from "../Form/CardAvatar";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";

export interface ForgotPasswordProps {
  onSuccess?: () => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
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
    if(res.code === 201 && onSuccess) await onSuccess();
    else if(onFailed) await onFailed(res.message, res.trace);
    setLoadingSubmit(false);
  }

  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <CardAvatar 
        mainText="Esqueci a senha" 
        subText="Insira seu e-mail para receber um link de redefinição." 
      />

      <div className="trieoh-card__fields">
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
        <div className="trieoh-card__divider">
          <hr />
          OU
          <hr />
        </div>

        <span className="trieoh-card__other">
          {"Lembrou-se da sua senha? "}
          <span onClick={loginRedirect}>
            Login
          </span>
        </span>
      </>}
    </form>
  );
}