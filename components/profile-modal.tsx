"use client"

import { useState, useEffect } from "react"
import { useStore } from "@/lib/store"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { AlertCircle } from "lucide-react"

interface ProfileModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function ProfileModal({ open, onOpenChange }: ProfileModalProps) {
  const { currentUser, changePassword, updateMyProfile } = useStore()
  const [activeTab, setActiveTab] = useState<"profile" | "password">("profile")

  // Profile fields
  const [name, setName] = useState(currentUser?.name ?? "")
  const [phone, setPhone] = useState(currentUser?.phone ?? "")
  const [profileError, setProfileError] = useState("")
  const [profileSuccess, setProfileSuccess] = useState(false)
  const [profileLoading, setProfileLoading] = useState(false)

  // Password fields
  const [currentPassword, setCurrentPassword] = useState("")
  const [newPassword, setNewPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [passwordError, setPasswordError] = useState("")
  const [passwordSuccess, setPasswordSuccess] = useState(false)
  const [passwordLoading, setPasswordLoading] = useState(false)

  const handleProfileSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setProfileError("")
    setProfileSuccess(false)
    if (!name.trim()) {
      setProfileError("Name is required")
      return
    }
    setProfileLoading(true)
    try {
      await updateMyProfile({
        name: name.trim(),
        phone: phone.trim() || undefined,
      })
      setProfileSuccess(true)
    } catch (err) {
      setProfileError(err instanceof Error ? err.message : "Failed to update profile")
    } finally {
      setProfileLoading(false)
    }
  }

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setPasswordError("")
    setPasswordSuccess(false)
    if (newPassword !== confirmPassword) {
      setPasswordError("New passwords do not match")
      return
    }
    if (newPassword.length < 6) {
      setPasswordError("New password must be at least 6 characters")
      return
    }
    setPasswordLoading(true)
    try {
      await changePassword(currentPassword, newPassword)
      setPasswordSuccess(true)
      setCurrentPassword("")
      setNewPassword("")
      setConfirmPassword("")
    } catch (err) {
      setPasswordError(err instanceof Error ? err.message : "Failed to change password")
    } finally {
      setPasswordLoading(false)
    }
  }

  useEffect(() => {
    if (open && currentUser) {
      setName(currentUser.name)
      setPhone(currentUser.phone ?? "")
    }
  }, [open, currentUser])

  if (!currentUser) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="uppercase tracking-wider">My Profile</DialogTitle>
        </DialogHeader>

        <div className="flex gap-1 border-b border-border mb-4">
          <button
            type="button"
            onClick={() => setActiveTab("profile")}
            className={`px-3 py-1.5 text-xs font-medium tracking-wider uppercase transition-colors ${
              activeTab === "profile"
                ? "text-primary border-b-2 border-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            Profile
          </button>
          <button
            type="button"
            onClick={() => setActiveTab("password")}
            className={`px-3 py-1.5 text-xs font-medium tracking-wider uppercase transition-colors ${
              activeTab === "password"
                ? "text-primary border-b-2 border-primary"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            Change Password
          </button>
        </div>

        {activeTab === "profile" && (
          <form onSubmit={handleProfileSubmit} className="space-y-4">
            {profileError && (
              <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-2 text-sm text-destructive">
                <AlertCircle className="h-4 w-4 shrink-0" />
                {profileError}
              </div>
            )}
            {profileSuccess && (
              <div className="rounded-md border border-primary/30 bg-primary/10 p-2 text-sm text-primary">
                Profile updated successfully.
              </div>
            )}

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Name
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Your name"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Email (read-only)
              </label>
              <input
                type="email"
                value={currentUser.email}
                readOnly
                className="w-full rounded-md border border-border bg-muted/50 p-3 text-sm text-muted-foreground"
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Phone
              </label>
              <input
                type="tel"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="Your phone number"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
              />
            </div>

            <button
              type="submit"
              disabled={profileLoading}
              className="w-full rounded-md border border-primary/30 bg-primary/10 p-3 text-sm font-bold tracking-wider text-primary uppercase transition-all hover:bg-primary/20 disabled:opacity-50"
            >
              {profileLoading ? "Saving..." : "Save Profile"}
            </button>
          </form>
        )}

        {activeTab === "password" && (
          <form onSubmit={handlePasswordSubmit} className="space-y-4">
            {passwordError && (
              <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-2 text-sm text-destructive">
                <AlertCircle className="h-4 w-4 shrink-0" />
                {passwordError}
              </div>
            )}
            {passwordSuccess && (
              <div className="rounded-md border border-primary/30 bg-primary/10 p-2 text-sm text-primary">
                Password changed successfully.
              </div>
            )}

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Current Password
              </label>
              <input
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                placeholder="Enter current password"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                New Password
              </label>
              <input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                placeholder="Enter new password (min 6 characters)"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
                required
                minLength={6}
              />
            </div>

            <div className="space-y-2">
              <label className="text-xs font-bold tracking-wider text-muted-foreground uppercase">
                Confirm New Password
              </label>
              <input
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Confirm new password"
                className="w-full rounded-md border border-border bg-card p-3 text-sm text-foreground placeholder:text-muted-foreground focus:border-primary/50 focus:outline-none"
                required
                minLength={6}
              />
            </div>

            <button
              type="submit"
              disabled={passwordLoading}
              className="w-full rounded-md border border-primary/30 bg-primary/10 p-3 text-sm font-bold tracking-wider text-primary uppercase transition-all hover:bg-primary/20 disabled:opacity-50"
            >
              {passwordLoading ? "Updating..." : "Change Password"}
            </button>
          </form>
        )}
      </DialogContent>
    </Dialog>
  )
}
