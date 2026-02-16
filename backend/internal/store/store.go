package store

import (
	"context"

	"github.com/supporttickr/backend/internal/models"
)

// Store is the data access interface (DynamoDB).
type Store interface {
	// Users
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context, role, orgID string) ([]models.UserResponse, error)
	CreateUser(ctx context.Context, u *models.User) error
	UpdateUser(ctx context.Context, id string, name, email, role *string, orgID *string, avatar *string) error
	UpdateMyProfile(ctx context.Context, id string, name, phone *string) error
	UpdatePassword(ctx context.Context, id string, newHash string) error
	DeleteUser(ctx context.Context, id string) error

	// Organizations
	ListOrgs(ctx context.Context, role, orgID string) ([]models.Organization, error)
	GetOrg(ctx context.Context, id string) (*models.Organization, error)
	CreateOrg(ctx context.Context, o *models.Organization) error
	UpdateOrg(ctx context.Context, id string, name, plan, contactEmail *string) error
	DeleteOrg(ctx context.Context, id string) error

	// Tickets
	ListTickets(ctx context.Context, status, priority, category, organizationID, assignedTo, search string) ([]models.Ticket, error)
	GetTicket(ctx context.Context, id string) (*models.Ticket, error)
	CreateTicket(ctx context.Context, t *models.Ticket) error
	UpdateTicket(ctx context.Context, id string, status, priority, assignedTo *string, hoursWorked *float64) error
	UpdateTicketCategory(ctx context.Context, id, category string) error

	// Messages
	GetMessagesByTicketID(ctx context.Context, ticketID string) ([]models.Message, error)
	AddMessage(ctx context.Context, m *models.Message) error

	// Time entries
	GetTimeEntriesByTicketID(ctx context.Context, ticketID string) ([]models.TimeEntry, error)
	AddTimeEntry(ctx context.Context, te *models.TimeEntry) error

	// Conversion requests
	GetConversionByTicketID(ctx context.Context, ticketID string) (*models.ConversionRequest, error)
	GetConversionByID(ctx context.Context, id string) (*models.ConversionRequest, error)
	CreateConversionRequest(ctx context.Context, cr *models.ConversionRequest) error
	UpdateConversionRequest(ctx context.Context, id string, internalApproval, clientApproval *string) error
	ListConversionRequestsPending(ctx context.Context) ([]models.ConversionRequest, error)

	// Invoices
	ListInvoices(ctx context.Context, role, orgID string) ([]models.Invoice, error)
	CreateInvoice(ctx context.Context, inv *models.Invoice) error
	UpdateInvoiceStatus(ctx context.Context, id, status string) error

	// Activities
	ListActivities(ctx context.Context, limit int) ([]models.ActivityItem, error)
	CreateActivity(ctx context.Context, a *models.ActivityItem) error
}
