import { useState } from "react";
import { RiEyeCloseLine, RiEyeLine } from "react-icons/ri";
import type { RuleStatus } from "../../../utils/field-validator";

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
  /** Current Input Value On Change */
  onValueChange?: (value: string) => void;
  /** OnBlur event handler */
  onBlur?: React.FocusEventHandler<HTMLInputElement>;
  /** AutoComplete */
  autoComplete?: string;
  /** Validations and their results */
  rulesStatus?: RuleStatus[];
  /** Form submission status */
  submitted?: boolean;
  /** Ref to the input element */
  inputRef?: React.Ref<HTMLInputElement>;
}

export default function BasicInputField({
  name,
  label,
  placeholder,
  type = "text",
  value,
  onValueChange,
  onBlur,
  autoComplete,
  rulesStatus = [],
  submitted = false,
  inputRef,
}: BasicInputFieldProps) {
  const [isSecretVisible, setIsSecretVisible] = useState(false);
  const hasAnyFailing = rulesStatus.some(r => !r.passed);

  return (
    <div className="font-inter relative w-full flex flex-col gap-1 text-trieoh-neutral2">
      <label htmlFor={name} className="text-[1rem] font-semibold">
        {label}
      </label>
      <div 
        className={`flex justify-between items-center px-[0.625rem] py-[0.0625rem] gap-[0.625rem] border-b-2 border-trieoh-neutral2 ${
          (hasAnyFailing && submitted) ? "!border-[#e53935]" : ""
        }`}
      >
        <input 
          type={isSecretVisible ? "text" : type} 
          name={name} 
          id={name} 
          placeholder={placeholder}
          value={value}
          onChange={(e) => onValueChange && onValueChange(e.target.value)}
          onBlur={onBlur}
          autoComplete={autoComplete}
          aria-invalid={hasAnyFailing && submitted}
          ref={inputRef}
          className="min-w-[10rem] flex-1 text-trieoh-base font-light text-trieoh-neutral2 appearance-none bg-transparent outline-none border-none !shadow-none py-[0.125rem]" 
        />
        {type === "password" && (
          isSecretVisible ?
            <RiEyeCloseLine 
              className="cursor-pointer shrink-0 select-none"
              size={24}
              onClick={() => setIsSecretVisible(false)} 
            />
          :
            <RiEyeLine 
              className="cursor-pointer shrink-0 select-none"
              size={24}
              onClick={() => setIsSecretVisible(true)} 
            />
          )
        }
      </div>

      <div className="text-[0.75rem] text-[#6b7280] transition-opacity duration-200 ease-in-out">
        {rulesStatus.map((r, i) => {
          const classes = [
            "transition-[color,text-decoration,opacity] duration-[120ms] ease opacity-95 m-[0.125rem]",
            r.passed ? "line-through opacity-60 text-[#10b981]" : "",
            !r.passed && submitted ? "text-[#e53935] font-semibold opacity-100" : "",
          ]
            .filter(Boolean)
            .join(" ");
          return (
            <p key={r.id ?? i} className={classes}>
              {r.message}
            </p>
          );
        })}
      </div>

    </div>
  )
}