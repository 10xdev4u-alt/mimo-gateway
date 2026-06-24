import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

const endpoints = [
  { method: 'POST', path: '/v1/chat/completions', desc: 'OpenAI-compatible chat endpoint' },
  { method: 'GET', path: '/v1/models', desc: 'List available models' },
  { method: 'GET', path: '/health', desc: 'Gateway health check' },
]

export function ApiReference() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>API Reference</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {endpoints.map((ep) => (
          <div key={ep.path} className="flex items-center gap-3 p-3 rounded-lg bg-muted">
            <span className={`px-2 py-1 text-xs font-bold rounded ${ep.method === 'POST' ? 'bg-blue-500/20 text-blue-400' : 'bg-green-500/20 text-green-400'}`}>
              {ep.method}
            </span>
            <code className="text-sm font-mono text-accent">{ep.path}</code>
            <span className="text-sm text-muted-foreground ml-auto">{ep.desc}</span>
          </div>
        ))}
      </CardContent>
    </Card>
  )
}
