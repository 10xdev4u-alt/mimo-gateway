import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { api } from '@/lib/api-client'

export const Route = createFileRoute('/_auth/sign-up')({
  component: SignUpPage,
})

function SignUpPage() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit, formState: { errors } } = useForm()

  const onSubmit = async (data: any) => {
    setLoading(true)
    setError('')
    try {
      const res = await api.post('/api/auth/register', data)
      localStorage.setItem('access_token', res.data.data.tokens.access_token)
      localStorage.setItem('refresh_token', res.data.data.tokens.refresh_token)
      navigate({ to: '/dashboard' })
    } catch (err: any) {
      setError(err.response?.data?.error?.message || 'Registration failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="w-full max-w-md mx-auto p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold">Create account</h1>
        <p className="text-muted-foreground mt-2">Get started with your new account</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        {error && (
          <div className="p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-sm text-destructive">
            {error}
          </div>
        )}

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1.5">First Name</label>
            <input type="text" {...register('first_name', { required: 'Required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1.5">Last Name</label>
            <input type="text" {...register('last_name', { required: 'Required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Email</label>
          <input type="email" {...register('email', { required: 'Email is required' })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="you@example.com" />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1.5">Password</label>
          <input type="password" {...register('password', { required: 'Password is required', minLength: { value: 8, message: 'Min 8 characters' } })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="••••••••" />
          {errors.password && <p className="text-xs text-destructive mt-1">{String(errors.password.message)}</p>}
        </div>

        <button type="submit" disabled={loading} className="w-full py-2.5 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors disabled:opacity-50">
          {loading ? 'Creating account...' : 'Create account'}
        </button>

        <p className="text-center text-sm text-muted-foreground">
          Already have an account? <Link to="/login" className="text-accent hover:underline">Sign in</Link>
        </p>
      </form>
    </div>
  )
}
