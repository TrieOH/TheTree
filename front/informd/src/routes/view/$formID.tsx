import { requireAuth } from '#/features/auths/lib/route-guard';
import FormContainer from '#/features/submissions/ui/form-view/form-container';
import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

const searchSchema = z.object({
  namespace_id: z.string().optional(),
})

export const Route = createFileRoute('/view/$formID')({
  validateSearch: (search) => searchSchema.parse(search),
  beforeLoad: requireAuth,
  component: RouteComponent,
})

function RouteComponent() {
  const { formID } = Route.useParams()
  const { namespace_id } = Route.useSearch()

  return (
    <main className='flex min-h-screen items-center justify-center bg-background p-4'>
      <FormContainer formId={formID} namespaceId={namespace_id} />
    </main>
  )
}
