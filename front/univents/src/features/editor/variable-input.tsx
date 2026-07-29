import { AlertCircle, HelpCircle } from "lucide-react";
import { type RefObject, useCallback, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/shadcn/badge";
import { Input } from "@/shared/ui/shadcn/input";
import { DEFAULT_VARIABLES, validateVariables } from "./types";

interface VariableInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  id?: string;
  multiline?: boolean;
  textareaRef?: RefObject<HTMLTextAreaElement | null>;
}

export default function VariableInput({
  value,
  onChange,
  placeholder = "Use {{variavel}} para variáveis",
  label,
  id,
  multiline = false,
  textareaRef,
}: VariableInputProps) {
  const [showHelp, setShowHelp] = useState(false);
  const { invalid, valid } = validateVariables(value, DEFAULT_VARIABLES);

  const insertVariable = useCallback(
    (key: string) => {
      const textareaValue = value;
      const cursorPos =
        textareaRef?.current?.selectionStart ?? textareaValue.length;
      const endPos = textareaRef?.current?.selectionEnd ?? cursorPos;
      const before = textareaValue.slice(0, cursorPos);
      const after = textareaValue.slice(endPos);
      onChange(`${before}{{${key}}}${after}`);
    },
    [value, onChange, textareaRef],
  );

  const inputClassName = cn(
    invalid.length > 0 && "border-destructive focus-visible:ring-destructive",
  );

  return (
    <div className="space-y-2">
      {label && (
        <div className="flex items-center justify-between">
          <label
            htmlFor={id}
            className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
          >
            {label}
          </label>
          <button
            type="button"
            onClick={() => setShowHelp(!showHelp)}
            className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1"
          >
            <HelpCircle className="size-3" />
            {showHelp ? "Fechar ajuda" : "Variáveis disponíveis"}
          </button>
        </div>
      )}

      {multiline ? (
        <textarea
          id={id}
          ref={textareaRef}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className={cn(
            "flex min-h-20 w-full rounded-lg border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
            inputClassName,
          )}
          rows={3}
        />
      ) : (
        <Input
          id={id}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className={inputClassName}
        />
      )}

      {invalid.length > 0 && (
        <div className="flex items-start gap-2 rounded-lg border border-destructive/30 bg-destructive/5 p-2 text-xs text-destructive">
          <AlertCircle className="mt-0.5 size-3.5 shrink-0" />
          <div>
            <span className="font-medium">Variáveis inválidas: </span>
            {invalid.map((v) => (
              <code
                key={v}
                className="mx-0.5 rounded bg-destructive/10 px-1 font-mono"
              >
                {`{{${v}}}`}
              </code>
            ))}
            <span className="block mt-1 text-muted-foreground">
              Use apenas:{" "}
              {DEFAULT_VARIABLES.map((v) => `{{${v.key}}}`).join(", ")}
            </span>
          </div>
        </div>
      )}

      {valid.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {valid.map((v) => {
            const def = DEFAULT_VARIABLES.find((d) => d.key === v);
            return (
              <Badge key={v} variant="secondary" className="text-xs">
                {`{{${v}}}`}
                {def && (
                  <span className="ml-1 text-muted-foreground font-normal">
                    {def.label}
                  </span>
                )}
              </Badge>
            );
          })}
        </div>
      )}

      {showHelp && (
        <div className="rounded-lg border border-border bg-muted/30 p-3 space-y-2">
          <p className="text-xs font-medium text-muted-foreground">
            Variáveis disponíveis para usar com {"{{}}"}:
          </p>
          <div className="grid gap-1.5">
            {DEFAULT_VARIABLES.map((v) => (
              <button
                key={v.key}
                type="button"
                onClick={() => insertVariable(v.key)}
                className="flex items-center justify-between rounded-md border border-border/50 bg-background px-2.5 py-1.5 text-left text-xs hover:bg-accent hover:text-accent-foreground transition-colors"
              >
                <div>
                  <code className="font-mono font-medium">{`{{${v.key}}}`}</code>
                  <span className="ml-2 text-muted-foreground">{v.label}</span>
                </div>
                <span className="text-muted-foreground text-[10px]">
                  {v.description}
                </span>
              </button>
            ))}
          </div>
          <p className="text-[10px] text-muted-foreground">
            Clique em uma variável para inseri-la no campo de texto acima.
          </p>
        </div>
      )}
    </div>
  );
}
