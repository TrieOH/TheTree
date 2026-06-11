import FormContainer from '#/features/submissions/ui/form-view/form-container';
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/view/$formID')({
  component: RouteComponent,
})

function RouteComponent() {
  const { formID } = Route.useParams()

  return (
    <main className='flex min-h-screen items-center justify-center bg-background p-4'>
      <FormContainer formId={formID} />
    </main>
  )
}
