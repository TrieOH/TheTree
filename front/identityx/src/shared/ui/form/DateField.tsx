import { useRef, useState } from "react";
import { useFieldContext } from "@/shared/lib/forms";
import { CalendarIcon } from "lucide-react";
import type { RuleStatus } from "./types";

interface PropsI {
  label: string;
  placeholder?: string;
  required?: boolean;
  getRulesStatus?: (value: unknown) => RuleStatus[]
  submitted?: boolean;
}

const toDateInputValue = (iso?: string) => {
  if (!iso) return '';
  return iso.split('T')[0];
};

export default function DateField({ label, placeholder, required, getRulesStatus, submitted }: PropsI) {
  const field = useFieldContext<string>();
  const inputRef = useRef<HTMLInputElement>(null);
  const [focused, setFocused] = useState(false);

  const rulesStatus = getRulesStatus ? getRulesStatus(field.state.value) : [];
  const hasAnyFailing = rulesStatus.some(r => !r.passed);

  const openDatePicker = () => {
    inputRef.current?.showPicker();
  };

  return (
    <div className="font-sans relative w-full flex flex-col gap-1 text-foreground mt-2">
      <label htmlFor={field.name} className="text-base font-semibold">
        {required ? `${label} *` : label}
      </label>
      <div
        className={`flex justify-between items-center px-2.5 py-px gap-2.5 border-b-2 border-foreground ${(hasAnyFailing && submitted) ? "border-destructive!" : focused ? "border-primary" : ""}`}
      >
        <input
          ref={inputRef}
          type="date"
          id={field.name}
          name={field.name}
          value={toDateInputValue(field.state.value)}
          onChange={(e) => {
            const rawValue = e.target.value;
            if (!rawValue) {
              field.handleChange('');
              return;
            }
            field.handleChange(`${rawValue}T23:59:59Z`);
          }}
          onBlur={() => { setFocused(false); field.handleBlur(); }}
          onFocus={() => setFocused(true)}
          placeholder={placeholder}
          aria-invalid={hasAnyFailing && submitted}
          className="min-w-40 flex-1 text-base font-light text-foreground appearance-none bg-transparent outline-none border-none shadow-none! py-0.5 dark:scheme-dark"
        />
        <button
          type="button"
          onClick={openDatePicker}
          tabIndex={-1}
          className="p-0 m-0 bg-transparent border-none outline-none cursor-pointer shrink-0"
        >
          <CalendarIcon className="size-5 text-muted-foreground" />
        </button>
      </div>

      <div className="text-xs text-muted-foreground transition-opacity duration-200 ease-in-out">
        {rulesStatus.map((r) => {
          const classes = [
            "transition-[color,text-decoration,opacity] duration-[120ms] ease opacity-95 m-[0.125rem]",
            r.passed ? "line-through opacity-60 text-green-500" : "",
            !r.passed && submitted ? "text-destructive font-semibold opacity-100" : "",
          ]
            .filter(Boolean)
            .join(" ");
          return (
            <p key={r.message} className={classes}>
              {r.message}
            </p>
          );
        })}
      </div>
    </div>
  )
}
