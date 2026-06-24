import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Activity, Zap, Clock, Server } from 'lucide-react'

interface StatsGridProps {
  totalRequests: number
  avgLatency: number
  activeModels: number
  uptime: string
}

export function StatsGrid({ totalRequests, avgLatency, activeModels, uptime }: StatsGridProps) {
  const stats = [
    { title: 'Total Requests', value: totalRequests.toLocaleString(), icon: Activity, trend: '+12%' },
    { title: 'Avg Latency', value: `${(avgLatency / 1000).toFixed(1)}s`, icon: Clock, trend: 'Binary backend' },
    { title: 'Active Models', value: activeModels.toString(), icon: Zap, trend: 'mimo-auto' },
    { title: 'Uptime', value: uptime, icon: Server, trend: 'Since restart' },
  ]

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat) => (
        <Card key={stat.title}>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
            <stat.icon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stat.value}</div>
            <p className="text-xs text-muted-foreground">{stat.trend}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
