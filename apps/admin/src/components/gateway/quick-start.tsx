import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export function QuickStart() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Quick Start</CardTitle>
      </CardHeader>
      <CardContent>
        <pre className="p-4 rounded-lg bg-muted font-mono text-sm overflow-x-auto">
{`curl http://localhost:4200/v1/chat/completions \\
  -H "Content-Type: application/json" \\
  -d '{"model":"mimo-auto","messages":[{"role":"user","content":"hi"}]}'`}
        </pre>
      </CardContent>
    </Card>
  )
}
