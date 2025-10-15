interface BasicSubmitButtonProps {
  label: string;
}

export default function BasicSubmitButton({
  label
}: BasicSubmitButtonProps) {
  return (
    <button type="submit" className="trieoh trieoh-button">
      {label}
    </button>
  )
}