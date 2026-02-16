"use client"

import { useState } from "react"
import { useStore } from "@/lib/store"
import { Headphones, AlertCircle } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"

export function LoginScreen() {
  const { loginWithCredentials, forgotPassword } = useStore()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const [forgotOpen, setForgotOpen] = useState(false)
  const [forgotEmail, setForgotEmail] = useState("")
  const [forgotError, setForgotError] = useState("")
  const [forgotSuccess, setForgotSuccess] = useState(false)
  const [forgotLoading, setForgotLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)
    try {
      await loginWithCredentials(email, password)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed")
    } finally {
      setLoading(false)
    }
  }

  const handleForgotSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setForgotError("")
    setForgotSuccess(false)
    setForgotLoading(true)
    try {
      await forgotPassword(forgotEmail)
      setForgotSuccess(true)
    } catch (err) {
      setForgotError(err instanceof Error ? err.message : "Request failed")
    } finally {
      setForgotLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <div className="mb-4 flex items-center justify-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded border border-primary/30 bg-primary/10">
              <Headphones className="h-5 w-5 text-primary" />
            </div>
            <h1 className="text-2xl font-bold tracking-wider text-foreground uppercase">
              SupportFIX
            </h1>
          </div>
          <p className="text-sm tracking-wide text-muted-foreground uppercase">
            Sign in to your account
          </p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
              <AlertCircle className="h-4 w-4 shrink-0" />
              {error}
            </div>
          )}

          <div className="space-y-2">
            <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter email"
              className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
              required
            />
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Password
              </label>
              <button
                type="button"
                onClick={() => {
                  setForgotOpen(true)
                  setForgotEmail(email)
                  setForgotError("")
                  setForgotSuccess(false)
                }}
                className="text-[11px] text-primary hover:underline"
              >
                Forgot password?
              </button>
            </div>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter password"
              className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-md border border-primary/30 bg-primary/10 p-3 text-sm font-bold tracking-wider text-primary uppercase transition-all hover:bg-primary/20 disabled:opacity-50"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
      </div>

      <Dialog open={forgotOpen} onOpenChange={setForgotOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="uppercase tracking-wider">Forgot Password</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Enter your email address and we&apos;ll send you instructions to reset your password.
          </p>
          <form onSubmit={handleForgotSubmit} className="space-y-4">
            {forgotError && (
              <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-2 text-sm text-destructive">
                <AlertCircle className="h-4 w-4 shrink-0" />
                {forgotError}
              </div>
            )}
            {forgotSuccess && (
              <div className="rounded-md border border-primary/30 bg-primary/10 p-2 text-sm text-primary">
                If an account exists with that email, you will receive password reset instructions shortly.
              </div>
            )}
            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Email
              </label>
              <input
                type="email"
                value={forgotEmail}
                onChange={(e) => setForgotEmail(e.target.value)}
                placeholder="Enter email"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
                required
              />
            </div>
            <button
              type="submit"
              disabled={forgotLoading}
              className="w-full rounded-md border border-primary/30 bg-primary/10 p-3 text-sm font-bold tracking-wider text-primary uppercase transition-all hover:bg-primary/20 disabled:opacity-50"
            >
              {forgotLoading ? "Sending..." : "Send Reset Instructions"}
            </button>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
