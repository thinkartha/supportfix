const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

let authToken: string | null = null

export function setToken(token: string | null) {
  authToken = token
  if (token) {
    if (typeof window !== "undefined") localStorage.setItem("st_token", token)
  } else {
    if (typeof window !== "undefined") localStorage.removeItem("st_token")
  }
}

export function getToken(): string | null {
  if (authToken) return authToken
  if (typeof window !== "undefined") {
    authToken = localStorage.getItem("st_token")
  }
  return authToken
}

async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken()
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> || {}),
  }
  if (token) {
    headers["Authorization"] = `Bearer ${token}`
  }

  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || `API error: ${res.status}`)
  }

  return res.json()
}

// ============================================================================
// Auth
// ============================================================================

export interface LoginResult {
  token: string
  user: import("./types").User
}

export async function login(email: string, password: string): Promise<LoginResult> {
  const result = await apiFetch<LoginResult>("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  })
  setToken(result.token)
  return result
}

export async function getMe(): Promise<import("./types").User> {
  return apiFetch("/api/auth/me")
}

export async function changePassword(currentPassword: string, newPassword: string): Promise<{ status: string }> {
  return apiFetch("/api/auth/change-password", {
    method: "PUT",
    body: JSON.stringify({ currentPassword, newPassword }),
  })
}

export async function updateMyProfile(data: { name?: string; phone?: string }): Promise<import("./types").User> {
  return apiFetch("/api/auth/me", {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

export async function forgotPassword(email: string): Promise<{ message: string }> {
  return apiFetch("/api/auth/forgot-password", {
    method: "POST",
    body: JSON.stringify({ email }),
  })
}

export function logout() {
  setToken(null)
}

// ============================================================================
// Users
// ============================================================================

export async function getUsers(): Promise<import("./types").User[]> {
  return apiFetch("/api/users")
}

export async function getUser(id: string): Promise<import("./types").User> {
  return apiFetch(`/api/users/${id}`)
}

export async function createUser(data: {
  name: string
  email: string
  role: string
  organizationId: string | null
  password?: string
}): Promise<import("./types").User> {
  return apiFetch("/api/users", {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function updateUser(id: string, data: {
  name?: string
  email?: string
  role?: string
  organizationId?: string | null
}): Promise<import("./types").User> {
  return apiFetch(`/api/users/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

export async function deleteUser(id: string): Promise<void> {
  await apiFetch(`/api/users/${id}`, { method: "DELETE" })
}

// ============================================================================
// Organizations
// ============================================================================

export async function getOrganizations(): Promise<import("./types").Organization[]> {
  return apiFetch("/api/organizations")
}

export async function getOrganization(id: string): Promise<import("./types").Organization> {
  return apiFetch(`/api/organizations/${id}`)
}

export async function createOrganization(data: {
  name: string
  plan: string
  contactEmail: string
}): Promise<import("./types").Organization> {
  return apiFetch("/api/organizations", {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function updateOrganization(id: string, data: {
  name?: string
  plan?: string
  contactEmail?: string
}): Promise<import("./types").Organization> {
  return apiFetch(`/api/organizations/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

export async function deleteOrganization(id: string): Promise<void> {
  await apiFetch(`/api/organizations/${id}`, { method: "DELETE" })
}

// ============================================================================
// Tickets
// ============================================================================

export async function getTickets(params?: Record<string, string>): Promise<import("./types").Ticket[]> {
  const qs = params ? "?" + new URLSearchParams(params).toString() : ""
  return apiFetch(`/api/tickets${qs}`)
}

export async function getTicket(id: string): Promise<import("./types").Ticket> {
  return apiFetch(`/api/tickets/${id}`)
}

export async function createTicket(data: {
  title: string
  description: string
  priority: string
  category: string
  organizationId: string
}): Promise<import("./types").Ticket> {
  return apiFetch("/api/tickets", {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function updateTicket(id: string, data: {
  status?: string
  priority?: string
  assignedTo?: string
}): Promise<import("./types").Ticket> {
  return apiFetch(`/api/tickets/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

export async function addMessage(ticketId: string, data: {
  content: string
  isInternal: boolean
}): Promise<{ id: string }> {
  return apiFetch(`/api/tickets/${ticketId}/messages`, {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function addTimeEntry(ticketId: string, data: {
  hours: number
  description: string
  date: string
}): Promise<{ id: string }> {
  return apiFetch(`/api/tickets/${ticketId}/time-entries`, {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function requestConversion(ticketId: string, data: {
  proposedType: string
  reason: string
}): Promise<{ id: string }> {
  return apiFetch(`/api/tickets/${ticketId}/convert`, {
    method: "POST",
    body: JSON.stringify(data),
  })
}

// ============================================================================
// Approvals
// ============================================================================

export async function getApprovals(): Promise<import("./types").ConversionRequest[]> {
  return apiFetch("/api/approvals")
}

export async function updateApproval(id: string, data: {
  side: "internal" | "client"
  status: "approved" | "rejected"
}): Promise<{ status: string }> {
  return apiFetch(`/api/approvals/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
}

// ============================================================================
// Invoices
// ============================================================================

export async function getInvoices(): Promise<import("./types").Invoice[]> {
  return apiFetch("/api/invoices")
}

export async function createInvoice(data: {
  organizationId: string
  month: number
  year: number
  ticketsClosed: number
  totalHours: number
  ratePerHour: number
  totalAmount: number
}): Promise<{ id: string }> {
  return apiFetch("/api/invoices", {
    method: "POST",
    body: JSON.stringify(data),
  })
}

export async function updateInvoiceStatus(id: string, status: string): Promise<{ status: string }> {
  return apiFetch(`/api/invoices/${id}`, {
    method: "PUT",
    body: JSON.stringify({ status }),
  })
}

// ============================================================================
// Dashboard
// ============================================================================

export interface DashboardStats {
  totalTickets: number
  openTickets: number
  inProgress: number
  resolved: number
  closed: number
  avgResponseTime: string
  totalHours: number
  pendingApprovals: number
}

export async function getDashboardStats(): Promise<DashboardStats> {
  return apiFetch("/api/dashboard/stats")
}

export async function getActivities(): Promise<import("./types").ActivityItem[]> {
  return apiFetch("/api/dashboard/activities")
}
