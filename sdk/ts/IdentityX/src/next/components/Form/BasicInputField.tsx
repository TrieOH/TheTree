import { Activity, useState } from "react";
import EyeIcon from "../Icons/EyeIcon";
import EyeClosedIcon from "../Icons/EyeClosedIcon";

interface BasicInputFieldProps {
  /** The Input ID/Name */
  name: string;
  /** The label text (name of the field) */
  label: string;
  /** The placeholder text (a default text to help the user) */
  placeholder: string;
  /** Input Type */
  type?: "text" | "email" | "number" | "password";
  /** Current Input Value */
  value?: string;
  /** Hint/AutoComplete */
  autoComplete?: string;
}
export default function BasicInputField({
  name,
  label,
  placeholder,
  type = "text",
  value,
  autoComplete
}: BasicInputFieldProps) {
  const [isSecretVisible, setIsSecretVisible] = useState(false);
  return (
    <div className="trieoh trieoh-input">
      <label htmlFor={name} className="trieoh-input__label">
        {label}
      </label>
      <div className="trieoh-input__container">
        <input 
          type={isSecretVisible ? "text" : type} 
          name={name} 
          id={name} 
          placeholder={placeholder}
          value={value}
          autoComplete={autoComplete}
          className="trieoh-input__container-field" 
        />
        <Activity mode={type === "password" ? "visible" : "hidden"}>
          { isSecretVisible ?
            <EyeClosedIcon 
              className="trieoh-input__container-icon" 
              onClick={() => setIsSecretVisible(false)} 
            />
            :
            <EyeIcon 
              className="trieoh-input__container-icon"
              onClick={() => setIsSecretVisible(true)} 
            />
          }
        </Activity>
      </div>
    </div>
  )
}