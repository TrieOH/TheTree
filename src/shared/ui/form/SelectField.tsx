import { useFieldContext } from "@/shared/lib/forms";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/ui/shadcn/select";

interface Option {
  label: string;
  value: string;
}

interface PropsI {
  label: string;
  placeholder?: string;
  options: Option[];
}

export default function SelectField({ label, placeholder, options }: PropsI) {
  const field = useFieldContext<string>();

  return (
    <div className="flex flex-col gap-2 mb-4">
      <label 
        htmlFor={field.name}
        className="text-sm font-medium text-foreground"
      >
        {label}
      </label>
      <Select
        value={field.state.value}
        onValueChange={(v) => field.handleChange(v)}
      >
        <SelectTrigger id={field.name} className="w-full">
          <SelectValue placeholder={placeholder || "Select an option"} />
        </SelectTrigger>
        <SelectContent>
          {options.map((option) => (
            <SelectItem key={option.value} value={option.value}>
              {option.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
