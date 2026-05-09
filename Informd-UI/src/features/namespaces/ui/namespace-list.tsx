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
      <Card className="rounded-none border-4 border-dashed border-border flex flex-col items-center justify-center p-12 md:p-24 text-center bg-muted/20 relative overflow-hidden">
        <div className="absolute top-0 left-0 w-full h-1 bg-border/50" />
        <div className="absolute top-0 right-0 w-32 h-32 bg-primary/5 -mr-16 -mt-16 -rotate-45" />

        <div className="w-20 h-20 rounded-none bg-primary text-primary-foreground flex items-center justify-center mb-8 border-4 border-primary shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]">
          <Globe className="w-10 h-10" />
        </div>
        <CardTitle className="text-3xl md:text-4xl mb-4 font-black uppercase tracking-tighter">
          No namespaces found
        </CardTitle>
        <CardDescription className="max-w-sm mb-10 text-sm uppercase tracking-widest font-bold opacity-60">
          Create your first namespace to start organizing your forms in dedicated environments.
        </CardDescription>
        <Button className="rounded-none gap-3 h-14 px-10 font-black uppercase tracking-[0.2em] shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:translate-x-1 hover:translate-y-1 hover:shadow-none transition-all" onClick={openModal}>
          <Plus className="w-5 h-5" />
          Create First Namespace
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Namespaces</h2>
          <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">Manage your forms environments and integrations.</p>
        </div>
        <Button className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all sm:w-auto w-full" onClick={openModal}>
          <Plus className="w-4 h-4" />
          New Namespace
        </Button>
      </div>

      <div className="grid gap-6 grid-cols-[repeat(auto-fill,minmax(min(100%,320px),1fr))]">
        {namespaces.map((namespace) => (
          <Link
            key={namespace.id}
            to="/admin/$namespaceID"
            params={{ namespaceID: namespace.id }}
            className="block group"
          >
            <Card className="rounded-none border-2 border-border group-hover:border-primary group-hover:shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] transition-all bg-card flex flex-col h-full group-hover:-translate-x-1 group-hover:-translate-y-1 min-w-0 overflow-hidden relative">
              <CardHeader className="p-6 pb-4 min-w-0">
                <div className="grid grid-cols-[1fr_auto] items-start gap-4 min-w-0">
                  <div className="flex flex-col gap-1 min-w-0 overflow-hidden">
                    <CardTitle className="text-xl sm:text-2xl font-black uppercase tracking-tighter leading-none mb-1 truncate">
                      {namespace.name}
                    </CardTitle>
                    <div className="flex items-center gap-2 min-w-0">
                      <span className="text-[10px] font-mono font-bold text-muted-foreground/60 truncate bg-muted/50 px-1.5 py-0.5">
                        ID: {namespace.id}
                      </span>
                    </div>
                  </div>
                </div>
              </CardHeader>

              <CardContent className="p-6 pt-0 mt-auto">
                <div className="flex items-center justify-between text-[10px] font-black uppercase tracking-widest gap-2 pt-4 border-t border-border/40">
                  <span className="text-muted-foreground flex items-center gap-1.5">
                    <div className="w-1 h-1 bg-primary" />
                    {new Date(namespace.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
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
