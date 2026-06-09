import {
  Clock,
  Star,
  FileText,
  Mail,
  Phone,
  Link2,
  Hash,
  ToggleLeft,
  Calendar,
  List,
  User,
} from "lucide-react";
import type { FullFieldI, FullFormI } from "../model";
import type { FieldTypeI } from "#/features/fields/model";
import { cn } from "#/shared/lib/utils";
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
} from "#/shared/ui/shadcn/drawer";

interface SubmissionDetailProps {
  fullForm: FullFormI;
  responder: string | null;
  onClose: () => void;
}

function getFieldAnswer(fields: FullFieldI[], fieldId: string, responder: string): string | null {
  const field = fields.find((f) => f.field.id === fieldId);
  if (!field) return null;
  const answer = field.answers.find((a) => a.responder === responder);
  return answer?.answer.answer ?? null;
}

const typeIcons: Record<FieldTypeI, typeof FileText> = {
  string: FileText,
  email: Mail,
  int: Hash,
  float: Hash,
  bool: ToggleLeft,
  date: Calendar,
  time: Clock,
  datetime: Calendar,
  select: List,
  file: FileText,
  phone: Phone,
  url: Link2,
};

function renderFieldValue(fieldType: FieldTypeI, value: string): React.ReactNode {
  if (value === "") return <span className="text-sm text-muted-foreground italic">No response</span>;

  switch (fieldType) {
    case "int": {
      const num = parseInt(value, 10);
      if (isNaN(num)) return <span className="text-sm text-foreground">{value}</span>;
      if (num >= 1 && num <= 5) {
        return (
          <div className="flex items-center gap-1">
            {Array.from({ length: 5 }).map((_, i) => (
              <Star
                key={i}
                className={cn("w-4 h-4", i < num ? "text-accent fill-accent" : "text-muted")}
              />
            ))}
            <span className="ml-2 text-sm font-medium text-foreground">{num}/5</span>
          </div>
        );
      }
      return <span className="text-sm text-foreground">{num}</span>;
    }
    case "float": {
      const floatNum = parseFloat(value);
      return <span className="text-sm text-foreground">{isNaN(floatNum) ? value : floatNum.toFixed(2)}</span>;
    }
    case "bool": {
      const boolVal = value.toLowerCase() === "true" || value === "1";
      return (
        <span className={cn(
          "inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider",
          boolVal ? "bg-primary/10 text-primary" : "bg-destructive/10 text-destructive"
        )}>
          {boolVal ? "Yes" : "No"}
        </span>
      );
    }
    case "email":
      return <a href={`mailto:${value}`} className="text-sm text-secondary hover:text-primary transition-colors underline decoration-secondary/30 underline-offset-4">{value}</a>;
    case "url":
      return <a href={value} target="_blank" rel="noopener noreferrer" className="text-sm text-secondary hover:text-primary transition-colors underline decoration-secondary/30 underline-offset-4 break-all">{value}</a>;
    case "phone":
      return <a href={`tel:${value}`} className="text-sm text-secondary hover:text-primary transition-colors underline decoration-secondary/30 underline-offset-4">{value}</a>;
    case "date":
      return <span className="text-sm text-foreground tabular-nums">{new Date(value).toLocaleDateString()}</span>;
    case "datetime":
      return <span className="text-sm text-foreground tabular-nums">{new Date(value).toLocaleString()}</span>;
    case "time":
      return <span className="text-sm text-foreground tabular-nums">{value}</span>;
    default:
      return <span className="text-sm text-foreground whitespace-pre-wrap">{value}</span>;
  }
}

export function SubmissionDetail({ fullForm, responder, onClose }: SubmissionDetailProps) {
  const isOpen = !!responder;

  const allFields = fullForm.steps.flatMap((s) => s.fields ?? []);
  const answers = responder ? allFields
    .flatMap((f) => f.answers)
    .filter((a) => a.responder === responder) : [];

  const completedAt = answers.length > 0
    ? new Date(Math.max(...answers.map((a) => new Date(a.answer.answered_at).getTime())))
    : null;

  return (
    <Drawer open={isOpen} onOpenChange={(open) => !open && onClose()} direction="right">
      <DrawerContent className="h-full flex flex-col focus:outline-none rounded-none border-none">
        <DrawerHeader className="border-b border-border bg-card py-6">
          <div className="space-y-1">
            <DrawerTitle className="text-xl font-bold tracking-tight">Submission Review</DrawerTitle>
            <DrawerDescription className="font-medium">
              {responder ? `Response from ${responder.split("@")[0]}` : ""}
            </DrawerDescription>
          </div>
        </DrawerHeader>

        <div className="flex-1 overflow-y-auto px-6 py-6 space-y-8">
          {/* Metadata Stack */}
          {responder && (
            <div className="flex flex-col gap-3 pb-8 border-b border-border/60">
              <div className="flex items-center gap-3 p-4 rounded-lg bg-muted/40 border border-border/40">
                <div className="p-2.5 bg-primary/10 text-primary rounded-lg">
                  <User className="w-5 h-5" />
                </div>
                <div className="min-w-0">
                  <div className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Respondent</div>
                  <div className="text-sm font-semibold text-foreground truncate">{responder}</div>
                </div>
              </div>
              <div className="flex items-center gap-3 p-4 rounded-lg bg-muted/40 border border-border/40">
                <div className="p-2.5 bg-accent/10 text-accent rounded-lg">
                  <Clock className="w-5 h-5" />
                </div>
                <div className="min-w-0">
                  <div className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Submitted At</div>
                  <div className="text-sm font-semibold text-foreground tabular-nums">
                    {completedAt?.toLocaleString(undefined, {
                      month: "short",
                      day: "numeric",
                      year: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    }) ?? "—"}
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Form Content */}
          {responder && fullForm.steps.map((fullStep, stepIdx) => {
            const stepAnswers = fullStep.fields?.map((f) => ({
              field: f.field,
              value: getFieldAnswer(fullStep.fields ?? [], f.field.id, responder),
            })) ?? [];

            const hasAnyAnswer = stepAnswers.some((a) => a.value !== null && a.value !== "");
            if (!hasAnyAnswer) return null;

            return (
              <div key={fullStep.step.id} className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
                <div className="flex items-center gap-4">
                  <div className="flex items-center justify-center shrink-0 w-10 h-10 rounded-full bg-primary text-primary-foreground text-sm font-bold shadow-sm">
                    {stepIdx + 1}
                  </div>
                  <div className="min-w-0">
                    <h3 className="text-base font-bold text-foreground leading-tight truncate">{fullStep.step.title}</h3>
                    {fullStep.step.description && (
                      <p className="text-xs text-muted-foreground mt-0.5 line-clamp-1">{fullStep.step.description}</p>
                    )}
                  </div>
                </div>

                <div className="sm:ml-14 space-y-4">
                  {stepAnswers.map(({ field, value }) => {
                    const TypeIcon = typeIcons[field.type];
                    return (
                      <div key={field.id} className="group bg-card hover:bg-accent/5 rounded-xl p-5 border border-border/60 hover:border-accent/30 transition-all shadow-sm">
                        <div className="flex items-start justify-between gap-4 mb-2">
                          <label className="text-sm font-bold text-foreground leading-snug">
                            {field.title}
                            {field.required && <span className="text-destructive ml-1">*</span>}
                          </label>
                          <span className="shrink-0 inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md bg-muted text-[10px] font-bold uppercase tracking-wider text-muted-foreground group-hover:bg-accent/10 group-hover:text-accent transition-colors">
                            <TypeIcon className="w-3 h-3" />
                            {field.type}
                          </span>
                        </div>
                        {field.description && (
                          <p className="text-xs text-muted-foreground mb-3">{field.description}</p>
                        )}
                        <div className="mt-1 pt-3 border-t border-border/40 group-hover:border-accent/20 transition-colors">
                          {value === null || value === "" ? (
                            <span className="text-sm text-muted-foreground italic">No response</span>
                          ) : (
                            renderFieldValue(field.type, value)
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      </DrawerContent>
    </Drawer>
  );
}