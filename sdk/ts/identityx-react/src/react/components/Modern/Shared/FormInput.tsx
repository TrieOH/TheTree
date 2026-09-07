import { useState, useId, forwardRef } from "react";
import { Eye, EyeOff } from "lucide-react";
import { cn } from "../../../../utils/cn";

interface PropsI extends React.InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: boolean;
}

const FormInput = forwardRef<HTMLInputElement, PropsI>(({
  label,
  id: propId,
  type = "text",
  error,
  className,
  value,
  onChange,
  onFocus,
  onBlur,
  ...props
}, ref) => {
  const [showPassword, setShowPassword] = useState(false);
  const generatedId = useId();
  const id = propId ?? generatedId;

  const isPassword = type === "password";
  const inputType = isPassword ? (showPassword ? "text" : "password") : type;

  return (
    <div className="w-full font-body group relative">
      <input
        ref={ref}
        id={id}
        type={inputType}
        value={value}
        placeholder=" "
        onChange={onChange}
        onFocus={onFocus}
        onBlur={onBlur}
        className={cn(
          "peer h-13 w-full rounded-sm text-sm text-foreground outline-none",
          "border border-secondary/20 bg-muted/40 transition-all duration-200",
          "placeholder:text-transparent focus-visible:ring-0 focus:bg-muted/60",
          "px-3.5 pt-4 pb-1",
          isPassword && "pr-10",
          error && "bg-destructive/5 ring-1 ring-destructive/20 border-destructive/50",
          className
        )}
        {...props}
      />

      <label
        htmlFor={id}
        className={cn(
          "pointer-events-none absolute transition-all duration-200 ease-out left-3.5",

          // Float state (Default)
          "top-1.5 text-[10px] font-semibold tracking-wider text-primary/70",

          // Default (Float 'disabled')
          "peer-placeholder-shown:top-3.5 peer-placeholder-shown:text-[15px] peer-placeholder-shown:font-normal peer-placeholder-shown:tracking-normal peer-placeholder-shown:text-muted-foreground",

          // Float state on focus
          "peer-focus:top-1.5 peer-focus:text-[10px] peer-focus:font-semibold peer-focus:tracking-wider peer-focus:text-primary/70",

          // Float error
          error && "text-destructive/80",
          error && "peer-placeholder-shown:text-destructive/60",
          error && "peer-focus:text-destructive/80"
        )}
      >
        {label}
      </label>

      {isPassword && (
        <button
          type="button"
          tabIndex={-1}
          onClick={() => { setShowPassword((prev) => !prev); }}
          className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground transition-colors hover:text-primary/70 focus:outline-none"
          aria-label={showPassword ? "Ocultar senha" : "Mostrar senha"}
        >
          {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
        </button>
      )}
    </div>
  );
});

FormInput.displayName = "FormInput";

export default FormInput;
