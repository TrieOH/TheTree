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
  onSuccess?: () => Promise<void>;
  onFailed?: (message: string) => Promise<void>;
  signUpRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
}

export function SignIn({
  onSuccess,
  onFailed,
  signUpRedirect,
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
    if(res.code === 200 && onSuccess) await onSuccess();
    else if(onFailed) await onFailed(res.message);
    setLoadingSubmit(false);
  }
  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <CardAvatar mainText="Faça login na sua conta" subText="Insira seus dados para fazer login" />
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
      <BasicSubmitButton label="Entrar" onSubmit={handleSubmit} loading={loadingSubmit}/>
      <div className="trieoh-card__divider">
        <hr />
        OU
        <hr />
      </div>
      <span className="trieoh-card__other">
        {"Ainda não possui uma conta? "}
        <span onClick={signUpRedirect}>Cadastra-se</span>
      </span>
    </form>
  );
}
