import type { JSX } from "@solidjs/web";
import { Combobox, type ComboboxOption } from "@trieoh/ui-solid";

export interface ToolbarComboboxOption {
  value: string;
  label: string;
  description?: string;
}

export interface ToolbarComboboxProps {
  value?: string;
  options: readonly ToolbarComboboxOption[];
  placeholder: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  icon?: JSX.Element;
  iconOnly?: boolean;
  class?: string;
  className?: string;
  triggerClass?: string;
  triggerClassName?: string;
  dropdownClass?: string;
  dropdownClassName?: string;
  onChange: (value: string) => void;
}

export function ToolbarCombobox(props: ToolbarComboboxProps): JSX.Element {
  return (
    <Combobox
      size="sm"
      value={props.value}
      options={props.options as readonly ComboboxOption[]}
      placeholder={props.placeholder}
      searchPlaceholder={props.searchPlaceholder}
      disabled={props.disabled}
      icon={props.icon}
      iconOnly={props.iconOnly}
      class={props.class}
      className={props.className}
      triggerClass={props.triggerClass}
      triggerClassName={props.triggerClassName}
      dropdownClass={props.dropdownClass}
      dropdownClassName={props.dropdownClassName}
      onChange={props.onChange}
    />
  );
}
