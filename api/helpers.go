package api

// import (
// 	"database/sql"

// 	"github.com/gin-gonic/gin"
// )

// // Helper function to check if user has access to a datasource
// func (server *Server) userHasAccessToDatasource(ctx *gin.Context, datasourceID int32, cognitoSub string) (bool, error) {
// 	// This is a complex check because datasources can be linked to companies, contacts, or projects

// 	// First, check if the datasource is linked to a company owned by the user
// 	companyAssociations, err := server.store.ListCompanyDatasourceAssociations(ctx, datasourceID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return false, err
// 	}

// 	// Check each company association
// 	for _, association := range companyAssociations {
// 		hasAccess, err := server.userHasAccessToCompany(ctx, association.CompanyID, cognitoSub)
// 		if err == nil && hasAccess {
// 			return true, nil
// 		}
// 	}

// 	// Next, check if the datasource is linked to a contact in a company owned by the user
// 	contactAssociations, err := server.store.ListContactDatasourceAssociations(ctx, datasourceID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return false, err
// 	}

// 	// Check each contact association
// 	for _, association := range contactAssociations {
// 		contact, err := server.store.GetContactByID(ctx, association.ContactID)
// 		if err == nil {
// 			hasAccess, err := server.userHasAccessToCompany(ctx, contact.CompanyID, cognitoSub)
// 			if err == nil && hasAccess {
// 				return true, nil
// 			}
// 		}
// 	}

// 	// Finally, check if it's linked to a project owned by the user
// 	projectAssociations, err := server.store.ListProjectDatasourceAssociations(ctx, datasourceID)
// 	if err != nil && err != sql.ErrNoRows {
// 		return false, err
// 	}

// 	// Check each project association
// 	for _, association := range projectAssociations {
// 		project, err := server.store.GetProjectByID(ctx, association.ProjectID)
// 		if err == nil && project.CognitoSub.Valid && project.CognitoSub.String == cognitoSub {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }
