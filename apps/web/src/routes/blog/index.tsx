import { createFileRoute, Link } from '@tanstack/react-router'
import { useBlogs } from '@/hooks/use-blogs'

export const Route = createFileRoute('/blog/')({
  component: BlogListPage,
})

function BlogListPage() {
  const { data: blogs, isLoading } = useBlogs()

  return (
    <div className="max-w-4xl mx-auto py-16 px-6">
      <h1 className="text-4xl font-bold mb-2">Blog</h1>
      <p className="text-muted-foreground mb-10">Latest articles and updates.</p>

      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-32 rounded-xl bg-card/50 animate-pulse" />
          ))}
        </div>
      ) : !blogs?.length ? (
        <p className="text-muted-foreground">No posts yet. Check back soon!</p>
      ) : (
        <div className="space-y-6">
          {blogs.map((blog: any) => (
            <Link
              key={blog.id}
              to="/blog/$slug"
              params={{ slug: blog.slug }}
              className="block rounded-xl border border-border/40 bg-card/50 p-6 hover:border-accent/30 transition-colors"
            >
              <h2 className="text-xl font-semibold mb-2">{blog.title}</h2>
              <p className="text-sm text-muted-foreground line-clamp-2">{blog.excerpt || blog.content?.substring(0, 150)}</p>
              <span className="text-xs text-muted-foreground/50 mt-3 block">
                {new Date(blog.created_at).toLocaleDateString()}
              </span>
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
