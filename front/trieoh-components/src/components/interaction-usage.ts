import type { Values } from "../catalog";
import exampleSource from "../examples/InteractionExamples.tsx?raw";
import { previewUsage } from "./preview-usage";
const examples: Record<string, string> = {
  Hotkeys: "HotkeysExample",
  Kbd: "KbdExample",
  DropdownMenu: "DropdownExample",
  Popover: "PopoverExample",
  DatePicker: "DateExample",
  DateRangePicker: "DateRangeExample",
  FileUpload: "FileUploadExample",
  CommandPalette: "CommandExample",
  Timeline: "TimelineExample",
};
export function interactionUsage(
  name: string,
  values: Values,
): string | undefined {
  const example = examples[name];
  return example
    ? previewUsage(
        exampleSource,
        example,
        values,
        exampleSource.slice(
          exampleSource.indexOf("const localDate ="),
          exampleSource.indexOf("export function DropdownExample"),
        ),
      )
    : undefined;
}
