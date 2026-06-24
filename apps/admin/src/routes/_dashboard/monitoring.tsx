import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { StatsGrid } from '@/components/gateway/stats-grid'

export const Route = createFileRoute('/_dashboard/monitoring')({
  component: MonitoringPage,
})

function MonitoringPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Monitoring</h1>
        <p className="text-muted-foreground">Real-time gateway metrics and health</p>
      </div>

      <StatsGrid
        totalRequests={1247}
        avgLatency={14200}
        activeModels={1}
        uptime="2h 34m"
      />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Request Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Success (2xx)</span>
                <span className="text-green-500">1,180 (94.6%)</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Client Error (4xx)</span>
                <span className="text-yellow-500">45 (3.6%)</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Server Error (5xx)</span>
                <span className="text-red-500">22 (1.8%)</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Model Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>mimo-auto</span>
                <span>1,247 requests</span>
              </div>
              <div className="flex justify-between text-sm text-muted-foreground">
                <span>Total tokens</span>
                <span>156,420</span>
              </div>
              <div className="flex justify-between text-sm text-muted-foreground">
                <span>Avg tokens/request</span>
                <span>125</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
