interface BasicSubmitButtonProps {
  /** The label text (Submit Text) */
  label: string;
  /** What will happen when the user click on the button */
  onSubmit: () => void;
}

export default function BasicSubmitButton({
  label,
  onSubmit
}: BasicSubmitButtonProps) {
  return (
    <button type="submit" onClick={onSubmit} className="trieoh trieoh-button">
      {label}
    </button>
  )
}