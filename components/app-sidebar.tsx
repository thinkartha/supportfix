"use client"

import { useState } from "react"
import { useStore } from "@/lib/store"
import { ProfileModal } from "@/components/profile-modal"
import type { UserRole } from "@/lib/types"
import {
  LayoutDashboard,
  Ticket,
  Building2,
  Users,
  GitPullRequestArrow,
  Receipt,
  Settings,
  LogOut,
  Headphones,
  ChevronRight,
} from "lucide-react"

export type Section =
  | "dashboard"
  | "tickets"
  | "organizations"
  | "users"
  | "approvals"
  | "billing"
  | "settings"

interface NavItem {
  id: Section
  label: string
  icon: React.ReactNode
  roles: UserRole[]
}

const navItems: NavItem[] = [
  {
    id: "dashboard",
    label: "OVERVIEW",
    icon: <LayoutDashboard className="h-4 w-4" />,
    roles: ["admin", "support-lead", "support-staff", "client"],
  },
  {
    id: "tickets",
    label: "TICKETS",
    icon: <Ticket className="h-4 w-4" />,
    roles: ["admin", "support-lead", "support-staff", "client"],
  },
  {
    id: "organizations",
    label: "ORGANIZATIONS",
    icon: <Building2 className="h-4 w-4" />,
    roles: ["admin", "support-lead"],
  },
  {
    id: "users",
    label: "USERS",
    icon: <Users className="h-4 w-4" />,
    roles: ["admin"],
  },
  {
    id: "approvals",
    label: "APPROVALS",
    icon: <GitPullRequestArrow className="h-4 w-4" />,
    roles: ["admin", "support-lead", "client"],
  },
  {
    id: "billing",
    label: "BILLING",
    icon: <Receipt className="h-4 w-4" />,
    roles: ["admin", "client"],
  },
  {
    id: "settings",
    label: "SETTINGS",
    icon: <Settings className="h-4 w-4" />,
    roles: ["admin"],
  },
]

interface AppSidebarProps {
  activeSection: Section
  onSectionChange: (section: Section) => void
}

export function AppSidebar({ activeSection, onSectionChange }: AppSidebarProps) {
  const { currentUser, logout, getUserById, getOrgById, tickets } = useStore()
  const [profileOpen, setProfileOpen] = useState(false)

  if (!currentUser) return null

  const filteredItems = navItems.filter((item) =>
    item.roles.includes(currentUser.role)
  )

  const orgName = currentUser.organizationId
    ? getOrgById(currentUser.organizationId)?.name
    : "Internal Team"

  const openTicketCount = tickets.filter((t) => {
    if (currentUser.role === "client") {
      return (
        t.organizationId === currentUser.organizationId &&
        t.status !== "closed" &&
        t.status !== "resolved"
      )
    }
    return t.status !== "closed" && t.status !== "resolved"
  }).length

  const pendingApprovals = tickets.filter((t) => {
    if (!t.conversionRequest) return false
    if (currentUser.role === "client") {
      return (
        t.organizationId === currentUser.organizationId &&
        t.conversionRequest.clientApproval === "pending"
      )
    }
    if (currentUser.role === "admin" || currentUser.role === "support-lead") {
      return t.conversionRequest.internalApproval === "pending"
    }
    return false
  }).length

  const getBadge = (id: Section) => {
    if (id === "tickets" && openTicketCount > 0) return openTicketCount
    if (id === "approvals" && pendingApprovals > 0) return pendingApprovals
    return null
  }

  return (
    <aside className="flex h-screen w-60 shrink-0 flex-col border-r border-border bg-card">
      {/* Logo */}
      <div className="flex items-center gap-2.5 border-b border-border px-4 py-4">
        <div className="flex h-8 w-8 items-center justify-center rounded border border-primary/30 bg-primary/10">
          <Headphones className="h-4 w-4 text-primary" />
        </div>
        <span className="text-sm font-bold tracking-wider text-foreground uppercase">
          SupportFIX
        </span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-0.5 overflow-y-auto px-2 py-3">
        {filteredItems.map((item) => {
          const isActive = activeSection === item.id
          const badge = getBadge(item.id)
          return (
            <button
              key={item.id}
              onClick={() => onSectionChange(item.id)}
              className={`flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-left text-xs font-medium tracking-wider transition-all ${
                isActive
                  ? "bg-primary/10 text-primary border border-primary/20"
                  : "text-muted-foreground hover:bg-secondary hover:text-foreground border border-transparent"
              }`}
            >
              <span className={isActive ? "text-primary" : ""}>{item.icon}</span>
              <span className="flex-1">{item.label}</span>
              {badge !== null && (
                <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-primary/20 px-1.5 text-[10px] font-bold text-primary">
                  {badge}
                </span>
              )}
              {isActive && <ChevronRight className="h-3 w-3 text-primary" />}
            </button>
          )
        })}
      </nav>

      {/* User section */}
      <div className="border-t border-border p-3">
        <button
          type="button"
          onClick={() => setProfileOpen(true)}
          className="mb-2 flex w-full items-center gap-2.5 rounded-md bg-secondary/50 px-3 py-2.5 text-left transition-colors hover:bg-secondary"
        >
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded border border-border bg-card text-[10px] font-bold tracking-wider text-foreground">
            {currentUser.avatar}
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-xs font-medium text-foreground truncate">
              {currentUser.name}
            </p>
            <p className="text-[10px] text-muted-foreground truncate">{orgName}</p>
          </div>
        </button>
        <ProfileModal open={profileOpen} onOpenChange={setProfileOpen} />
        <button
          onClick={logout}
          className="flex w-full items-center gap-2 rounded-md px-3 py-1.5 text-[11px] font-medium tracking-wider text-muted-foreground transition-colors hover:bg-secondary hover:text-destructive"
        >
          <LogOut className="h-3.5 w-3.5" />
          SWITCH USER
        </button>
      </div>
    </aside>
  )
}
