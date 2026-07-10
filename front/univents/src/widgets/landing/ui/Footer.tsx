import { Link } from '@tanstack/react-router'

export function Footer() {
  return (
    <footer className="border-t border-border py-8 md:py-12 pb-24!">
      <div className="max-w-5xl mx-auto flex flex-col md:flex-row justify-between items-center gap-4 md:gap-6">
        <div className="text-xs md:text-sm text-muted-foreground/70 order-2 md:order-1">
          © 2026 Univents. Todos os direitos reservados.
        </div>
        <div className="flex gap-6 md:gap-8 text-xs md:text-sm order-1 md:order-2">
          <Link to="/terms" className="text-muted-foreground hover:text-foreground transition-colors">
            Termos
          </Link>
          <Link to="/privacy" className="text-muted-foreground hover:text-foreground transition-colors">
            Privacidade
          </Link>
          <Link to="/contact" className="text-muted-foreground hover:text-foreground transition-colors">
            Contato
          </Link>
        </div>
      </div>
    </footer>
  )
}
