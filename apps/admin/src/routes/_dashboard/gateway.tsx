import { createFileRoute } from '@tanstack/react-router'
import { StatsCard } from '@/components/widgets/stats-card'
import { Activity, Zap, Clock, Server } from 'lucide-react'
import { useState, useEffect } from 'react'

export const Route = createFileRoute('/_dashboard/gateway')({
  component: GatewayPage,
})

interface GatewayStats {
  totalRequests: number
  avgLatency: number
  activeModels: number
  uptime: string
}

interface RecentRequest {
  id: string
  model: string
  tokens: number
  latency: number
  timestamp: string
  status: 'success' | 'error'
}

function GatewayPage() {
  const [stats, setStats] = useState<GatewayStats>({
    totalRequests: 0,
    avgLatency: 0,
    activeModels: 1,
    uptime: '0m',
  })
  const [recentRequests, setRecentRequests] = useState<RecentRequest[]>([])
  const [isConnected, setIsConnected] = useState(false)

  // Simulated data - replace with real API calls
  useEffect(() => {
    // Check gateway health
    fetch('/health')
      .then(res => res.json())
      .then(data => {
        setIsConnected(data.status === 'ok')
      })
      .catch(() => setIsConnected(false))

    // Simulated stats
    setStats({
      totalRequests: 1247,
      avgLatency: 14200,
      activeModels: 1,
      uptime: '2h 34m',
    })

    // Simulated recent requests
    setRecentRequests([
      { id: '1', model: 'mimo-auto', tokens: 156, latency: 12400, timestamp: '2m ago', status: 'success' },
      { id: '2', model: 'mimo-auto', tokens: 89, latency: 14800, timestamp: '5m ago', status: 'success' },
      { id: '3', model: 'mimo-auto', tokens: 234, latency: 16200, timestamp: '8m ago', status: 'success' },
      { id: '4', model: 'mimo-auto', tokens: 67, latency: 11900, timestamp: '12m ago', status: 'error' },
      { id: '5', model: 'mimo-auto', tokens: 178, latency: 13500, timestamp: '15m ago', status: 'success' },
    ])
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">MiMo Gateway</h1>
          <p className="text-muted-foreground">OpenAI-compatible proxy for MiMo Auto Free API</p>
        </div>
        <div className="flex items-center gap-2">
          <div className={`h-2 w-2 rounded-full ${isConnected ? 'bg-success' : 'bg-danger'}`} />
          <span className="text-sm text-muted-foreground">
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard
          title="Total Requests"
          value={stats.totalRequests.toLocaleString()}
          icon={<Activity className="h-5 w-5" />}
          trend="+12% from last hour"
        />
        <StatsCard
          title="Avg Latency"
          value={`${(stats.avgLatency / 1000).toFixed(1)}s`}
          icon={<Clock className="h-5 w-5" />}
          trend="Binary backend"
        />
        <StatsCard
          title="Active Models"
          value={stats.activeModels.toString()}
          icon={<Zap className="h-5 w-5" />}
          trend="mimo-auto"
        />
        <StatsCard
          title="Uptime"
          value={stats.uptime}
          icon={<Server className="h-5 w-5" />}
          trend="Since last restart"
        />
      </div>

      {/* Quick Start */}
      <div className="rounded-lg border bg-card p-6">
        <h2 className="text-lg font-semibold mb-4">Quick Start</h2>
        <div className="bg-muted rounded-lg p-4 font-mono text-sm">
          <p className="text-muted-foreground mb-2"># Test the gateway</p>
          <code className="text-accent">
            curl http://localhost:4200/v1/chat/completions \<br />
            &nbsp;&nbsp;-H "Content-Type: application/json" \<br />
            &nbsp;&nbsp;-d {'{'}"model":"mimo-auto","messages":[{'}'}{'}'}{']'}{']'}
          </code>
        </div>
      </div>

      {/* Recent Requests */}
      <div className="rounded-lg border bg-card">
        <div className="p-4 border-b">
          <h2 className="text-lg font-semibold">Recent Requests</h2>
        </div>
        <div className="divide-y">
          {recentRequests.map((req) => (
            <div key={req.id} className="p-4 flex items-center justify-between hover:bg-muted/50">
              <div className="flex items-center gap-4">
                <div className={`h-2 w-2 rounded-full ${req.status === 'success' ? 'bg-success' : 'bg-danger'}`} />
                <div>
                  <p className="font-medium">{req.model}</p>
                  <p className="text-sm text-muted-foreground">{req.tokens} tokens</p>
                </div>
              </div>
              <div className="text-right">
                <p className="font-mono text-sm">{(req.latency / 1000).toFixed(1)}s</p>
                <p className="text-xs text-muted-foreground">{req.timestamp}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* API Reference */}
      <div className="rounded-lg border bg-card p-6">
        <h2 className="text-lg font-semibold mb-4">API Reference</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="p-4 rounded-lg bg-muted">
            <h3 className="font-medium mb-2">Chat Completions</h3>
            <code className="text-sm text-accent">POST /v1/chat/completions</code>
            <p className="text-sm text-muted-foreground mt-2">OpenAI-compatible chat endpoint</p>
          </div>
          <div className="p-4 rounded-lg bg-muted">
            <h3 className="font-medium mb-2">List Models</h3>
            <code className="text-sm text-accent">GET /v1/models</code>
            <p className="text-sm text-muted-foreground mt-2">Available models</p>
          </div>
          <div className="p-4 rounded-lg bg-muted">
            <h3 className="font-medium mb-2">Health Check</h3>
            <code className="text-sm text-accent">GET /health</code>
            <p className="text-sm text-muted-foreground mt-2">Gateway status</p>
          </div>
        </div>
      </div>
    </div>
  )
}
