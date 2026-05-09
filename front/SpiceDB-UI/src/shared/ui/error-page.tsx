import { ArrowLeft, RefreshCw, XCircle } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { useNavigate } from '@tanstack/react-router'
import type { ErrorComponentProps } from '@tanstack/react-router'

export default function ErrorPage({
  error,
  reset,
}: ErrorComponentProps) {
  const navigate = useNavigate()
  const message = error instanceof Error ? error.message : String(error)

  return (
    <main className="min-h-screen flex flex-col items-center justify-center gap-6 px-6 py-12 text-center">
      <XCircle className="w-16 h-16 text-destructive" />
      <div className="max-w-2xl">
        <h1 className="text-3xl font-semibold">
          Something went wrong
        </h1>
        <p className="mt-3 text-sm text-muted-foreground">{message || 'Ocorreu um erro inesperado. Tente novamente.'}</p>
      </div>
      <div className="flex flex-wrap justify-center gap-3">
        <Button variant="secondary" className="cursor-pointer" onClick={() => navigate({ to: '/' })}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Go Back
        </Button>
        <Button className="cursor-pointer" onClick={() => reset()}>
          <RefreshCw className="mr-2 h-4 w-4" /> Try Again
        </Button>
      </div>
    </main>
  )
}
