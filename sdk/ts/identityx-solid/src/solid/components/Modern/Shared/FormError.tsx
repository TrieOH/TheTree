export interface FormErrorProps {
  message?: string;
  class?: string;
}
export function FormError(props: FormErrorProps) {
  return props.message ? (
    <p
      class={`text-[11px] font-medium text-destructive tracking-tight leading-none px-1 mt-1.5 ${props.class ?? ""}`}
    >
      {props.message}
    </p>
  ) : null;
}
