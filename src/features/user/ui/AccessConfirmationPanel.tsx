import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";
import { Check, X } from "lucide-react";

interface PropsI {
  title: string;
  subTitle: string;
  state: "success" | "error",
  onExit: () => void;
}

export default function AccessConfirmationPanel({ 
  title, 
  subTitle,
  state,
  onExit
}: PropsI) {
  return (
    <div className="flex flex-col justify-center items-center gap-3 py-2">
      { state === "success" ?
        <Check 
          strokeWidth="4px" 
          size={32} 
          className="p-2 border-2 rounded-full border-green-400 text-green-400" 
        /> :
        <X 
          strokeWidth="4px" 
          size={32} 
          className="p-2 border-2 rounded-full border-destructive text-destructive" 
        />
      }
      <span className="text-primary font-medium text-sm text-center">{title}</span>
      <span className="text-muted-foreground text-xs">{subTitle}</span>
      <ShadowButton 
        value="Done" 
        onClick={onExit} 
        variant="ghost" 
        className="text-primary outline-0 font-medium"
      />
    </div>
  )
}