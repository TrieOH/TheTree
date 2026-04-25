import { Plus, FileText } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'
import type { FormI } from '../model'
import { Badge } from '#/shared/ui/shadcn/badge'

interface FormListProps {
  forms: FormI[]
  openModal: () => void;
  projectID: string;
}

export function FormList({
  forms,
  openModal,
  projectID,
}: FormListProps) {

  if (forms.length === 0) {
    return (
      <Card className="rounded-sm border-dashed flex flex-col items-center justify-center p-8 md:p-16 text-center bg-muted/30">
        <div className="w-12 h-12 rounded-none bg-primary/10 flex items-center justify-center mb-4">
          <FileText className="w-6 h-6 text-primary" />
        </div>
        <CardTitle className="text-xl md:text-2xl mb-2 font-bold tracking-tight">
          No forms found
        </CardTitle>
        <CardDescription className="max-w-xs mb-6 text-sm">
          Create your first form to start collecting data.
        </CardDescription>
        <Button className="rounded-sm gap-2 h-10 px-6" onClick={openModal}>
          <Plus className="w-4 h-4" />
          Create Form
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-extrabold tracking-tight">Forms</h2>
          <p className="text-muted-foreground text-sm md:text-base">Manage your forms and view responses.</p>
        </div>
        <Button className="rounded-sm gap-2 h-10 sm:w-auto w-full" onClick={openModal}>
          <Plus className="w-4 h-4" />
          New Form
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {forms.map((form) => (
          <Link
            key={form.id}
            to="/admin/$project" // Update this later if there's a specific form route
            params={{ project: projectID }}
            className="block group"
          >
            <Card className="rounded-none border border-border group-hover:border-primary group-hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] transition-all bg-card flex flex-col h-full group-hover:-translate-x-0.5 group-hover:-translate-y-0.5 min-w-0 overflow-hidden">
              <CardHeader className="p-4 sm:p-5 pb-3 sm:pb-4 min-w-0">
                <div className="grid grid-cols-[1fr_auto] items-start gap-2 min-w-0">
                  <div className="flex flex-col gap-1 min-w-0 overflow-hidden">
                    <CardTitle className="text-lg sm:text-xl font-black uppercase tracking-tighter leading-none mb-1 truncate">
                      {form.title}
                    </CardTitle>
                    <div className="flex items-center gap-2 min-w-0">
                      <span className="text-[10px] font-mono text-muted-foreground/50 truncate">
                        #{form.id}
                      </span>
                    </div>
                  </div>
                  <Badge variant={form.status === 'open' ? 'default' : 'secondary'} className="rounded-none uppercase text-[10px]">
                    {form.status}
                  </Badge>
                </div>
              </CardHeader>

              <CardContent className="p-4 sm:p-5 pt-0 mt-auto">
                <div className="flex items-center justify-between text-[10px] font-black uppercase tracking-widest gap-2">
                  <span className="text-muted-foreground truncate">
                    {new Date(form.created_at).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}
                  </span>
                </div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>
    </div>
  )
}
