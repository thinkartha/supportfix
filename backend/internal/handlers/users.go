package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/supporttickr/backend/internal/middleware"
	"github.com/supporttickr/backend/internal/models"
	"github.com/supporttickr/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Store store.Store
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	orgID := middleware.GetOrgID(r.Context())

	users, err := h.Store.ListUsers(r.Context(), role, orgID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query users")
		return
	}
	if users == nil {
		users = []models.UserResponse{}
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	userIDParam := r.PathValue("id")

	u, err := h.Store.GetUser(r.Context(), userIDParam)
	if err != nil || u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, u.ToResponse())
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		writeError(w, http.StatusForbidden, "admin access required")
		return
	}

	var input struct {
		Name           string  `json:"name"`
		Email          string  `json:"email"`
		Role           string  `json:"role"`
		OrganizationID *string `json:"organizationId"`
		Password       *string `json:"password"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" || input.Email == "" || input.Role == "" {
		writeError(w, http.StatusBadRequest, "name, email, and role are required")
		return
	}

	// Default password
	password := "changeme"
	if input.Password != nil && *input.Password != "" {
		password = *input.Password
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	avatar := initialsFromName(input.Name)
	id := "user-" + generateID()

	u := &models.User{
		ID:             id,
		Name:           input.Name,
		Email:          input.Email,
		PasswordHash:   string(passwordHash),
		Role:           input.Role,
		OrganizationID: input.OrganizationID,
		Avatar:         avatar,
	}
	if err := h.Store.CreateUser(r.Context(), u); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user: "+err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, u.ToResponse())
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		writeError(w, http.StatusForbidden, "admin access required")
		return
	}

	userIDParam := r.PathValue("id")

	var input struct {
		Name           *string `json:"name"`
		Email          *string `json:"email"`
		Role           *string `json:"role"`
		OrganizationID *string `json:"organizationId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var avatar *string
	if input.Name != nil {
		a := initialsFromName(*input.Name)
		avatar = &a
	}

	// Determine organization_id update: clear when role is internal, or when client has none
	orgIDToUpdate := input.OrganizationID
	if input.Role != nil {
		if *input.Role != "client" {
			empty := ""
			orgIDToUpdate = &empty // clear org for admin/support-staff/support-lead
		} else if input.OrganizationID == nil {
			empty := ""
			orgIDToUpdate = &empty // client with no org = remove attribute
		}
	}

	if err := h.Store.UpdateUser(r.Context(), userIDParam, input.Name, input.Email, input.Role, orgIDToUpdate, avatar); err != nil {
		log.Printf("UpdateUser error: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to update user: "+err.Error())
		return
	}

	u, err := h.Store.GetUser(r.Context(), userIDParam)
	if err != nil || u == nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, u.ToResponse())
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	role := middleware.GetRole(r.Context())
	if role != "admin" {
		writeError(w, http.StatusForbidden, "admin access required")
		return
	}

	userIDParam := r.PathValue("id")
	currentUserID := middleware.GetUserID(r.Context())
	if userIDParam == currentUserID {
		writeError(w, http.StatusBadRequest, "cannot delete your own account")
		return
	}

	// Unassign tickets assigned to this user
	tickets, _ := h.Store.ListTickets(r.Context(), "", "", "", "", userIDParam, "")
	empty := ""
	for _, t := range tickets {
		_ = h.Store.UpdateTicket(r.Context(), t.ID, nil, nil, &empty, nil)
	}

	if err := h.Store.DeleteUser(r.Context(), userIDParam); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func initialsFromName(name string) string {
	words := splitWords(name)
	avatar := ""
	for _, w := range words {
		if len(w) > 0 {
			avatar += string([]rune(w)[0])
		}
	}
	if len(avatar) > 2 {
		avatar = avatar[:2]
	}
	return strings.ToUpper(avatar)
}

func splitWords(s string) []string {
	words := []string{}
	current := ""
	for _, c := range s {
		if c == ' ' || c == '\t' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
