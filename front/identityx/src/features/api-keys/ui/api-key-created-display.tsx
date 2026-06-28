import { Copy, KeySquare } from "lucide-react";
import { toast } from "sonner";
import { ShadowButton } from "@/shared/ui/buttons/ShadowButton";

interface ApiKeyCreatedDisplayProps {
  name: string;
  rawKey: string;
  onClose: () => void;
}

export function ApiKeyCreatedDisplay({ name, rawKey, onClose }: ApiKeyCreatedDisplayProps) {
  const handleCopy = () => {
    navigator.clipboard.writeText(rawKey);
    toast.success("API key copied to clipboard");
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3 p-4 rounded-sm bg-amber-500/10 border border-amber-500/20">
        <KeySquare className="size-5 text-amber-500 shrink-0" />
        <p className="text-xs text-amber-600 font-medium">
          Copy this key now. You won't be able to see it again.
        </p>
      </div>

      <div>
        <p className="text-xs font-medium text-muted-foreground mb-1">Key Name</p>
        <p className="text-sm font-semibold">{name}</p>
      </div>

      <div>
        <p className="text-xs font-medium text-muted-foreground mb-1">API Key</p>
        <div className="flex items-center gap-2 p-3 rounded-sm bg-muted border border-border">
          <code className="flex-1 text-xs font-mono break-all select-all">
            {rawKey}
          </code>
          <ShadowButton
            variant="ghost"
            onClick={handleCopy}
            className="p-1.5 h-auto shrink-0"
            leftIcon={<Copy className="size-4" />}
          />
        </div>
      </div>

      <div className="flex justify-end pt-2">
        <ShadowButton
          variant="default"
          onClick={onClose}
          className="rounded-sm font-medium text-xs px-6"
          value="Done"
        />
      </div>
    </div>
  );
}