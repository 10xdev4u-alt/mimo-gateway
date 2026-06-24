import { useState, useEffect } from 'react'

interface GatewayStats {
  totalRequests: number
  avgLatency: number
  activeModels: number
  uptime: string
}

interface Request {
  id: string
  model: string
  tokens: number
  latency: number
  timestamp: string
  status: 'success' | 'error'
}

export function useGateway() {
  const [stats, setStats] = useState<GatewayStats>({
    totalRequests: 0,
    avgLatency: 0,
    activeModels: 1,
    uptime: '0m',
  })
  const [requests, setRequests] = useState<Request[]>([])
  const [isConnected, setIsConnected] = useState(false)

  useEffect(() => {
    // Check health
    fetch('/health')
      .then(res => res.json())
      .then(data => setIsConnected(data.status === 'ok'))
      .catch(() => setIsConnected(false))

    // Simulated data
    setStats({
      totalRequests: 1247,
      avgLatency: 14200,
      activeModels: 1,
      uptime: '2h 34m',
    })

    setRequests([
      { id: '1', model: 'mimo-auto', tokens: 156, latency: 12400, timestamp: '2m ago', status: 'success' },
      { id: '2', model: 'mimo-auto', tokens: 89, latency: 14800, timestamp: '5m ago', status: 'success' },
      { id: '3', model: 'mimo-auto', tokens: 234, latency: 16200, timestamp: '8m ago', status: 'success' },
      { id: '4', model: 'mimo-auto', tokens: 67, latency: 11900, timestamp: '12m ago', status: 'error' },
      { id: '5', model: 'mimo-auto', tokens: 178, latency: 13500, timestamp: '15m ago', status: 'success' },
    ])
  }, [])

  return { stats, requests, isConnected }
}
