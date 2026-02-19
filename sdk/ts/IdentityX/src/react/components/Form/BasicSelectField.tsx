import { useState, useRef, useEffect } from "react";
import { RiArrowDownSLine } from "react-icons/ri";
import type { RuleStatus } from "../../../utils/field-validator";

interface Option {
  id: string | number;
  value: string;
  label: string;
}

interface CustomSelectProps {
  /** The Input ID/Name */
  name: string;
  /** The label text (name of the field) */
  label: string;
  /** The placeholder text (a default text to help the user) */
  placeholder?: string;
  /** Current selected value */
  value?: string;
  /** Available options */
  options: Option[];
  /** Current Value On Change */
  onValueChange?: (value: string) => void;
  /** OnBlur event handler */
  onBlur?: React.FocusEventHandler<HTMLDivElement>;
  /** Validations and their results */
  rulesStatus?: RuleStatus[];
  /** Form submission status */
  submitted?: boolean;
  /** Ref to the trigger element */
  triggerRef?: React.Ref<HTMLDivElement>;
  /** Disabled state */
  disabled?: boolean;
}

export default function BasicSelectField({
  name,
  label,
  placeholder = "Selecione uma opção",
  value,
  options,
  onValueChange,
  onBlur,
  rulesStatus = [],
  submitted = false,
  triggerRef,
  disabled = false,
}: CustomSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedLabel, setSelectedLabel] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);
  const hasAnyFailing = rulesStatus.some((r) => !r.passed);

  useEffect(() => {
    const selected = options.find((opt) => opt.value === value);
    setSelectedLabel(selected ? selected.label : "");
  }, [value, options]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) setIsOpen(false);
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSelect = (optionValue: string) => {
    onValueChange && onValueChange(optionValue);
    setIsOpen(false);
  };

  const handleToggle = () => {
    if (!disabled) setIsOpen(!isOpen);
  };

  const handleBlur: React.FocusEventHandler<HTMLDivElement> = (e) => {
    if (!containerRef.current?.contains(e.relatedTarget)) onBlur && onBlur(e);
  };

  const displayValue = selectedLabel || placeholder;

  return (
    <div className="trieoh trieoh-input">
      <label htmlFor={name} className="trieoh-input__label">
        {label}
      </label>
      
      <div className="trieoh-custom-select__wrapper" ref={containerRef}>
        <div
          className={
            (hasAnyFailing && submitted ? "trieoh-input__container--error " : "") +
            "trieoh-input__container trieoh-custom-select " +
            (isOpen ? "is-open " : "") +
            (disabled ? "is-disabled " : "")
          }
          onClick={handleToggle}
          onBlur={handleBlur}
          ref={triggerRef}
          id={name}
          role="combobox"
          aria-expanded={isOpen}
          aria-haspopup="listbox"
          aria-disabled={disabled}
          tabIndex={disabled ? -1 : 0}
        >
          <div
            className={
              "trieoh-input__container-field trieoh-custom-select__trigger " +
              (!selectedLabel ? "placeholder" : "")
            }
          >
            {displayValue}
          </div>
          <RiArrowDownSLine
            className={`trieoh-input__container-icon trieoh-custom-select__arrow ${
              isOpen ? "is-open" : ""
            }`}
            size={24}
          />
        </div>

        {isOpen && (
          <div
            className="trieoh-custom-select__dropdown"
            role="listbox"
            aria-label={`Opções para ${label}`}
          >
            <div className="trieoh-custom-select__options">
              {options.map((opt) => (
                <div
                  key={opt.id}
                  className={`trieoh-custom-select__option ${
                    opt.value === value ? "is-selected" : ""
                  }`}
                  onClick={() => handleSelect(opt.value)}
                  role="option"
                  aria-selected={opt.value === value}
                  tabIndex={0}
                  onKeyDown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                      e.preventDefault();
                      handleSelect(opt.value);
                    }
                  }}
                >
                  {opt.label}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      <div className="trieoh-input__hint">
        {rulesStatus.map((r, i) => {
          const classes = [
            "hint-part",
            r.passed ? "passed" : "",
            !r.passed && submitted ? "failed-on-submit" : "",
          ]
            .filter(Boolean)
            .join(" ");
          return (
            <p key={r.id ?? i} className={classes}>
              {r.message}
            </p>
          );
        })}
      </div>
    </div>
  );
}