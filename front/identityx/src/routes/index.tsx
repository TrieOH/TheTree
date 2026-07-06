import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowRight, Lock, Shield } from 'lucide-react'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  const { auth } = Route.useRouteContext()
  const isAuthenticated = auth?.isAuthenticated

  return (
    <div className="min-h-screen bg-linear-to-br from-background via-background to-primary/5">
      {/* Nav */}
      <nav className="flex items-center justify-between px-6 py-4 max-w-6xl mx-auto">
        <span className="text-xl font-bold tracking-tight">IdentityX</span>
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <Link
              to="/admin"
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              Dashboard
            </Link>
          ) : (
            <>
              <Link
                to="/auth"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Sign in
              </Link>
              <Link
                to="/auth"
                className="text-sm px-4 py-2 bg-primary text-primary-foreground rounded-lg hover:opacity-90 transition-opacity font-medium"
              >
                Get started
              </Link>
            </>
          )}
        </div>
      </nav>

      {/* Hero */}
      <section className="max-w-6xl mx-auto px-6 pt-24 pb-16 text-center">
        <h1 className="text-5xl md:text-6xl font-bold tracking-tight mb-6">
          Auth infrastructure
          <br />
          <span className="bg-linear-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            for your next project
          </span>
        </h1>
        <p className="text-lg text-muted-foreground max-w-xl mx-auto mb-10 leading-relaxed">
          IdentityX is a modern, headless authentication platform that integrates
          seamlessly with any stack. OAuth, sessions, organizations — all out of the box.
        </p>
        <div className="flex items-center justify-center gap-4">
          {isAuthenticated ? (
            <Link
              to="/admin"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-xl font-medium hover:opacity-90 transition-opacity"
            >
              Go to Dashboard
              <ArrowRight size={18} />
            </Link>
          ) : (
            <>
              <Link
                to="/auth"
                className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-xl font-medium hover:opacity-90 transition-opacity"
              >
                Start building
                <ArrowRight size={18} />
              </Link>
              <Link
                to="/auth"
                className="inline-flex items-center gap-2 px-6 py-3 border border-border rounded-xl font-medium hover:bg-muted/50 transition-colors"
              >
                View docs
              </Link>
            </>
          )}
        </div>
      </section>

      {/* Features */}
      <section className="max-w-6xl mx-auto px-6 pb-24">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            {
              icon: Shield,
              title: "OAuth / SSO",
              description:
                "Google and GitHub out of the box — plus any OpenID Connect provider. One integration, all protocols.",
            },
            {
              icon: Lock,
              title: "Session management",
              description:
                "Secure HTTP-only cookies, refresh token rotation, device tracking, and instant revocation.",
            },
            {
              icon: ArrowRight,
              title: "API-first design",
              description:
                "Headless APIs, SDKs for React / TypeScript / Go. Bring your own UI or use our components.",
            },
          ].map((feat) => (
            <div
              key={feat.title}
              className="rounded-xl border border-border bg-card p-6 hover:shadow-md transition-shadow"
            >
              <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <feat.icon size={20} className="text-primary" />
              </div>
              <h3 className="font-semibold mb-2">{feat.title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">
                {feat.description}
              </p>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="border-t border-border">
        <div className="max-w-6xl mx-auto px-6 py-16 text-center">
          <h2 className="text-2xl font-bold mb-4">Ready to ship auth?</h2>
          <Link
            to={isAuthenticated ? "/admin" : "/auth"}
            className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-xl font-medium hover:opacity-90 transition-opacity"
          >
            {isAuthenticated ? "Go to Dashboard" : "Get started"}
            <ArrowRight size={18} />
          </Link>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border">
        <div className="max-w-6xl mx-auto px-6 py-6 flex items-center justify-between text-sm text-muted-foreground">
          <span>IdentityX by TrieOH</span>
          <Link
            to={isAuthenticated ? "/admin" : "/auth"}
            className="hover:text-foreground transition-colors"
          >
            {isAuthenticated ? "Dashboard" : "Sign in"}
          </Link>
        </div>
      </footer>
    </div>
  )
}