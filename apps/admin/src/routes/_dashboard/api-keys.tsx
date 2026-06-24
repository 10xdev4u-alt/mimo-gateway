import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Plus, Trash2, Copy } from 'lucide-react'
import { useState } from 'react'

export const Route = createFileRoute('/_dashboard/api-keys')({
  component: ApiKeysPage,
})

interface ApiKey {
  id: string
  name: string
  key: string
  created: string
  lastUsed: string
}

function ApiKeysPage() {
  const [keys, setKeys] = useState<ApiKey[]>([
    { id: '1', name: 'Development', key: 'mg_abc123...', created: '2 days ago', lastUsed: '5 min ago' },
    { id: '2', name: 'Production', key: 'mg_def456...', created: '1 week ago', lastUsed: '1 hour ago' },
  ])

  const copyKey = (key: string) => {
    navigator.clipboard.writeText(key)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">API Keys</h1>
          <p className="text-muted-foreground">Manage your API keys for gateway access</p>
        </div>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          Create Key
        </Button>
      </div>

      <div className="space-y-4">
        {keys.map((key) => (
          <Card key={key.id}>
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">{key.name}</h3>
                  <code className="text-sm text-muted-foreground">{key.key}</code>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="sm" onClick={() => copyKey(key.key)}>
                    <Copy className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" className="text-destructive">
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
              <div className="mt-2 text-xs text-muted-foreground">
                Created {key.created} · Last used {key.lastUsed}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
