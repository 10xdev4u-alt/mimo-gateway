import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <div className="relative">
      {/* Hero */}
      <section className="relative py-24 px-6">
        <div className="max-w-5xl mx-auto text-center">
          <span className="inline-flex items-center rounded-full bg-accent/10 border border-accent/20 px-4 py-1.5 text-sm font-medium text-accent mb-6">
            Go + React. Built with Grit.
          </span>
          <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-6">
            <span className="text-foreground">Build faster with</span>{' '}
            <span className="bg-gradient-to-r from-accent to-purple-400 bg-clip-text text-transparent">
              mimo-gateway
            </span>
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto mb-10 leading-relaxed">
            A full-stack application powered by Go, React, and the Grit framework.
            Production-ready from day one.
          </p>
          <div className="flex items-center justify-center gap-4">
            <Link
              to="/blog"
              className="inline-flex items-center px-6 py-3 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors"
            >
              Read the Blog
            </Link>
            <a
              href="/api/health"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center px-6 py-3 rounded-lg border border-border text-foreground font-medium hover:bg-muted transition-colors"
            >
              API Health Check
            </a>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-20 px-6 border-t border-border/30">
        <div className="max-w-5xl mx-auto">
          <h2 className="text-3xl font-bold text-center mb-12">What{"'"}s Included</h2>
          <div className="grid md:grid-cols-3 gap-6">
            {[
              { title: 'Go API', desc: 'Gin + GORM with JWT auth, RBAC, file uploads, and background jobs.' },
              { title: 'React Frontend', desc: 'TanStack Router + React Query + Tailwind CSS. Fast and lightweight.' },
              { title: 'Full-Stack DX', desc: 'Shared types, one-command resource generation, hot reload everywhere.' },
            ].map((f) => (
              <div key={f.title} className="rounded-xl border border-border/40 bg-card/50 p-6">
                <h3 className="text-lg font-semibold mb-2">{f.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  )
}
