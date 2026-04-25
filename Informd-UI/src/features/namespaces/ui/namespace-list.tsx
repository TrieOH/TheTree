import { Plus, Globe } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { Button } from '#/shared/ui/shadcn/button'
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle 
} from '#/shared/ui/shadcn/card'
import type { NamespaceI } from '../model'

interface NamespaceListProps {
  namespaces: NamespaceI[]
  openModal: () => void;
}

export default function NamespaceList({
  namespaces,
  openModal,
}: NamespaceListProps) {

  if (namespaces.length === 0) {
    return (
      <Card className="rounded-sm border-dashed flex flex-col items-center justify-center p-8 md:p-16 text-center bg-muted/30">
        <div className="w-12 h-12 rounded-none bg-primary/10 flex items-center justify-center mb-4">
          <Globe className="w-6 h-6 text-primary" />
        </div>
        <CardTitle className="text-xl md:text-2xl mb-2 font-bold tracking-tight">
          No namespaces found
        </CardTitle>
        <CardDescription className="max-w-xs mb-6 text-sm">
          Create your first namespace to start managing forms.
        </CardDescription>
        <Button className="rounded-sm gap-2 h-10 px-6" onClick={openModal}>
          <Plus className="w-4 h-4" />
          Create Namespace
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-extrabold tracking-tight">Namespace</h2>
          <p className="text-muted-foreground text-sm md:text-base">Manage your forms environments and integrations.</p>
        </div>
        <Button className="rounded-sm gap-2 h-10 sm:w-auto w-full" onClick={openModal}>
          <Plus className="w-4 h-4" />
          New Namespace
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {namespaces.map((namespace) => (
          <Link
            key={namespace.id}
            to="/admin/$namespaceID"
            params={{ namespaceID: namespace.id }}
            className="block group"
          >
            <Card className="rounded-none border border-border group-hover:border-primary group-hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] transition-all bg-card flex flex-col h-full group-hover:-translate-x-0.5 group-hover:-translate-y-0.5 min-w-0 overflow-hidden">
              <CardHeader className="p-4 sm:p-5 pb-3 sm:pb-4 min-w-0">
                <div className="grid grid-cols-[1fr_auto] items-start gap-2 min-w-0">
                  <div className="flex flex-col gap-1 min-w-0 overflow-hidden">
                    <CardTitle className="text-lg sm:text-xl font-black uppercase tracking-tighter leading-none mb-1 truncate">
                      {namespace.name}
                    </CardTitle>
                    <div className="flex items-center gap-2 min-w-0">
                      <span className="text-[10px] font-mono text-muted-foreground/50 truncate">
                        #{namespace.id}
                      </span>
                    </div>
                  </div>
                </div>
              </CardHeader>

              <CardContent className="p-4 sm:p-5 pt-0 mt-auto">
                <div className="flex items-center justify-between text-[10px] font-black uppercase tracking-widest gap-2">
                  <span className="text-muted-foreground truncate">
                    {new Date(namespace.created_at).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}
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
