import { Link } from "@tanstack/solid-router";
import Logo from "@/shared/ui/Logo";
export function Footer() {
  return (
    <footer class="border-t border-border py-8 pb-24 md:py-12">
      <div class="mx-auto flex max-w-5xl flex-col items-center justify-between gap-4 md:flex-row md:gap-6">
        <div class="order-2 flex items-center gap-3 text-xs text-muted-foreground/70 md:order-1 md:text-sm">
          <div class="h-6 w-6">
            <Logo variant="icon" imgClassName="grayscale opacity-70" />
          </div>
          <span>© 2026 Univents. Todos os direitos reservados.</span>
        </div>
        <nav class="order-1 flex gap-6 text-xs md:order-2 md:gap-8 md:text-sm">
          <Link
            to="/"
            search={{ as: "guest" }}
            class="text-muted-foreground hover:text-foreground"
          >
            Início
          </Link>
          <Link
            to="/auth"
            search={{}}
            class="text-muted-foreground hover:text-foreground"
          >
            Entrar
          </Link>
        </nav>
      </div>
    </footer>
  );
}
