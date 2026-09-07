import { createSignal } from "solid-js";
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
}
export function FormInput(props: FormInputProps) {
  const [visible, setVisible] = createSignal(false);
  const password = () => props.type === "password";
  return (
    <div class="w-full font-body group relative">
      <input
        id={props.id ?? props.name}
        name={props.name}
        type={password() && !visible() ? "password" : (props.type ?? "text")}
        value={props.value ?? ""}
        placeholder=" "
        autocomplete={props.autocomplete}
        onInput={props.onInput}
        class={`peer h-13 w-full rounded-sm text-sm text-foreground outline-none border border-secondary/20 bg-muted/40 transition-all px-3.5 pt-4 pb-1 ${props.error ? "bg-destructive/5 border-destructive/50" : ""}`}
      />
      <label
        for={props.id ?? props.name}
        class="pointer-events-none absolute transition-all left-3.5 top-1.5 text-[10px] font-semibold tracking-wider text-primary/70 peer-placeholder-shown:top-3.5 peer-placeholder-shown:text-[15px] peer-placeholder-shown:font-normal"
      >
        {props.label}
      </label>
      {password() && (
        <button
          type="button"
          tabindex={-1}
          onClick={() => {
            setVisible(!visible());
          }}
          class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground"
        >
          {visible() ? "Ocultar" : "Mostrar"}
        </button>
      )}
    </div>
  );
}
