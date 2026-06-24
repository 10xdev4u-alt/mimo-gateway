import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { api } from '@/lib/api-client'

export const Route = createFileRoute('/_auth/forgot-password')({
  component: ForgotPasswordPage,
})

function ForgotPasswordPage() {
  const [sent, setSent] = useState(false)
  const [loading, setLoading] = useState(false)
  const { register, handleSubmit } = useForm()

  const onSubmit = async (data: any) => {
    setLoading(true)
    try {
      await api.post('/api/auth/forgot-password', data)
      setSent(true)
    } catch {} finally {
      setLoading(false)
    }
  }

  if (sent) {
    return (
      <div className="w-full max-w-md mx-auto p-8 text-center">
        <h1 className="text-3xl font-bold mb-4">Check your email</h1>
        <p className="text-muted-foreground mb-6">We sent a password reset link to your email address.</p>
        <Link to="/login" className="text-accent hover:underline">Back to login</Link>
      </div>
    )
  }

  return (
    <div className="w-full max-w-md mx-auto p-8">
      <div className="text-center mb-8">
        <h1 className="text-3xl font-bold">Reset password</h1>
        <p className="text-muted-foreground mt-2">Enter your email to receive a reset link</p>
      </div>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1.5">Email</label>
          <input type="email" {...register('email', { required: true })} className="w-full px-3 py-2 rounded-lg border border-border bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-accent" placeholder="you@example.com" />
        </div>
        <button type="submit" disabled={loading} className="w-full py-2.5 rounded-lg bg-accent text-white font-medium hover:bg-accent-hover transition-colors disabled:opacity-50">
          {loading ? 'Sending...' : 'Send reset link'}
        </button>
        <p className="text-center text-sm text-muted-foreground">
          <Link to="/login" className="text-accent hover:underline">Back to login</Link>
        </p>
      </form>
    </div>
  )
}
