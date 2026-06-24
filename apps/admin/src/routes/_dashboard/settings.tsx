import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export const Route = createFileRoute('/_dashboard/settings')({
  component: SettingsPage,
})

function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Settings</h1>
        <p className="text-muted-foreground">Configure your gateway settings</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Gateway Configuration</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label>Binary Path</Label>
            <Input placeholder="/path/to/.mimocode" />
          </div>
          <div className="space-y-2">
            <Label>Rate Limit (requests/minute)</Label>
            <Input type="number" defaultValue={100} />
          </div>
          <div className="space-y-2">
            <Label>Default Model</Label>
            <Input defaultValue="mimo-auto" />
          </div>
          <Button>Save Changes</Button>
        </CardContent>
      </Card>
    </div>
  )
}
