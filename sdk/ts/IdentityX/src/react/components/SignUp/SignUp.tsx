import { type MouseEvent, useRef, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";
import CardAvatar from "../Form/CardAvatar";
import DynamicFields from "../Form/DynamicFields";
import { 
  evaluateRules, 
  type Rule,
} from "../../../utils/field-validator";
import type { FieldDefinitionResultI, FieldValue } from "../../../types/fields-types";

export interface SignUpProps {
  onSuccess?: () => Promise<void>;
  onFailed?: (message: string, trace?: string[]) => Promise<void>;
  loginRedirect?:(e: MouseEvent<HTMLSpanElement>) => void;
  emailRules?: Rule[];
  passwordRules?: Rule[];
  flow_id?: string;
  fields?: FieldDefinitionResultI[];
}

export function SignUp({
  onSuccess,
  onFailed,
  loginRedirect,
  emailRules,
  passwordRules,
  flow_id,
  fields = [],
}: SignUpProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [dynamicValues, setDynamicValues] = useState<Record<string, FieldValue>>({});
  const [submitted, setSubmitted] = useState(false);
  const [loadingSubmit, setLoadingSubmit] = useState(false);
  const emailRef = useRef<HTMLInputElement | null>(null);
  const passwordRef = useRef<HTMLInputElement | null>(null);
  const confirmPasswordRef = useRef<HTMLInputElement | null>(null);
  const { auth } = useAuth();

  const handleDynamicChange = (key: string, value: FieldValue) => {
    setDynamicValues(prev => ({ ...prev, [key]: value }));
  };

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
    
    // Simple validation for dynamic required fields
    const dynamicInvalid = fields.some(f => f.required && !dynamicValues[f.key]);

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
    if (dynamicInvalid) return;
    
    setLoadingSubmit(true);

    const res = await auth.register(email, password, flow_id, dynamicValues);
    if(res.code === 201 && onSuccess) await onSuccess();
    else if(onFailed) await onFailed(res.message, res.trace);
    setLoadingSubmit(false);
  }
  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <CardAvatar mainText="Crie sua conta" subText="Insira seus dados para começar" />
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
        <DynamicFields 
          fields={fields} 
          values={dynamicValues} 
          onValueChange={handleDynamicChange}
          submitted={submitted}
        />
      </div>
      <BasicSubmitButton label="Criar Conta" onSubmit={handleSubmit} loading={loadingSubmit}/>
      {loginRedirect && <>
        <div className="trieoh-card__divider">
          <hr />
          OU
          <hr />
        </div>
        <span className="trieoh-card__other">
          {"Já possui uma conta? "}
          <span onClick={loginRedirect}>Entre</span>
        </span>
      </>}
    </form>
  );
}
