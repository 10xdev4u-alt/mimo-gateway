import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_dashboard/system/files')({
  component: FilesPage,
})

function FilesPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">File Manager</h1>
      <div className="rounded-xl border border-border/40 bg-card/50 p-6">
        <p className="text-muted-foreground">System page content will be loaded here.</p>
      </div>
    </div>
  )
}
