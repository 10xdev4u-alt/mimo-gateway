import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '@/hooks/use-auth'

export const Route = createFileRoute('/_dashboard/profile')({
  component: ProfilePage,
})

function ProfilePage() {
  const { user } = useAuth()

  return (
    <div className="max-w-2xl">
      <h1 className="text-2xl font-bold mb-6">Profile</h1>
      <div className="rounded-xl border border-border/40 bg-card/50 p-6 space-y-4">
        <div>
          <label className="text-sm text-muted-foreground">Name</label>
          <p className="font-medium">{user?.first_name} {user?.last_name}</p>
        </div>
        <div>
          <label className="text-sm text-muted-foreground">Email</label>
          <p className="font-medium">{user?.email}</p>
        </div>
        <div>
          <label className="text-sm text-muted-foreground">Role</label>
          <p className="font-medium">{user?.role}</p>
        </div>
      </div>
    </div>
  )
}
