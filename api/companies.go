package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/mbaxamb3/nusli/db/sqlc"
)

// createCompanyRequest represents the request body for creating a company
type createCompanyRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	Industry    string `json:"industry" binding:"omitempty"`
	Website     string `json:"website" binding:"omitempty"`
	Address     string `json:"address" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

// companyResponse represents the API response structure for company data
type companyResponse struct {
	CompanyID   int32     `json:"company_id"`
	CognitoSub  string    `json:"cognito_sub,omitempty"`
	CompanyName string    `json:"company_name"`
	Industry    string    `json:"industry,omitempty"`
	Website     string    `json:"website,omitempty"`
	Address     string    `json:"address,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// updateCompanyRequest represents the request body for updating a company
type updateCompanyRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	Industry    string `json:"industry" binding:"omitempty"`
	Website     string `json:"website" binding:"omitempty"`
	Address     string `json:"address" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
}

// convertCompanyToResponse converts a database company model to an API response
func convertCompanyToResponse(company db.GetCompanyByIDRow) companyResponse {
	createdAt := time.Time{}
	if company.CreatedAt.Valid {
		createdAt = company.CreatedAt.Time
	}

	// Convert SQL null strings to regular strings
	industry := ""
	if company.Industry.Valid {
		industry = company.Industry.String
	}

	website := ""
	if company.Website.Valid {
		website = company.Website.String
	}

	address := ""
	if company.Address.Valid {
		address = company.Address.String
	}

	description := ""
	if company.Description.Valid {
		description = company.Description.String
	}

	cognitoSub := ""
	if company.CognitoSub.Valid {
		cognitoSub = company.CognitoSub.String
	}

	return companyResponse{
		CompanyID:   company.CompanyID,
		CognitoSub:  cognitoSub,
		CompanyName: company.CompanyName,
		Industry:    industry,
		Website:     website,
		Address:     address,
		Description: description,
		CreatedAt:   createdAt,
	}
}

// Function overload to handle different company struct types
func convertCompanyListToResponse(company db.ListCompaniesRow) companyResponse {
	createdAt := time.Time{}
	if company.CreatedAt.Valid {
		createdAt = company.CreatedAt.Time
	}

	// Convert SQL null strings to regular strings
	industry := ""
	if company.Industry.Valid {
		industry = company.Industry.String
	}

	website := ""
	if company.Website.Valid {
		website = company.Website.String
	}

	address := ""
	if company.Address.Valid {
		address = company.Address.String
	}

	description := ""
	if company.Description.Valid {
		description = company.Description.String
	}

	cognitoSub := ""
	if company.CognitoSub.Valid {
		cognitoSub = company.CognitoSub.String
	}

	return companyResponse{
		CompanyID:   company.CompanyID,
		CognitoSub:  cognitoSub,
		CompanyName: company.CompanyName,
		Industry:    industry,
		Website:     website,
		Address:     address,
		Description: description,
		CreatedAt:   createdAt,
	}
}

// Function overload to handle cognito sub company rows
func convertCompanyCognitoToResponse(company db.GetCompaniesByCognitoSubRow) companyResponse {
	createdAt := time.Time{}
	if company.CreatedAt.Valid {
		createdAt = company.CreatedAt.Time
	}

	// Convert SQL null strings to regular strings
	industry := ""
	if company.Industry.Valid {
		industry = company.Industry.String
	}

	website := ""
	if company.Website.Valid {
		website = company.Website.String
	}

	address := ""
	if company.Address.Valid {
		address = company.Address.String
	}

	description := ""
	if company.Description.Valid {
		description = company.Description.String
	}

	cognitoSub := ""
	if company.CognitoSub.Valid {
		cognitoSub = company.CognitoSub.String
	}

	return companyResponse{
		CompanyID:   company.CompanyID,
		CognitoSub:  cognitoSub,
		CompanyName: company.CompanyName,
		Industry:    industry,
		Website:     website,
		Address:     address,
		Description: description,
		CreatedAt:   createdAt,
	}
}

// createCompany handles requests to create a new company
func (server *Server) createCompany(ctx *gin.Context) {
	var req createCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get cognito_sub from authentication context
	// In a real application, you would get this from authenticated user
	cognitoSub := "default_cognito_sub" // TODO: Replace with actual authenticated user's cognito_sub

	// Convert request to database params
	arg := db.CreateCompanyParams{
		CognitoSub:  sql.NullString{String: cognitoSub, Valid: true},
		CompanyName: req.CompanyName,
		Industry:    sql.NullString{String: req.Industry, Valid: req.Industry != ""},
		Website:     sql.NullString{String: req.Website, Valid: req.Website != ""},
		Address:     sql.NullString{String: req.Address, Valid: req.Address != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	// Create company in database
	company, err := server.store.CreateCompany(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	// Return created company as response
	result := companyResponse{
		CompanyID:   company.CompanyID,
		CognitoSub:  cognitoSub,
		CompanyName: company.CompanyName,
		Industry:    req.Industry,
		Website:     req.Website,
		Address:     req.Address,
		Description: req.Description,
		CreatedAt:   company.CreatedAt.Time,
	}
	ctx.JSON(http.StatusCreated, result)
}

// getCompanyByID handles requests to get a specific company
func (server *Server) getCompanyByID(ctx *gin.Context) {
	// Get company ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Get company from database
	company, err := server.store.GetCompanyByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}

	// Return company response
	ctx.JSON(http.StatusOK, convertCompanyToResponse(company))
}

// listCompanies handles requests to get all companies with pagination
func (server *Server) listCompanies(ctx *gin.Context) {
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

	// Get companies from database
	companies, err := server.store.ListCompanies(ctx, db.ListCompaniesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	// Convert companies to response format
	companyResponses := make([]companyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = convertCompanyListToResponse(company)
	}

	ctx.JSON(http.StatusOK, companyResponses)
}

// getCompaniesByUser handles requests to get companies for a specific user
func (server *Server) getCompaniesByUser(ctx *gin.Context) {
	// Get cognito_sub from URL param
	cognitoSub := ctx.Param("cognito_sub")

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

	// Get companies from database
	companies, err := server.store.GetCompaniesByCognitoSub(ctx, db.GetCompaniesByCognitoSubParams{
		CognitoSub: sql.NullString{String: cognitoSub, Valid: true},
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	// Convert companies to response format
	companyResponses := make([]companyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = convertCompanyCognitoToResponse(company)
	}

	ctx.JSON(http.StatusOK, companyResponses)
}

// updateCompany handles requests to update an existing company
func (server *Server) updateCompany(ctx *gin.Context) {
	// Get company ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Parse request body
	var req updateCompanyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if company exists
	_, err = server.store.GetCompanyByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}

	// Update company in database
	arg := db.UpdateCompanyParams{
		CompanyID:   int32(id),
		CompanyName: req.CompanyName,
		Industry:    sql.NullString{String: req.Industry, Valid: req.Industry != ""},
		Website:     sql.NullString{String: req.Website, Valid: req.Website != ""},
		Address:     sql.NullString{String: req.Address, Valid: req.Address != ""},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	company, err := server.store.UpdateCompany(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	// Return updated company
	response := companyResponse{
		CompanyID:   company.CompanyID,
		CompanyName: company.CompanyName,
		CognitoSub:  company.CognitoSub.String,
		Industry:    req.Industry,
		Website:     req.Website,
		Address:     req.Address,
		Description: req.Description,
		CreatedAt:   company.CreatedAt.Time,
	}
	ctx.JSON(http.StatusOK, response)
}

// deleteCompany handles requests to delete a company
func (server *Server) deleteCompany(ctx *gin.Context) {
	// Get company ID from URL param
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID format"})
		return
	}

	// Check if company exists
	_, err = server.store.GetCompanyByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch company"})
		return
	}

	// Delete company
	err = server.store.DeleteCompany(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
