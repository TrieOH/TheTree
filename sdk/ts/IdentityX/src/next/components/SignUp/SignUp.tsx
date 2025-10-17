import { type MouseEvent, useRef, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import { 
  evaluateRules, 
  type Rule,
} from "../../../utils/field-validator";

export interface SignUpProps {
  onSuccess?: () => void;
  onFailed?: (message: string) => void;
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
  passwordRules?: Rule[];
}

export function SignUp({
  onSuccess,
  onFailed,
  loginRedirect,
  emailRules,
  passwordRules
}: SignUpProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const emailRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const { auth } = useAuth();

  const rules: Record<string, Rule[]> = {
    email: emailRules || [
      { message: "Digite um e-mail válido, ex: exemplo@dominio.com", test: v => /\S+@\S+\.\S+/.test(v) },
    ],
    password: passwordRules || [
      { message: "Mínimo de 6 caracteres.", test: v => v.length >= 6 },
      { message: "Deve conter uma letra maiúscula.", test: v => /[A-Z]/.test(v) },
      { 
        message: "Inclua pelo menos um caractere especial, ex: ! @ # $ % & * . ,", 
        test: v => /[!@#$%^&*(),.?":{}|<>_\-+=~`;/\\[\]]/.test(v) 
      },
      { message: "Deve conter um número.", test: v => /\d/.test(v) },
    ],
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

    const res = await auth.register(email, password);
    console.log(res)
    if(res.code === 201 && onSuccess) onSuccess();
    else if(onFailed) onFailed(res.message);
  }
  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <h3 className="trieoh-card__title">Crie uma Conta!</h3>
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
      <BasicSubmitButton label="Criar Conta" onSubmit={handleSubmit}/>
      <div className="trieoh-card__divider">
        <hr />
        OU
        <hr />
      </div>
      <span className="trieoh-card__other">
        {"Ainda não possui uma conta? "}
        <span onClick={loginRedirect}>Entre</span>
      </span>
    </form>
  );
}
