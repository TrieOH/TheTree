import {
  X,
  User,
  Clock,
  Star,
  FileText,
  Mail,
  Phone,
  Link2,
  Hash,
  ToggleLeft,
  Calendar,
  List
} from "lucide-react";
import type { FullFieldI, FullFormI } from "../model";
import type { FieldTypeI } from "#/features/fields/model";
import { cn } from "#/shared/lib/utils";

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
  if (value === "") return <span className="text-sm text-gray-400 italic">No response</span>;

  switch (fieldType) {
    case "int": {
      const num = parseInt(value, 10);
      if (isNaN(num)) return <span className="text-sm text-gray-900">{value}</span>;
      if (num >= 1 && num <= 5) {
        return (
          <div className="flex items-center gap-1">
            {Array.from({ length: 5 }).map((_, i) => (
              <Star
                key={i}
                className={cn("w-4 h-4", i < num ? "text-yellow-400 fill-yellow-400" : "text-gray-300")}
              />
            ))}
            <span className="ml-2 text-sm font-medium text-gray-900">{num}/5</span>
          </div>
        );
      }
      return <span className="text-sm text-gray-900">{num}</span>;
    }
    case "float": {
      const floatNum = parseFloat(value);
      return <span className="text-sm text-gray-900">{isNaN(floatNum) ? value : floatNum.toFixed(2)}</span>;
    }
    case "bool": {
      const boolVal = value.toLowerCase() === "true" || value === "1";
      return (
        <span className={cn("inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium", boolVal ? "bg-green-50 text-green-700" : "bg-red-50 text-red-700")}>
          {boolVal ? "Yes" : "No"}
        </span>
      );
    }
    case "email":
      return <a href={`mailto:${value}`} className="text-sm text-blue-600 hover:underline">{value}</a>;
    case "url":
      return <a href={value} target="_blank" rel="noopener noreferrer" className="text-sm text-blue-600 hover:underline break-all">{value}</a>;
    case "phone":
      return <a href={`tel:${value}`} className="text-sm text-blue-600 hover:underline">{value}</a>;
    case "date":
      return <span className="text-sm text-gray-900">{new Date(value).toLocaleDateString()}</span>;
    case "datetime":
      return <span className="text-sm text-gray-900">{new Date(value).toLocaleString()}</span>;
    case "time":
      return <span className="text-sm text-gray-900">{value}</span>;
    default:
      return <span className="text-sm text-gray-900 whitespace-pre-wrap">{value}</span>;
  }
}

export function SubmissionDetail({ fullForm, responder, onClose }: SubmissionDetailProps) {
  if (!responder) return null;

  const allFields = fullForm.steps.flatMap((s) => s.fields);
  const answers = allFields
    .flatMap((f) => f.answers)
    .filter((a) => a.responder === responder);

  const completedAt = answers.length > 0
    ? new Date(Math.max(...answers.map((a) => new Date(a.answer.answered_at).getTime())))
    : null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-2xl bg-white h-full shadow-2xl overflow-y-auto animate-in slide-in-from-right duration-300">
        <div className="sticky top-0 z-10 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">Submission Review</h2>
            <p className="text-sm text-gray-500 mt-0.5">Response #{responder.split("@")[0]}</p>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <div className="px-6 py-4 border-b border-gray-200 bg-gray-50/50">
          <div className="grid grid-cols-2 gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-white rounded-lg border border-gray-200">
                <User className="w-4 h-4 text-gray-600" />
              </div>
              <div>
                <div className="text-xs text-gray-500 uppercase tracking-wider font-medium">Respondent</div>
                <div className="text-sm font-medium text-gray-900">{responder}</div>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-white rounded-lg border border-gray-200">
                <Clock className="w-4 h-4 text-gray-600" />
              </div>
              <div>
                <div className="text-xs text-gray-500 uppercase tracking-wider font-medium">Submitted</div>
                <div className="text-sm font-medium text-gray-900">
                  {completedAt?.toLocaleString("en-US", {
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
        </div>

        <div className="px-6 py-6 space-y-8">
          {fullForm.steps.map((fullStep, stepIdx) => {
            const stepAnswers = fullStep.fields.map((f) => ({
              field: f.field,
              value: getFieldAnswer(fullStep.fields, f.field.id, responder),
            }));

            const hasAnyAnswer = stepAnswers.some((a) => a.value !== null && a.value !== "");
            if (!hasAnyAnswer) return null;

            return (
              <div key={fullStep.step.id} className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="flex items-center justify-center w-8 h-8 rounded-full bg-blue-600 text-white text-sm font-semibold">
                    {stepIdx + 1}
                  </div>
                  <div>
                    <h3 className="text-base font-semibold text-gray-900">{fullStep.step.title}</h3>
                    {fullStep.step.description && (
                      <p className="text-xs text-gray-500">{fullStep.step.description}</p>
                    )}
                  </div>
                </div>

                <div className="ml-11 space-y-4">
                  {stepAnswers.map(({ field, value }) => {
                    const TypeIcon = typeIcons[field.type];
                    return (
                      <div key={field.id} className="bg-gray-50 rounded-lg p-4 border border-gray-100">
                        <div className="flex items-start justify-between mb-1">
                          <label className="text-sm font-medium text-gray-700">
                            {field.title}
                            {field.required && <span className="text-red-500 ml-1">*</span>}
                          </label>
                          <span className="inline-flex items-center gap-1 text-[10px] uppercase tracking-wider text-gray-400 font-medium">
                            <TypeIcon className="w-3 h-3" />
                            {field.type}
                          </span>
                        </div>
                        {field.description && (
                          <p className="text-xs text-gray-500 mb-2">{field.description}</p>
                        )}
                        <div className="mt-1">
                          {value === null || value === "" ? (
                            <span className="text-sm text-gray-400 italic">No response</span>
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
      </div>
    </div>
  );
}