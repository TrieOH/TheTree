import { type MouseEvent, useState } from "react";
import { useAuth } from "../../AuthProvider";
import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";

export function SignIn() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { auth } = useAuth();
  const handleSubmit = async (e: MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    const res = await auth.login(email, password);
    console.log(res)
  }
  return (
    <form className="trieoh trieoh-card trieoh-card--full-rounded">
      <h3 className="trieoh-card__title">Faça seu Login!</h3>
      <div className="trieoh-card__fields">
        <BasicInputField 
          label="Email" 
          name="email"
          placeholder="teste@gmail.com"
          autoComplete="email"
          type="email"
          value={email}
          onValueChange={setEmail}
        />
        <BasicInputField 
          label="Senha" 
          name="password"
          placeholder="**********"
          autoComplete="current-password"
          type="password"
          value={password}
          onValueChange={setPassword}
        />
      </div>
      <BasicSubmitButton label="Entrar" onSubmit={handleSubmit}/>
      <div className="trieoh-card__divider">
        <hr />
        OU
        <hr />
      </div>
      <span className="trieoh-card__other">Ainda não possui uma conta? <span>Cadastra-se</span></span>
    </form>
  );
}
