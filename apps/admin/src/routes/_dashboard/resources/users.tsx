import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { usersResource } from '@/resources/users'

export const Route = createFileRoute('/_dashboard/resources/users')({
  component: () => <ResourcePage resource={usersResource} />,
})
