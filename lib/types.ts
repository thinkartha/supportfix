export type UserRole = "admin" | "support-lead" | "support-staff" | "client"

export type TicketStatus = "open" | "in-progress" | "awaiting-client" | "resolved" | "closed"

export type TicketPriority = "low" | "medium" | "high" | "critical"

export type TicketCategory = "bug" | "support" | "question" | "feature" | "enhancement"

export type ApprovalStatus = "pending" | "approved" | "rejected"

export type InvoiceStatus = "draft" | "sent" | "paid"

export interface Organization {
  id: string
  name: string
  plan: "starter" | "professional" | "enterprise"
  contactEmail: string
  createdAt: string
}

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  organizationId: string | null
  avatar: string
  phone?: string
}

export interface Message {
  id: string
  ticketId: string
  userId: string
  content: string
  createdAt: string
  isInternal: boolean
}

export interface TimeEntry {
  id: string
  ticketId: string
  userId: string
  hours: number
  description: string
  date: string
}

export interface ConversionRequest {
  id: string
  ticketId: string
  proposedType: "feature" | "enhancement"
  reason: string
  internalApproval: ApprovalStatus
  clientApproval: ApprovalStatus
  proposedBy: string
  createdAt: string
}

export interface Ticket {
  id: string
  title: string
  description: string
  status: TicketStatus
  priority: TicketPriority
  category: TicketCategory
  organizationId: string
  createdBy: string
  assignedTo: string | null
  createdAt: string
  updatedAt: string
  hoursWorked: number
  messages: Message[]
  timeEntries: TimeEntry[]
  conversionRequest?: ConversionRequest
}

export interface Invoice {
  id: string
  organizationId: string
  month: number
  year: number
  ticketsClosed: number
  totalHours: number
  ratePerHour: number
  totalAmount: number
  status: InvoiceStatus
  createdAt: string
}

export interface ActivityItem {
  id: string
  type: "ticket-created" | "ticket-updated" | "message-added" | "ticket-resolved" | "conversion-requested" | "conversion-approved"
  description: string
  userId: string
  ticketId?: string
  createdAt: string
}
