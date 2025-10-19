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
      className={`trieoh trieoh-button ${loading ? "trieoh-button--loading" : ""}`}
    >
      {label}
    </button>
  )
}