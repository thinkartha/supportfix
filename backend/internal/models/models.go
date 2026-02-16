package models

import "time"

// Organization represents a client organization
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Plan         string    `json:"plan"`
	ContactEmail string    `json:"contactEmail"`
	CreatedAt    time.Time `json:"createdAt"`
}

// User represents a system user
type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           string    `json:"role"`
	OrganizationID *string   `json:"organizationId"`
	Avatar         string    `json:"avatar"`
	Phone          string    `json:"phone"`
	CreatedAt      time.Time `json:"createdAt"`
}

// UserResponse is the JSON-safe version of User
type UserResponse struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	OrganizationID *string `json:"organizationId"`
	Avatar         string  `json:"avatar"`
	Phone          string  `json:"phone"`
}

func (u *User) ToResponse() UserResponse {
	r := UserResponse{
		ID:             u.ID,
		Name:           u.Name,
		Email:          u.Email,
		Role:           u.Role,
		Avatar:         u.Avatar,
		OrganizationID: u.OrganizationID,
		Phone:          u.Phone,
	}
	return r
}

// Ticket represents a support ticket
type Ticket struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	Priority       string    `json:"priority"`
	Category       string    `json:"category"`
	OrganizationID string    `json:"organizationId"`
	CreatedBy      string    `json:"createdBy"`
	AssignedTo     *string   `json:"assignedTo"`
	HoursWorked    float64   `json:"hoursWorked"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TicketResponse is the JSON-safe version
type TicketResponse struct {
	ID                string              `json:"id"`
	Title             string              `json:"title"`
	Description       string              `json:"description"`
	Status            string              `json:"status"`
	Priority          string              `json:"priority"`
	Category          string              `json:"category"`
	OrganizationID    string              `json:"organizationId"`
	CreatedBy         string              `json:"createdBy"`
	AssignedTo        *string             `json:"assignedTo"`
	HoursWorked       float64             `json:"hoursWorked"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
	Messages          []Message           `json:"messages"`
	TimeEntries       []TimeEntry         `json:"timeEntries"`
	ConversionRequest *ConversionRequest  `json:"conversionRequest,omitempty"`
}

func (t *Ticket) ToResponse() TicketResponse {
	r := TicketResponse{
		ID:             t.ID,
		Title:          t.Title,
		Description:    t.Description,
		Status:         t.Status,
		Priority:       t.Priority,
		Category:       t.Category,
		OrganizationID: t.OrganizationID,
		CreatedBy:      t.CreatedBy,
		AssignedTo:     t.AssignedTo,
		HoursWorked:    t.HoursWorked,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
		Messages:       []Message{},
		TimeEntries:    []TimeEntry{},
	}
	return r
}

// Message represents a ticket message
type Message struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticketId"`
	UserID     string    `json:"userId"`
	Content    string    `json:"content"`
	IsInternal bool      `json:"isInternal"`
	CreatedAt  time.Time `json:"createdAt"`
}

// TimeEntry represents logged time on a ticket
type TimeEntry struct {
	ID          string    `json:"id"`
	TicketID    string    `json:"ticketId"`
	UserID      string    `json:"userId"`
	Hours       float64   `json:"hours"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ConversionRequest represents a ticket-to-feature conversion request
type ConversionRequest struct {
	ID               string    `json:"id"`
	TicketID         string    `json:"ticketId"`
	ProposedType     string    `json:"proposedType"`
	Reason           string    `json:"reason"`
	InternalApproval string    `json:"internalApproval"`
	ClientApproval   string    `json:"clientApproval"`
	ProposedBy       string    `json:"proposedBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Month          int       `json:"month"`
	Year           int       `json:"year"`
	TicketsClosed  int       `json:"ticketsClosed"`
	TotalHours     float64   `json:"totalHours"`
	RatePerHour    float64   `json:"ratePerHour"`
	TotalAmount    float64   `json:"totalAmount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}

// ActivityItem represents an activity log entry
type ActivityItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      string    `json:"userId"`
	TicketID    *string   `json:"ticketId,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ActivityResponse is the JSON-safe version
type ActivityResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserID      string    `json:"userId"`
	TicketID    *string   `json:"ticketId,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (a *ActivityItem) ToResponse() ActivityResponse {
	r := ActivityResponse{
		ID:          a.ID,
		Type:        a.Type,
		Description: a.Description,
		UserID:      a.UserID,
		TicketID:    a.TicketID,
		CreatedAt:   a.CreatedAt,
	}
	return r
}

// Request/response types for API

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type CreateTicketRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Priority       string `json:"priority"`
	Category       string `json:"category"`
	OrganizationID string `json:"organizationId"`
}

type UpdateTicketRequest struct {
	Status   *string `json:"status,omitempty"`
	Priority *string `json:"priority,omitempty"`
	AssignTo *string `json:"assignedTo,omitempty"`
}

type CreateMessageRequest struct {
	Content    string `json:"content"`
	IsInternal bool   `json:"isInternal"`
}

type CreateTimeEntryRequest struct {
	Hours       float64 `json:"hours"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type ConversionRequestBody struct {
	ProposedType string `json:"proposedType"`
	Reason       string `json:"reason"`
}

type UpdateApprovalRequest struct {
	Side   string `json:"side"`   // "internal" or "client"
	Status string `json:"status"` // "approved" or "rejected"
}

type CreateInvoiceRequest struct {
	OrganizationID string  `json:"organizationId"`
	Month          int     `json:"month"`
	Year           int     `json:"year"`
	TicketsClosed  int     `json:"ticketsClosed"`
	TotalHours     float64 `json:"totalHours"`
	RatePerHour    float64 `json:"ratePerHour"`
	TotalAmount    float64 `json:"totalAmount"`
}

type UpdateInvoiceStatusRequest struct {
	Status string `json:"status"`
}

type DashboardStats struct {
	TotalTickets    int     `json:"totalTickets"`
	OpenTickets     int     `json:"openTickets"`
	InProgress      int     `json:"inProgress"`
	Resolved        int     `json:"resolved"`
	Closed          int     `json:"closed"`
	AvgResponseTime string  `json:"avgResponseTime"`
	TotalHours      float64 `json:"totalHours"`
	PendingApproval int     `json:"pendingApprovals"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type UpdateMyProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
