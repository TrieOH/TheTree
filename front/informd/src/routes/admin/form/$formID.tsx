import { createFileRoute } from '@tanstack/react-router'
import z from 'zod'

const formSearchSchema = z.object({
  namespaceID: z.string().optional(),
})

export const Route = createFileRoute('/admin/form/$formID')({
  validateSearch: (search) => formSearchSchema.parse(search),
  component: FormDetailComponent,
})

function FormDetailComponent() {

  return (
    <div className="p-6">
      <div className="mt-8 p-12 border-2 border-dashed rounded-lg flex items-center justify-center text-muted-foreground bg-muted/20">
        Form coming soon...
      </div>
    </div>
  )
}
