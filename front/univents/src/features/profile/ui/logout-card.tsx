import { useAuthActions } from "@trieoh/front-core";
import { LogOut } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/shared/ui/shadcn/button";

export function LogoutCard() {
  const { handleLogout } = useAuthActions();
  const [isLoading, setIsLoading] = useState(false);

  const onLogout = async () => {
    if (isLoading) return;
    setIsLoading(true);
    try {
      await handleLogout();
    } catch (_error) {
      toast.error("Ocorreu um erro ao tentar sair da sessão.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="pt-4 pb-2">
      <div className="flex flex-col gap-4 rounded-xl border border-destructive/20 bg-destructive/5 p-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="space-y-1">
          <p className="text-sm font-medium leading-none">Encerrar sessão</p>
          <p className="text-sm text-muted-foreground">
            Ao sair, a sessão ativa será encerrada.
          </p>
        </div>
        <Button
          variant="destructive"
          onClick={onLogout}
          disabled={isLoading}
          className="w-full shadow-sm sm:w-auto"
        >
          <LogOut className="mr-2 size-4" />
          {isLoading ? "Saindo..." : "Sair da conta"}
        </Button>
      </div>
    </div>
  );
}
