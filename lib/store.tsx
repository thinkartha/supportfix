"use client"

import React, { createContext, useContext, useState, useCallback, useEffect } from "react"
import type {
  Organization,
  User,
  Ticket,
  Invoice,
  ActivityItem,
  Message,
  TimeEntry,
  ConversionRequest,
  TicketStatus,
  TicketPriority,
  TicketCategory,
  ApprovalStatus,
  InvoiceStatus,
} from "./types"
import * as api from "./api"

interface StoreContextType {
  organizations: Organization[]
  users: User[]
  tickets: Ticket[]
  invoices: Invoice[]
  activities: ActivityItem[]
  currentUser: User | null
  ratePerHour: number
  loading: boolean

  loginWithCredentials: (email: string, password: string) => Promise<void>
  logout: () => void
  changePassword: (currentPassword: string, newPassword: string) => Promise<void>
  updateMyProfile: (data: { name?: string; phone?: string }) => Promise<void>
  forgotPassword: (email: string) => Promise<void>

  createTicket: (ticket: Omit<Ticket, "id" | "createdAt" | "updatedAt" | "hoursWorked" | "messages" | "timeEntries">) => Promise<void>
  updateTicketStatus: (ticketId: string, status: TicketStatus) => Promise<void>
  updateTicketPriority: (ticketId: string, priority: TicketPriority) => Promise<void>
  assignTicket: (ticketId: string, userId: string | null) => Promise<void>
  addMessage: (ticketId: string, message: Omit<Message, "id" | "createdAt">) => Promise<void>
  addTimeEntry: (ticketId: string, entry: Omit<TimeEntry, "id">) => Promise<void>

  requestConversion: (ticketId: string, request: Omit<ConversionRequest, "id" | "createdAt" | "internalApproval" | "clientApproval">) => Promise<void>
  updateConversionApproval: (ticketId: string, side: "internal" | "client", status: ApprovalStatus) => Promise<void>

  createInvoice: (invoice: Omit<Invoice, "id" | "createdAt">) => Promise<void>
  updateInvoiceStatus: (invoiceId: string, status: InvoiceStatus) => Promise<void>

  createOrganization: (org: { name: string; plan: "starter" | "professional" | "enterprise"; contactEmail: string }) => Promise<void>
  updateOrganization: (id: string, data: Partial<Pick<Organization, "name" | "plan" | "contactEmail">>) => Promise<void>
  deleteOrganization: (id: string) => Promise<void>

  createUser: (user: { name: string; email: string; role: import("./types").UserRole; organizationId: string | null; password?: string }) => Promise<void>
  updateUser: (id: string, data: Partial<Pick<User, "name" | "email" | "role" | "organizationId">>) => Promise<void>
  deleteUser: (id: string) => Promise<void>

  setRatePerHour: (rate: number) => void
  getUserById: (id: string) => User | undefined
  getOrgById: (id: string) => Organization | undefined
  getTicketsForOrg: (orgId: string) => Ticket[]
  getTicketsForUser: (userId: string) => Ticket[]
  refreshData: () => Promise<void>
}

const StoreContext = createContext<StoreContextType | null>(null)

export function StoreProvider({ children }: { children: React.ReactNode }) {
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [tickets, setTickets] = useState<Ticket[]>([])
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [activities, setActivities] = useState<ActivityItem[]>([])
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [ratePerHour, setRatePerHour] = useState(150)
  const [loading, setLoading] = useState(false)

  // Restore session from stored token on mount
  useEffect(() => {
    const token = api.getToken()
    if (token) {
      api
        .getMe()
        .then((user) => {
          setCurrentUser(user)
        })
        .catch(() => {
          api.logout()
        })
    }
  }, [])

  const refreshData = useCallback(async () => {
    setLoading(true)
    try {
      const [orgs, usrs, tkts, invs, acts] = await Promise.all([
        api.getOrganizations(),
        api.getUsers(),
        api.getTickets(),
        api.getInvoices(),
        api.getActivities(),
      ])
      setOrganizations(orgs)
      setUsers(usrs)
      setTickets(tkts)
      setInvoices(invs)
      setActivities(acts)
    } catch (err) {
      console.error("Failed to refresh data:", err)
    } finally {
      setLoading(false)
    }
  }, [])

  // Load data when user is authenticated
  useEffect(() => {
    if (currentUser) {
      refreshData()
    }
  }, [currentUser, refreshData])

  const loginWithCredentials = useCallback(
    async (email: string, password: string) => {
      const result = await api.login(email, password)
      setCurrentUser(result.user)
    },
    []
  )

  const logout = useCallback(() => {
    setCurrentUser(null)
    setOrganizations([])
    setUsers([])
    setTickets([])
    setInvoices([])
    setActivities([])
    api.logout()
  }, [])

  const changePassword = useCallback(
    async (currentPassword: string, newPassword: string) => {
      await api.changePassword(currentPassword, newPassword)
    },
    []
  )

  const updateMyProfile = useCallback(
    async (data: { name?: string; phone?: string }) => {
      const user = await api.updateMyProfile(data)
      setCurrentUser(user)
      setUsers((prev) => {
        const idx = prev.findIndex((u) => u.id === user.id)
        if (idx < 0) return prev
        const next = [...prev]
        next[idx] = user
        return next
      })
    },
    []
  )

  const forgotPassword = useCallback(async (email: string) => {
    await api.forgotPassword(email)
  }, [])

  // Tickets
  const createTicket = useCallback(
    async (ticket: Omit<Ticket, "id" | "createdAt" | "updatedAt" | "hoursWorked" | "messages" | "timeEntries">) => {
      await api.createTicket({
        title: ticket.title,
        description: ticket.description,
        priority: ticket.priority,
        category: ticket.category,
        organizationId: ticket.organizationId,
      })
      await refreshData()
    },
    [refreshData]
  )

  const updateTicketStatus = useCallback(
    async (ticketId: string, status: TicketStatus) => {
      await api.updateTicket(ticketId, { status })
      await refreshData()
    },
    [refreshData]
  )

  const updateTicketPriority = useCallback(
    async (ticketId: string, priority: TicketPriority) => {
      await api.updateTicket(ticketId, { priority })
      await refreshData()
    },
    [refreshData]
  )

  const assignTicket = useCallback(
    async (ticketId: string, userId: string | null) => {
      await api.updateTicket(ticketId, { assignedTo: userId || "" })
      await refreshData()
    },
    [refreshData]
  )

  const addMessage = useCallback(
    async (ticketId: string, message: Omit<Message, "id" | "createdAt">) => {
      await api.addMessage(ticketId, {
        content: message.content,
        isInternal: message.isInternal,
      })
      await refreshData()
    },
    [refreshData]
  )

  const addTimeEntry = useCallback(
    async (ticketId: string, entry: Omit<TimeEntry, "id">) => {
      await api.addTimeEntry(ticketId, {
        hours: entry.hours,
        description: entry.description,
        date: entry.date,
      })
      await refreshData()
    },
    [refreshData]
  )

  // Conversions
  const requestConversion = useCallback(
    async (ticketId: string, request: Omit<ConversionRequest, "id" | "createdAt" | "internalApproval" | "clientApproval">) => {
      await api.requestConversion(ticketId, {
        proposedType: request.proposedType,
        reason: request.reason,
      })
      await refreshData()
    },
    [refreshData]
  )

  const updateConversionApproval = useCallback(
    async (ticketId: string, side: "internal" | "client", status: ApprovalStatus) => {
      const ticket = tickets.find((t) => t.id === ticketId)
      if (ticket?.conversionRequest) {
        await api.updateApproval(ticket.conversionRequest.id, { side, status })
        await refreshData()
      }
    },
    [tickets, refreshData]
  )

  // Invoices
  const createInvoice = useCallback(
    async (invoice: Omit<Invoice, "id" | "createdAt">) => {
      await api.createInvoice({
        organizationId: invoice.organizationId,
        month: invoice.month,
        year: invoice.year,
        ticketsClosed: invoice.ticketsClosed,
        totalHours: invoice.totalHours,
        ratePerHour: invoice.ratePerHour,
        totalAmount: invoice.totalAmount,
      })
      await refreshData()
    },
    [refreshData]
  )

  const updateInvoiceStatus = useCallback(
    async (invoiceId: string, status: InvoiceStatus) => {
      await api.updateInvoiceStatus(invoiceId, status)
      await refreshData()
    },
    [refreshData]
  )

  // Organization CRUD
  const createOrganization = useCallback(
    async (org: { name: string; plan: "starter" | "professional" | "enterprise"; contactEmail: string }) => {
      await api.createOrganization(org)
      await refreshData()
    },
    [refreshData]
  )

  const updateOrganization = useCallback(
    async (id: string, data: Partial<Pick<Organization, "name" | "plan" | "contactEmail">>) => {
      await api.updateOrganization(id, data)
      await refreshData()
    },
    [refreshData]
  )

  const deleteOrganization = useCallback(
    async (id: string) => {
      await api.deleteOrganization(id)
      await refreshData()
    },
    [refreshData]
  )

  // User CRUD
  const createUser = useCallback(
    async (user: { name: string; email: string; role: import("./types").UserRole; organizationId: string | null; password?: string }) => {
      await api.createUser(user)
      await refreshData()
    },
    [refreshData]
  )

  const updateUser = useCallback(
    async (id: string, data: Partial<Pick<User, "name" | "email" | "role" | "organizationId">>) => {
      await api.updateUser(id, data)
      await refreshData()
    },
    [refreshData]
  )

  const deleteUser = useCallback(
    async (id: string) => {
      await api.deleteUser(id)
      await refreshData()
    },
    [refreshData]
  )

  // Helpers
  const getUserById = useCallback((id: string) => users.find((u) => u.id === id), [users])
  const getOrgById = useCallback((id: string) => organizations.find((o) => o.id === id), [organizations])
  const getTicketsForOrg = useCallback((orgId: string) => tickets.filter((t) => t.organizationId === orgId), [tickets])
  const getTicketsForUser = useCallback((userId: string) => tickets.filter((t) => t.assignedTo === userId), [tickets])

  return (
    <StoreContext.Provider
      value={{
        organizations,
        users,
        tickets,
        invoices,
        activities,
        currentUser,
        ratePerHour,
        loading,
        loginWithCredentials,
        logout,
        changePassword,
        updateMyProfile,
        forgotPassword,
        createTicket,
        updateTicketStatus,
        updateTicketPriority,
        assignTicket,
        addMessage,
        addTimeEntry,
        requestConversion,
        updateConversionApproval,
        createInvoice,
        updateInvoiceStatus,
        createOrganization,
        updateOrganization,
        deleteOrganization,
        createUser,
        updateUser,
        deleteUser,
        setRatePerHour,
        getUserById,
        getOrgById,
        getTicketsForOrg,
        getTicketsForUser,
        refreshData,
      }}
    >
      {children}
    </StoreContext.Provider>
  )
}

export function useStore() {
  const context = useContext(StoreContext)
  if (!context) {
    throw new Error("useStore must be used within a StoreProvider")
  }
  return context
}
