import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

export const Route = createFileRoute('/_dashboard/logs')({
  component: LogsPage,
})

const logs = [
  { id: '1', time: '14:32:15', method: 'POST', path: '/v1/chat/completions', status: 200, latency: '12.4s', model: 'mimo-auto' },
  { id: '2', time: '14:31:42', method: 'GET', path: '/v1/models', status: 200, latency: '2ms', model: '-' },
  { id: '3', time: '14:30:18', method: 'POST', path: '/v1/chat/completions', status: 200, latency: '14.8s', model: 'mimo-auto' },
  { id: '4', time: '14:29:55', method: 'GET', path: '/health', status: 200, latency: '1ms', model: '-' },
  { id: '5', time: '14:28:03', method: 'POST', path: '/v1/chat/completions', status: 500, latency: '30.0s', model: 'mimo-auto' },
]

function LogsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Request Logs</h1>
        <p className="text-muted-foreground">Monitor gateway requests and errors</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Recent Requests</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {logs.map((log) => (
              <div key={log.id} className="flex items-center gap-4 p-3 rounded-lg bg-muted/50 hover:bg-muted">
                <code className="text-xs text-muted-foreground w-16">{log.time}</code>
                <Badge variant={log.method === 'POST' ? 'default' : 'secondary'} className="w-16 justify-center">
                  {log.method}
                </Badge>
                <code className="text-sm flex-1">{log.path}</code>
                <Badge variant={log.status === 200 ? 'default' : 'destructive'} className="w-12 justify-center">
                  {log.status}
                </Badge>
                <span className="text-sm text-muted-foreground w-16 text-right">{log.latency}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
