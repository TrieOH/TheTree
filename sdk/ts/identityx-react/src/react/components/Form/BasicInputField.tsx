import { useState } from "react";
import { Eye, EyeOff } from "lucide-react";
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
    <div className="font-sans relative w-full flex flex-col gap-1 text-foreground">
      <label htmlFor={name} className="text-base font-semibold">
        {label}
      </label>
      <div
        className={`flex justify-between items-center px-2.5 py-px gap-2.5 border-b-2 border-foreground ${(hasAnyFailing && submitted) ? "border-destructive!" : ""
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
          className="min-w-40 flex-1 text-base font-light text-foreground appearance-none bg-transparent outline-none border-none shadow-none! py-0.5"
        />
        {type === "password" && (
          isSecretVisible ?
            <EyeOff
              className="cursor-pointer shrink-0 select-none"
              size={24}
              onClick={() => setIsSecretVisible(false)}
            />
            :
            <Eye
              className="cursor-pointer shrink-0 select-none"
              size={24}
              onClick={() => setIsSecretVisible(true)}
            />
        )
        }
      </div>

      <div className="text-xs text-muted-foreground transition-opacity duration-200 ease-in-out">
        {rulesStatus.map((r, i) => {
          const classes = [
            "transition-[color,text-decoration,opacity] duration-[120ms] ease opacity-95 m-[0.125rem]",
            r.passed ? "line-through opacity-60 text-green-500" : "",
            !r.passed && submitted ? "text-destructive font-semibold opacity-100" : "",
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