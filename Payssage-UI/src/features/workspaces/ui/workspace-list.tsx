import { Plus, Globe, Settings, Calendar, MoreHorizontal, TestTube2, ArrowRight, Copy, Check } from 'lucide-react'
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
          <Card key={workspace.id} className="rounded-sm border-border hover:border-primary/50 transition-colors bg-card group flex flex-col">
            <CardHeader className="p-5 pb-4">
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-center gap-3">
                  <div className={
                    `w-10 h-10 rounded-none flex items-center justify-center transition-colors 
                    ${workspace.sandbox ? 'bg-amber-500/10 text-amber-600' : 'bg-primary/10 text-primary'}`
                  }>
                    {workspace.sandbox ? <TestTube2 className="w-5 h-5" /> : <Globe className="w-5 h-5" />}
                  </div>
                  <div className="flex flex-col">
                    <CardTitle className="text-lg font-bold leading-none mb-1">
                      {workspace.name}
                    </CardTitle>
                    <button
                      onClick={() => copyToClipboard(workspace.id)}
                      className="flex items-center gap-1 group/id w-fit"
                    >
                      <Badge variant="outline" className="w-fit h-5 px-1.5 rounded-none text-[10px] font-mono uppercase tracking-wider text-muted-foreground border-muted-foreground/30 transition-colors group-hover/id:border-primary/50 group-hover/id:text-primary">
                        ID: {workspace.id.slice(0, 8)}...
                      </Badge>
                      {copiedId === workspace.id ? (
                        <Check className="w-3 h-3 text-emerald-500" />
                      ) : (
                        <Copy className="w-3 h-3 text-muted-foreground opacity-0 group-hover/id:opacity-100 transition-opacity" />
                      )}
                    </button>
                  </div>
                </div>

                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <Button variant="ghost" size="icon" className="h-8 w-8 rounded-none">
                      <MoreHorizontal className="w-4 h-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="rounded-none min-w-48 max-w-[90vw]">
                    <DropdownMenuItem className="gap-2 text-sm rounded-none py-2 px-3">
                      <Settings className="w-3.5 h-3.5" />
                      {workspace.sandbox ? 'Disable Sandbox' : 'Enable Sandbox'}
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="gap-2 text-sm rounded-none py-2 px-3"
                      onClick={() => copyToClipboard(workspace.id)}
                    >
                      <Copy className="w-3.5 h-3.5" />
                      Copy ID
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </CardHeader>

            <CardContent className="p-5 pt-0 mt-auto">
              <div className="space-y-3 mb-6">
                <div className="flex items-center justify-between text-xs text-muted-foreground border-b border-border/50 pb-2">
                  <span className="flex items-center gap-1.5">
                    <Calendar className="w-3.5 h-3.5" />
                    Created at
                  </span>
                  <span className="font-medium text-foreground">
                    {new Date(workspace.created_at).toLocaleDateString('en-US')}
                  </span>
                </div>
                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span className="flex items-center gap-1.5">
                    <div className={`w-1.5 h-1.5 rounded-full ${workspace.sandbox ? 'bg-amber-500' : 'bg-emerald-500'}`} />
                    Environment
                  </span>
                  <span className={`font-bold ${workspace.sandbox ? 'text-amber-600' : 'text-emerald-600'}`}>
                    {workspace.sandbox ? 'SANDBOX' : 'PRODUCTION'}
                  </span>
                </div>
              </div>

              <Button variant="default" className="w-full rounded-none h-10 flex items-center justify-between px-4 group/btn">
                <span>Manage</span>
                <ArrowRight className="w-4 h-4 transition-transform group-hover/btn:translate-x-1" />
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
