export function Footer() {
  return (
    <footer className="border-t border-neutral-200 py-12 px-6">
      <div className="max-w-5xl mx-auto flex flex-col md:flex-row justify-between items-center gap-6">
        <div className="text-sm text-neutral-500">
          © 2026 Univents. Todos os direitos reservados.
        </div>
        <div className="flex gap-8 text-sm">
          <a href="#" className="text-neutral-500 hover:text-neutral-900 transition-colors">Termos</a>
          <a href="#" className="text-neutral-500 hover:text-neutral-900 transition-colors">Privacidade</a>
          <a href="#" className="text-neutral-500 hover:text-neutral-900 transition-colors">Contato</a>
        </div>
      </div>
    </footer>
  )
}