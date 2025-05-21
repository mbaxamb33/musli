package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// createContactRequest represents the request body for creating a contact
type createContactRequest struct {
	CompanyID int32  `json:"company_id" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Position  string `json:"position" binding:"omitempty"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty"`
	Notes     string `json:"notes" binding:"omitempty"`
}

// contactResponse represents the API response structure for contact data
type contactResponse struct {
	ContactID int32     `json:"contact_id"`
	CompanyID int32     `json:"company_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Position  string    `json:"position,omitempty"`
	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// updateContactRequest represents the request body for updating a contact
type updateContactRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Position  string `json:"position" binding:"omitempty"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty"`
	Notes     string `json:"notes" binding:"omitempty"`
}

// convertContactToResponse converts a database contact model to an API response
func convertContactToResponse(contact db.Contact) contactResponse {
	createdAt := time.Time{}
	if contact.CreatedAt.Valid {
		createdAt = contact.CreatedAt.Time
	}

	// Convert SQL null strings to regular strings
	position := ""
	if contact.Position.Valid {
		position = contact.Position.String
	}

	email := ""
	if contact.Email.Valid {
		email = contact.Email.String
	}

	phone := ""
	if contact.Phone.Valid {
		phone = contact.Phone.String
	}

	notes := ""
	if contact.Notes.Valid {
		notes = contact.Notes.String
	}

	return contactResponse{
		ContactID: contact.ContactID,
		CompanyID: contact.CompanyID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Position:  position,
		Email:     email,
		Phone:     phone,
		Notes:     notes,
		CreatedAt: createdAt,
	}
}

// Helper function to check if a user has access to a company
func (server *Server) userHasAccessToCompany(ctx *gin.Context, companyID int32, cognitoSub string) (bool, error) {
	company, err := server.store.GetCompanyByID(ctx, companyID)
	if err != nil {
		return false, err
	}
	return company.CognitoSub.Valid && company.CognitoSub.String == cognitoSub, nil
}

// createContact handles requests to create a new contact
func (server *Server) createContact(ctx *gin.Context) {
	var req createContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Check if company exists and belongs to the authenticated user
	company, err := server.store.GetCompanyByID(ctx, req.CompanyID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}

	// Verify that the company belongs to the authenticated user
	if !company.CognitoSub.Valid || company.CognitoSub.String != cognitoSub.(string) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to add contacts to this company"})
		return
	}

	// Convert request to database params
	arg := db.CreateContactParams{
		CompanyID: req.CompanyID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Position:  sql.NullString{String: req.Position, Valid: req.Position != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Notes:     sql.NullString{String: req.Notes, Valid: req.Notes != ""},
	}

	// Create contact in database
	contact, err := server.store.CreateContact(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	// Return created contact as response
	ctx.JSON(http.StatusCreated, convertContactToResponse(contact))
}

// getContactByID handles requests to get a specific contact
func (server *Server) getContactByID(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Get contact from database
	contact, err := server.store.GetContactByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}

	// Check if the contact's company belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify company ownership"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this contact"})
		return
	}

	// Return contact response
	ctx.JSON(http.StatusOK, convertContactToResponse(contact))
}

// listContacts handles requests to get all contacts with pagination
func (server *Server) listContacts(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	_, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// This is a placeholder implementation since there's no direct "list all contacts" query
	// In a real implementation, you would need to add a new query in your SQL files
	// For now, we'll return a 501 Not Implemented
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "This endpoint is not yet implemented"})
}

// listContactsByCompany handles requests to get contacts for a specific company
func (server *Server) listContactsByCompany(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get company ID from URL param
	companyIDParam := ctx.Param("company_id")
	companyID, err := strconv.Atoi(companyIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if the company belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, int32(companyID), cognitoSub.(string))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify company ownership"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access contacts for this company"})
		return
	}

	// Default pagination settings
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	limitParam := ctx.DefaultQuery("limit", "10")
	offsetParam := ctx.DefaultQuery("offset", "0")

	// Convert string params to integers
	limit, err = strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get contacts from database
	contacts, err := server.store.ListContactsByCompanyID(ctx, db.ListContactsByCompanyIDParams{
		CompanyID: int32(companyID),
		Limit:     int32(limit),
		Offset:    int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	// Convert contacts to response format
	contactResponses := make([]contactResponse, len(contacts))
	for i, contact := range contacts {
		contactResponses[i] = convertContactToResponse(contact)
	}

	ctx.JSON(http.StatusOK, contactResponses)
}

// searchContactsByName handles requests to search contacts by name
func (server *Server) searchContactsByName(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get search query from URL param
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Default pagination settings
	limit := 10
	offset := 0

	// Parse query parameters for pagination
	limitParam := ctx.DefaultQuery("limit", "10")
	offsetParam := ctx.DefaultQuery("offset", "0")

	// Convert string params to integers
	var err error
	limit, err = strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset, err = strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get contacts from database
	// Note: This needs to be modified to only return contacts for companies owned by the user
	// You may need to create a new database query for this
	contacts, err := server.store.SearchContactsByName(ctx, db.SearchContactsByNameParams{
		Column1: sql.NullString{String: query, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search contacts"})
		return
	}

	// Filter contacts to only include those from companies owned by the user
	// Note: This is not efficient - ideally, this filtering should be done at the database level
	var authorizedContacts []db.Contact
	for _, contact := range contacts {
		hasAccess, err := server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub.(string))
		if err == nil && hasAccess {
			authorizedContacts = append(authorizedContacts, contact)
		}
	}

	// Convert contacts to response format
	contactResponses := make([]contactResponse, len(authorizedContacts))
	for i, contact := range authorizedContacts {
		contactResponses[i] = convertContactToResponse(contact)
	}

	ctx.JSON(http.StatusOK, contactResponses)
}

// updateContact handles requests to update an existing contact
func (server *Server) updateContact(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Parse request body
	var req updateContactRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if contact exists and get its company ID
	contact, err := server.store.GetContactByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}

	// Check if the contact's company belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify company ownership"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this contact"})
		return
	}

	// Update contact in database
	arg := db.UpdateContactParams{
		ContactID: int32(id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Position:  sql.NullString{String: req.Position, Valid: req.Position != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Notes:     sql.NullString{String: req.Notes, Valid: req.Notes != ""},
	}

	updatedContact, err := server.store.UpdateContact(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact"})
		return
	}

	// Return updated contact
	ctx.JSON(http.StatusOK, convertContactToResponse(updatedContact))
}

// deleteContact handles requests to delete a contact
func (server *Server) deleteContact(ctx *gin.Context) {
	// Get authenticated user's cognito_sub from context
	cognitoSub, exists := ctx.Get("cognito_sub")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	// Get contact ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID format"})
		return
	}

	// Check if contact exists
	contact, err := server.store.GetContactByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contact"})
		return
	}

	// Check if the contact's company belongs to the authenticated user
	hasAccess, err := server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify company ownership"})
		return
	}
	if !hasAccess {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this contact"})
		return
	}

	// Delete contact
	err = server.store.DeleteContact(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
