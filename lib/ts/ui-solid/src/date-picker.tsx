import { Button } from "./button";
import { Popover } from "./popover";
import { Calendar, formatDate, validDate } from "./date-picker/calendar";
import CalendarIcon from "~icons/lucide/calendar-days";
import type { Shape } from "./lib/shape";
export interface DatePickerProps {
  value?: Date | null;
  onChange: (date: Date | null) => void;
  min?: Date;
  max?: Date;
  isDateDisabled?: (date: Date) => boolean;
  placeholder?: string;
  ariaLabel?: string;
  disabled?: boolean;
  clearable?: boolean;
  shape?: Shape;
  class?: string;
}
export function DatePicker(props: DatePickerProps) {
  return (
    <Popover
      ariaLabel={props.ariaLabel ?? "Selecionar data"}
      disabled={props.disabled}
      shape={props.shape}
      width={288}
      trigger={(trigger) => (
        <Button
          {...trigger}
          variant="outline"
          shape={props.shape}
          class={props.class}
          aria-label={props.ariaLabel ?? "Selecionar data"}
        >
          <CalendarIcon class="size-4 shrink-0" />
          <span class="truncate">
            {validDate(props.value)
              ? formatDate(props.value)
              : (props.placeholder ?? "Selecionar data")}
          </span>
        </Button>
      )}
    >
      {(close) => (
        <>
          <Calendar
            value={props.value}
            min={props.min}
            max={props.max}
            isDateDisabled={props.isDateDisabled}
            onSelect={(date) => {
              props.onChange(date);
              close();
            }}
          />
          {props.clearable !== false && (
            <Button
              class="mt-2 w-full"
              size="sm"
              variant="ghost"
              onClick={() => {
                props.onChange(null);
                close();
              }}
            >
              Limpar data
            </Button>
          )}
        </>
      )}
    </Popover>
  );
}
