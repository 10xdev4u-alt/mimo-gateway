import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'

interface Request {
  id: string
  model: string
  tokens: number
  latency: number
  timestamp: string
  status: 'success' | 'error'
}

interface RequestTableProps {
  requests: Request[]
}

export function RequestTable({ requests }: RequestTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Status</TableHead>
          <TableHead>Model</TableHead>
          <TableHead>Tokens</TableHead>
          <TableHead>Latency</TableHead>
          <TableHead>Time</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {requests.map((req) => (
          <TableRow key={req.id}>
            <TableCell>
              <div className={`h-2 w-2 rounded-full ${req.status === 'success' ? 'bg-green-500' : 'bg-red-500'}`} />
            </TableCell>
            <TableCell className="font-medium">{req.model}</TableCell>
            <TableCell>{req.tokens}</TableCell>
            <TableCell className="font-mono">{(req.latency / 1000).toFixed(1)}s</TableCell>
            <TableCell className="text-muted-foreground">{req.timestamp}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
