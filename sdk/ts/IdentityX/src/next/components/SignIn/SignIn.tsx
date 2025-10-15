import BasicInputField from "../Form/BasicInputField";
import BasicSubmitButton from "../Form/BasicSubmitButton";

export interface SignInProps {
  /** Click Handler - Perform Login */
  onSubmit: () => void;
}

export function SignIn({

}: SignInProps) {
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
        />
        <BasicInputField 
          label="Senha" 
          name="password"
          placeholder="**********"
          autoComplete="current-password"
          type="password"
        />
      </div>
      <BasicSubmitButton label="Entrar"/>
      <div className="trieoh-card__divider">
        <hr />
        OU
        <hr />
      </div>
      <span className="trieoh-card__other">Ainda não possui uma conta? <span>Cadastra-se</span></span>
    </form>
  );
}
