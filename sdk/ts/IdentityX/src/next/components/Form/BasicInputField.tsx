export interface BasicInputFieldProps {
  /** The Input ID/Name */
  name: string;
  /** The label text (name of the field) */
  label: string;
  /** The placeholder text (a default text to help the user) */
  placeholder: string;
  /** Input Type */
  type?: "text" | "email" | "number";
  /** Current Input Value */
  value?: string;
}
export function BasicInputField({
  name,
  label,
  placeholder,
  type = "text",
  value
}: BasicInputFieldProps) {
  return (
    <div className="trieoh trieoh-input">
      <label htmlFor={name} className="trieoh-input__label">
        {label}
      </label>
      <div className="trieoh-input__container">
        <input 
          type={type} 
          name={name} 
          id={name} 
          placeholder={placeholder}
          value={value}
          className="trieoh-input__container-field" 
        />
      </div>
    </div>
  )
}