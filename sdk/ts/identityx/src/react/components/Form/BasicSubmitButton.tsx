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
      className={`font-inter w-full h-[3.25rem] text-trieoh-xl font-semibold outline-none bg-transparent relative overflow-hidden min-w-[10rem] shrink-0 border-2 border-trieoh-neutral2 text-trieoh-neutral2 cursor-pointer px-[1.5rem] transition-transform duration-500 rounded-[0.25rem] hover:scale-[1.02] active:scale-[0.99] disabled:opacity-60 disabled:cursor-not-allowed disabled:!transform-none ${
        loading ? "trieoh-button-loading" : ""
      }`}
    >
      {label}
    </button>
  )
}