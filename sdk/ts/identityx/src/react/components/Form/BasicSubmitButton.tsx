import type { MouseEvent } from "react";

interface BasicSubmitButtonProps {
  /** The label text (Submit Text) */
  label: string;
  /** What will happen when the user click on the button */
  onSubmit: (e: MouseEvent<HTMLButtonElement>) => void;
  /** Is performing the submit */
  loading: boolean;
}

export default function BasicSubmitButton({
  label,
  onSubmit,
  loading
}: BasicSubmitButtonProps) {
  return (
    <button
      type="submit"
      onClick={onSubmit}
      disabled={loading}
      className={`font-sans w-full h-13 text-xl font-semibold outline-none bg-transparent relative overflow-hidden min-w-40 shrink-0 border-2 border-foreground text-foreground cursor-pointer px-6 transition-transform duration-500 rounded-lg hover:scale-[1.02] active:scale-[0.99] disabled:opacity-60 disabled:cursor-not-allowed disabled:transform-none! ${loading ? "button-loading" : ""
        }`}
    >
      {label}
    </button>
  )
}