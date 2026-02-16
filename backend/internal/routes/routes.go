package routes

import (
	"net/http"

	"github.com/supporttickr/backend/internal/config"
	"github.com/supporttickr/backend/internal/handlers"
	"github.com/supporttickr/backend/internal/middleware"
	"github.com/supporttickr/backend/internal/store"
)

func Setup(st store.Store, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Initialize handlers
	authH := &handlers.AuthHandler{Store: st, JWTSecret: cfg.JWTSecret}
	ticketH := &handlers.TicketHandler{Store: st}
	orgH := &handlers.OrgHandler{Store: st}
	userH := &handlers.UserHandler{Store: st}
	approvalH := &handlers.ApprovalHandler{Store: st}
	invoiceH := &handlers.InvoiceHandler{Store: st}
	dashboardH := &handlers.DashboardHandler{Store: st}

	// Auth middleware
	authMW := middleware.Auth(cfg.JWTSecret)

	// Public routes
	mux.HandleFunc("POST /api/auth/login", authH.Login)
	mux.HandleFunc("POST /api/auth/forgot-password", authH.ForgotPassword)

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Protected routes - wrap with auth middleware
	mux.Handle("GET /api/auth/me", authMW(http.HandlerFunc(authH.Me)))
	mux.Handle("PUT /api/auth/change-password", authMW(http.HandlerFunc(authH.ChangePassword)))
	mux.Handle("PUT /api/auth/me", authMW(http.HandlerFunc(authH.UpdateMyProfile)))

	mux.Handle("GET /api/users", authMW(http.HandlerFunc(userH.List)))
	mux.Handle("GET /api/users/{id}", authMW(http.HandlerFunc(userH.Get)))
	mux.Handle("POST /api/users", authMW(http.HandlerFunc(userH.Create)))
	mux.Handle("PUT /api/users/{id}", authMW(http.HandlerFunc(userH.Update)))
	mux.Handle("DELETE /api/users/{id}", authMW(http.HandlerFunc(userH.Delete)))

	mux.Handle("GET /api/organizations", authMW(http.HandlerFunc(orgH.List)))
	mux.Handle("GET /api/organizations/{id}", authMW(http.HandlerFunc(orgH.Get)))
	mux.Handle("POST /api/organizations", authMW(http.HandlerFunc(orgH.Create)))
	mux.Handle("PUT /api/organizations/{id}", authMW(http.HandlerFunc(orgH.Update)))
	mux.Handle("DELETE /api/organizations/{id}", authMW(http.HandlerFunc(orgH.Delete)))

	mux.Handle("GET /api/tickets", authMW(http.HandlerFunc(ticketH.List)))
	mux.Handle("GET /api/tickets/{id}", authMW(http.HandlerFunc(ticketH.Get)))
	mux.Handle("POST /api/tickets", authMW(http.HandlerFunc(ticketH.Create)))
	mux.Handle("PUT /api/tickets/{id}", authMW(http.HandlerFunc(ticketH.Update)))
	mux.Handle("POST /api/tickets/{id}/messages", authMW(http.HandlerFunc(ticketH.AddMessage)))
	mux.Handle("POST /api/tickets/{id}/time-entries", authMW(http.HandlerFunc(ticketH.AddTimeEntry)))
	mux.Handle("POST /api/tickets/{id}/convert", authMW(http.HandlerFunc(ticketH.RequestConversion)))

	mux.Handle("GET /api/approvals", authMW(http.HandlerFunc(approvalH.List)))
	mux.Handle("PUT /api/approvals/{id}", authMW(http.HandlerFunc(approvalH.Update)))

	mux.Handle("GET /api/invoices", authMW(http.HandlerFunc(invoiceH.List)))
	mux.Handle("POST /api/invoices", authMW(http.HandlerFunc(invoiceH.Create)))
	mux.Handle("PUT /api/invoices/{id}", authMW(http.HandlerFunc(invoiceH.UpdateStatus)))

	mux.Handle("GET /api/dashboard/stats", authMW(http.HandlerFunc(dashboardH.Stats)))
	mux.Handle("GET /api/dashboard/activities", authMW(http.HandlerFunc(dashboardH.Activities)))

	corsHandler := middleware.CORS(cfg.FrontendURL)(mux)
	return corsHandler
}
