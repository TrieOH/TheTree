import { Copy } from "lucide-react";
import { truncateString } from "../lib/utils";
import { toast } from "sonner";
import { ShadowButton } from "./buttons/ShadowButton";

interface TruncatedIdProps {
  id: string;
}

export default function TruncatedId({ id }: TruncatedIdProps) {
  const truncatedId = truncateString(id, 5, 4);

  const handleCopy = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id);
    toast.success("ID copied to clipboard");
  };

  return (
    <div className="flex items-center gap-2">
      <span>{truncatedId}</span>
      <ShadowButton
        variant="ghost"
        onClick={handleCopy}
        className="p-0"
        leftIcon={<Copy className="h-4 w-4"/>}
      />
    </div>
  );
}
