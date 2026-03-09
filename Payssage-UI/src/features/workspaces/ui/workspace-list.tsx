import { Plus, Globe, Settings, MoreHorizontal, ArrowRight, Copy, Check } from 'lucide-react'
import { Link } from '@tanstack/react-router'
import { Button } from '#/shared/ui/shadcn/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger
} from '#/shared/ui/shadcn/dropdown-menu'
import type { WorkspaceI } from '../model'
import { toast } from 'sonner'
import { useState } from 'react'

interface WorkspaceListProps {
  workspaces: WorkspaceI[]
}

export function WorkspaceList({ workspaces }: WorkspaceListProps) {
  const [copiedId, setCopiedId] = useState<string | null>(null)

  const copyToClipboard = (id: string) => {
    navigator.clipboard.writeText(id)
    setCopiedId(id)
    toast.success('ID copied to clipboard')
    setTimeout(() => setCopiedId(null), 2000)
  }

  if (workspaces.length === 0) {
    return (
      <Card className="rounded-sm border-dashed flex flex-col items-center justify-center p-8 md:p-16 text-center bg-muted/30">
        <div className="w-12 h-12 rounded-none bg-primary/10 flex items-center justify-center mb-4">
          <Globe className="w-6 h-6 text-primary" />
        </div>
        <CardTitle className="text-xl md:text-2xl mb-2 font-bold tracking-tight">No workspaces found</CardTitle>
        <CardDescription className="max-w-xs mb-6 text-sm">
          Create your first workspace to start processing payments and managing keys.
        </CardDescription>
        <Button className="rounded-sm gap-2 h-10 px-6">
          <Plus className="w-4 h-4" />
          Create Workspace
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-extrabold tracking-tight">Workspaces</h2>
          <p className="text-muted-foreground text-sm md:text-base">Manage your payment environments and integrations.</p>
        </div>
        <Button className="rounded-sm gap-2 h-10 sm:w-auto w-full">
          <Plus className="w-4 h-4" />
          New Workspace
        </Button>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {workspaces.map((workspace) => (
          <Link
            key={workspace.id}
            to="/admin/$name"
            params={{ name: workspace.name }}
            className="block group"
          >
            <Card className="rounded-none border border-border group-hover:border-primary group-hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,0.1)] transition-all bg-card flex flex-col h-full group-hover:-translate-x-0.5 group-hover:-translate-y-0.5">
              <CardHeader className="p-5 pb-4">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex flex-col gap-1">
                    <CardTitle className="text-xl font-black uppercase tracking-tighter leading-none mb-1">
                      {workspace.name}
                    </CardTitle>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className="w-fit h-4 px-1 rounded-none text-[9px] font-mono uppercase tracking-wider text-muted-foreground border-muted-foreground/30">
                        {workspace.id}
                      </Badge>
                    </div>
                  </div>

                  <div className="flex items-center gap-2 lg:opacity-0 lg:group-hover:opacity-100 transition-opacity">
                    <DropdownMenu>
                      <DropdownMenuTrigger>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 rounded-none hover:bg-muted"
                          onClick={(e) => {
                            e.preventDefault()
                            e.stopPropagation()
                          }}
                        >
                          <MoreHorizontal className="w-4 h-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent
                        align="end"
                        className="rounded-none min-w-48"
                        onClick={(e) => {
                          e.preventDefault()
                          e.stopPropagation()
                        }}
                      >
                        <DropdownMenuItem className="gap-2 text-xs font-bold uppercase tracking-widest rounded-none py-2 px-3">
                          <Settings className="w-3.5 h-3.5" />
                          {workspace.sandbox ? 'Disable Sandbox' : 'Enable Sandbox'}
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className="gap-2 text-xs font-bold uppercase tracking-widest rounded-none py-2 px-3"
                          onClick={() => copyToClipboard(workspace.id)}
                        >
                          {copiedId === workspace.id ? (
                            <Check className="w-3.5 h-3.5 text-emerald-500" />
                          ) : (
                            <Copy className="w-3.5 h-3.5" />
                          )}
                          Copy ID
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                    <ArrowRight className="w-5 h-5 text-primary" />
                  </div>
                </div>
              </CardHeader>

              <CardContent className="p-5 pt-0 mt-auto">
                <div className="flex items-center justify-between text-[10px] font-black uppercase tracking-widest">
                  <span className={`flex items-center gap-1.5 ${workspace.sandbox ? 'text-amber-600' : 'text-emerald-600'}`}>
                    <div className={`w-1.5 h-1.5 rounded-full ${workspace.sandbox ? 'bg-amber-500' : 'bg-emerald-500'}`} />
                    {workspace.sandbox ? 'Sandbox' : 'Production'}
                  </span>
                  <span className="text-muted-foreground">
                    {new Date(workspace.created_at).toLocaleDateString('en-US', { month: 'short', year: 'numeric' })}
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
