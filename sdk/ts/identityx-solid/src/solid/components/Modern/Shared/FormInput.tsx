import { createSignal, createUniqueId, omit } from "solid-js";
import { Eye, EyeOff } from "./Icons";

export interface FormInputProps {
  label: string;
  type?: string;
  value?: string;
  id?: string;
  name?: string;
  placeholder?: string;
  error?: boolean;
  autocomplete?: string;
  onInput?: (event: InputEvent & { currentTarget: HTMLInputElement }) => void;
  class?: string;
}

export function FormInput(props: FormInputProps) {
  const rest = omit(props, "label", "type", "value", "id", "name", "error", "autocomplete", "onInput", "class");
  const [visible, setVisible] = createSignal(false);
  const generatedId = createUniqueId();
  const id = () => props.id ?? props.name ?? generatedId;
  const isPassword = () => props.type === "password";

  return (
    <div class={`w-full font-body group relative ${props.class ?? ""}`}>
      <input
        id={id()}
        name={props.name ?? id()}
        aria-invalid={props.error ? "true" : undefined}
        aria-describedby={props.error ? `${id()}-error` : undefined}
        type={isPassword() ? (visible() ? "text" : "password") : (props.type ?? "text")}
        value={props.value ?? ""}
        autocomplete={props.autocomplete ?? (isPassword() ? "new-password" : undefined)}
        placeholder=" "
        onInput={props.onInput}
        class={`peer h-13 w-full rounded-sm text-sm text-foreground selection:bg-primary/20 selection:text-foreground outline-none border border-secondary/20 bg-muted/40 transition-all duration-200 placeholder:text-transparent focus-visible:ring-0 focus:bg-muted/60 px-3.5 pt-4 pb-1 ${isPassword() ? "pr-10" : ""} ${props.error ? "ring-1 ring-destructive/20 border-destructive/50" : ""}`}
        style={props.error ? { "background-color": "color-mix(in oklab, var(--color-destructive) 5%, transparent)" } : undefined}
        {...rest}
      />
      <label
        for={id()}
        class={`pointer-events-none absolute transition-all duration-200 ease-out left-3.5 top-1.5 text-[10px] font-semibold tracking-wider text-primary/70 peer-placeholder-shown:top-3.5 peer-placeholder-shown:text-[15px] peer-placeholder-shown:font-normal peer-placeholder-shown:tracking-normal peer-placeholder-shown:text-muted-foreground peer-focus:top-1.5 peer-focus:text-[10px] peer-focus:font-semibold peer-focus:tracking-wider peer-focus:text-primary/70 ${props.error ? "text-destructive/80 peer-placeholder-shown:text-destructive/60 peer-focus:text-destructive/80" : ""}`}
      >
        {props.label}
      </label>
      {isPassword() && (
        <button
          type="button"
          tabindex={-1}
          onClick={() => setVisible(!visible())}
          class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground transition-colors hover:text-primary/70 focus:outline-none"
          aria-label={visible() ? "Ocultar senha" : "Mostrar senha"}
        >
          {visible() ? <EyeOff size={16} /> : <Eye size={16} />}
        </button>
      )}
    </div>
  );
}
