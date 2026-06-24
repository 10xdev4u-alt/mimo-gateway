import { createFileRoute } from '@tanstack/react-router'
import { ResourcePage } from '@/components/resource/resource-page'
import { blogsResource } from '@/resources/blogs'

export const Route = createFileRoute('/_dashboard/resources/blogs')({
  component: () => <ResourcePage resource={blogsResource} />,
})
