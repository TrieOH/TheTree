import { requireAuth } from '#/features/auths/lib/route-guard';
import { FormContainer } from '#/features/submissions/ui/form-view/form-container';
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/view/')({
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <main className='flex min-h-screen items-center justify-center bg-background p-4'>
      <FormContainer />
    </main>
  )
}
