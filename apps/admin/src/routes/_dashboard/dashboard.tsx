import { createFileRoute } from '@tanstack/react-router'
import { StatsCard } from '@/components/widgets/stats-card'
import { useAuth } from '@/hooks/use-auth'
import { Users, FileText, Upload, Activity } from 'lucide-react'

export const Route = createFileRoute('/_dashboard/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const { user } = useAuth()

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">Welcome back, {user?.first_name || 'Admin'}</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard title="Total Users" value="--" icon={<Users className="h-5 w-5" />} />
        <StatsCard title="Blog Posts" value="--" icon={<FileText className="h-5 w-5" />} />
        <StatsCard title="Uploads" value="--" icon={<Upload className="h-5 w-5" />} />
        <StatsCard title="Active Sessions" value="--" icon={<Activity className="h-5 w-5" />} />
      </div>
    </div>
  )
}
