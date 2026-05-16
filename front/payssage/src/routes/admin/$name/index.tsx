import { createFileRoute } from '@tanstack/react-router'
import { Globe, Zap, CreditCard, Activity, Clock, CheckCircle2, XCircle, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '#/shared/ui/shadcn/card'
import { Badge } from '#/shared/ui/shadcn/badge'
import { useQuery } from '@tanstack/react-query'
import { allWorkspacePaymentIntentsQueryOptions } from '#/features/payment-intents/api'
import { cn } from '#/shared/lib/utils'

export const Route = createFileRoute('/admin/$name/')({
  component: WorkspaceOverview,
})

function WorkspaceOverview() {
  const { name } = Route.useParams()
  const { data: intents = [], isLoading } = useQuery(
    allWorkspacePaymentIntentsQueryOptions(name),
  )

  const totalPayments = intents.length
  const succeededIntents = intents.filter((i) => i.status === 'succeeded')
  const successRate = totalPayments > 0 ? (succeededIntents.length / totalPayments) * 100 : 0
  const totalVolume = succeededIntents.reduce((acc, i) => acc + i.amount, 0)

  // Get 5 most recent activities
  const recentActivity = [...intents]
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5)

  const formatCurrency = (amount: number, currency: string = 'USD') => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount / 100)
  }

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="space-y-1">
        <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">
          Workspace Overview
        </h2>
        <p className="text-muted-foreground text-xs font-mono uppercase tracking-widest">
          NAME: {name}
        </p>
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
            <div className="text-3xl font-black">{isLoading ? '...' : totalPayments}</div>
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
            <div className="text-3xl font-black">
              {isLoading ? '...' : `${successRate.toFixed(1)}%`}
            </div>
          </CardContent>
        </Card>

        <Card className="rounded-none border-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-muted-foreground flex items-center gap-2">
              <Activity className="w-3 h-3" />
              Recent Intent Count
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black">{isLoading ? '...' : recentActivity.length}</div>
          </CardContent>
        </Card>

        <Card className="rounded-none bg-primary/5">
          <CardHeader className="pb-2">
            <CardTitle className="text-[10px] font-black uppercase tracking-[0.2em] text-primary flex items-center gap-2">
              <CreditCard className="w-3 h-3" />
              Total Volume
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black text-primary">
              {isLoading ? '...' : formatCurrency(totalVolume)}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-1">
        <Card className="rounded-none border-border bg-muted/20 border-dashed">
          <CardHeader className="border-b border-border/50 pb-4">
            <CardTitle className="text-lg font-black uppercase tracking-tight flex items-center gap-2">
              <Clock className="w-5 h-5" />
              Recent Activity
            </CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            {isLoading ? (
              <div className="space-y-4 py-6">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-12 bg-muted/50 animate-pulse" />
                ))}
              </div>
            ) : recentActivity.length > 0 ? (
              <div className="divide-y divide-border/50">
                {recentActivity.map((activity) => (
                  <div key={activity.id} className="py-4 flex items-center justify-between gap-4">
                    <div className="flex items-center gap-4 min-w-0">
                      <div
                        className={cn(
                          'w-10 h-10 flex items-center justify-center shrink-0 border',
                          activity.status === 'succeeded'
                            ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20'
                            : activity.status === 'failed'
                              ? 'bg-destructive/10 text-destructive border-destructive/20'
                              : activity.status === 'pending'
                                ? 'bg-amber-500/10 text-amber-500 border-amber-500/20'
                                : 'bg-muted text-muted-foreground border-border',
                        )}
                      >
                        {activity.status === 'succeeded' && <CheckCircle2 className="w-5 h-5" />}
                        {activity.status === 'failed' && <XCircle className="w-5 h-5" />}
                        {activity.status === 'pending' && <Clock className="w-5 h-5" />}
                        {activity.status === 'cancelled' && <AlertCircle className="w-5 h-5" />}
                      </div>
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-2 mb-1">
                          <span className="font-mono text-[10px] font-black uppercase tracking-tighter text-muted-foreground">
                            ID: {activity.id.split('_').pop()}
                          </span>
                          <Badge
                            variant="outline"
                            className={cn(
                              'rounded-none text-[8px] font-black uppercase tracking-widest px-1.5 h-4',
                              activity.status === 'succeeded' &&
                              'border-emerald-500/30 text-emerald-500 bg-emerald-500/5',
                              activity.status === 'failed' &&
                              'border-destructive/30 text-destructive bg-destructive/5',
                              activity.status === 'pending' &&
                              'border-amber-500/30 text-amber-500 bg-amber-500/5',
                            )}
                          >
                            {activity.status}
                          </Badge>
                        </div>
                        <p className="text-[10px] text-muted-foreground uppercase font-bold tracking-widest">
                          {new Date(activity.created_at).toLocaleString()}
                        </p>
                      </div>
                    </div>
                    <div className="text-right shrink-0">
                      <div className="text-sm font-black tracking-tight">
                        {formatCurrency(activity.amount, activity.currency)}
                      </div>
                      <div className="text-[9px] text-muted-foreground uppercase font-black tracking-widest mt-0.5">
                        {activity.provider}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-[10px] uppercase tracking-widest font-bold text-muted-foreground py-12 text-center">
                No transactions found
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
