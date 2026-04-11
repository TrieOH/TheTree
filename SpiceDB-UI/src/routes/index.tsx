import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({ component: App })

function App() {
  return (
    <main className="min-h-screen bg-background text-foreground">
    </main>
  )
}
