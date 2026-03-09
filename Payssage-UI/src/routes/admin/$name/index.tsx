import { createFileRoute } from '@tanstack/react-router'
import { Globe, Shield, Zap, CreditCard } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'

export const Route = createFileRoute('/admin/$name/')({
  component: WorkspaceOverview,
})

function WorkspaceOverview() {
  const { name } = Route.useParams()

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="space-y-1">
        <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Workspace Overview</h2>
        <p className="text-muted-foreground text-xs font-mono uppercase tracking-widest">ID: {name}</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="rounded-none border-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground flex items-center gap-2">
              <Zap className="w-3 h-3" />
              Total Payments
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">0</div>
          </CardContent>
        </Card>

        <Card className="rounded-none border-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground flex items-center gap-2">
              <Globe className="w-3 h-3" />
              Success Rate
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">0%</div>
          </CardContent>
        </Card>

        <Card className="rounded-none border-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground flex items-center gap-2">
              <Shield className="w-3 h-3" />
              Active Keys
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">0</div>
          </CardContent>
        </Card>

        <Card className="rounded-none border-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground flex items-center gap-2">
              <CreditCard className="w-3 h-3" />
              Volume
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">$0.00</div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card className="rounded-none border-border bg-muted/20 border-dashed">
          <CardHeader>
            <CardTitle className="text-lg font-black uppercase tracking-tight">Recent Activity</CardTitle>
            <div className="text-[10px] uppercase tracking-widest font-bold text-muted-foreground mt-1">No transactions found</div>
          </CardHeader>
        </Card>

        <Card className="rounded-none border-border bg-muted/20 border-dashed">
          <CardHeader>
            <CardTitle className="text-lg font-black uppercase tracking-tight">Integration</CardTitle>
            <div className="text-[10px] uppercase tracking-widest font-bold text-muted-foreground mt-1">Configure your API keys to start</div>
          </CardHeader>
        </Card>
      </div>
    </div>
  )
}
