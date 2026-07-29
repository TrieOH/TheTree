import type { EventColor, ProgramI } from "../model";
import { DraggableProgram } from "./DraggableProgram";

interface ProgramListProps {
  programs: ProgramI[];
  programColors: Record<string, EventColor>;
}

export function ProgramList({ programs, programColors }: ProgramListProps) {
  return (
    <div className="border-t border-border pt-3">
      <div className="text-sm font-medium text-foreground mb-2 px-2">
        Programas
      </div>
      {programs.map((p) => {
        const color = programColors[p.id];
        return <DraggableProgram key={p.id} program={p} color={color} />;
      })}
    </div>
  );
}
