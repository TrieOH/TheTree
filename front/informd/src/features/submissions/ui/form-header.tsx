import type { FormI } from "#/features/forms/model";
import { Clock, Calendar } from "lucide-react";

interface FormHeaderProps {
  form: FormI;
}

export default function FormHeader({ form }: FormHeaderProps) {
  return (
    <div className="bg-card border-b border-border shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3">
        <div className="flex justify-center items-center flex-wrap xs:justify-between! gap-2 text-sm">
          {form.opened_at && (
            <div className="flex items-center gap-2">
              <div className="p-1.5 rounded-xs bg-primary/10 text-primary">
                <Calendar className="w-4 h-4" />
              </div>
              <div className="flex items-center gap-1.5">
                <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground/80">
                  Opened
                </span>
                <span className="font-semibold text-foreground text-xs sm:text-sm">
                  {new Date(form.opened_at).toLocaleDateString(undefined, {
                    day: "2-digit",
                    month: "short",
                    year: "numeric",
                  })}
                </span>
              </div>
            </div>
          )}

          {form.closed_at && (
            <div className="flex items-center gap-2">
              <div className="p-1.5 rounded-xs bg-accent/10 text-accent">
                <Clock className="w-4 h-4" />
              </div>
              <div className="flex items-center gap-1.5">
                <span className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground/80">
                  Closes
                </span>
                <span className="font-semibold text-foreground text-xs sm:text-sm">
                  {new Date(form.closed_at).toLocaleDateString(undefined, {
                    day: "2-digit",
                    month: "short",
                    year: "numeric",
                  })}
                </span>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}