package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomCompany(t *testing.T) Company {
	// First, create a random user to associate with the company
	user := createRandomUser(t)

	arg := CreateCompanyParams{
		UserID:      user.UserID,
		CompanyName: randomString(10),
		Industry:    sql.NullString{String: randomString(8), Valid: true},
		Website:     sql.NullString{String: "https://" + randomString(6) + ".com", Valid: true},
		Address:     sql.NullString{String: randomString(15) + " Street", Valid: true},
		Description: sql.NullString{String: "A test company description", Valid: true},
	}

	company, err := testQueries.CreateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company)

	require.Equal(t, arg.UserID, company.UserID)
	require.Equal(t, arg.CompanyName, company.CompanyName)
	require.Equal(t, arg.Industry, company.Industry)
	require.Equal(t, arg.Website, company.Website)
	require.Equal(t, arg.Address, company.Address)
	require.Equal(t, arg.Description, company.Description)
	require.NotZero(t, company.CompanyID)
	require.NotZero(t, company.CreatedAt)

	return company
}

func TestCreateCompany(t *testing.T) {
	createRandomCompany(t)
}

func TestGetCompanyByID(t *testing.T) {
	// Create a random company
	company1 := createRandomCompany(t)

	// Retrieve the company by ID
	company2, err := testQueries.GetCompanyByID(context.Background(), company1.CompanyID)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, company1.UserID, company2.UserID)
	require.Equal(t, company1.CompanyName, company2.CompanyName)
	require.Equal(t, company1.Industry, company2.Industry)
	require.Equal(t, company1.Website, company2.Website)
	require.Equal(t, company1.Address, company2.Address)
	require.Equal(t, company1.Description, company2.Description)
	require.WithinDuration(t, company1.CreatedAt.Time, company2.CreatedAt.Time, time.Second)
}

func TestGetCompanyByName(t *testing.T) {
	// Create a random company
	company1 := createRandomCompany(t)

	// Retrieve the company by name
	company2, err := testQueries.GetCompanyByName(context.Background(), company1.CompanyName)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, company1.UserID, company2.UserID)
	require.Equal(t, company1.CompanyName, company2.CompanyName)
	require.Equal(t, company1.Industry, company2.Industry)
	require.Equal(t, company1.Website, company2.Website)
	require.Equal(t, company1.Address, company2.Address)
	require.Equal(t, company1.Description, company2.Description)
	require.WithinDuration(t, company1.CreatedAt.Time, company2.CreatedAt.Time, time.Second)
}

func TestUpdateCompany(t *testing.T) {
	// Create a random company
	company1 := createRandomCompany(t)

	// Prepare updated company details
	arg := UpdateCompanyParams{
		CompanyID:   company1.CompanyID,
		CompanyName: randomString(12),
		Industry:    sql.NullString{String: randomString(10), Valid: true},
		Website:     sql.NullString{String: "https://" + randomString(8) + ".org", Valid: true},
		Address:     sql.NullString{String: randomString(20) + " Avenue", Valid: true},
		Description: sql.NullString{String: "An updated test company description", Valid: true},
	}

	// Update the company
	company2, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	// Verify updated details
	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, company1.UserID, company2.UserID)
	require.Equal(t, arg.CompanyName, company2.CompanyName)
	require.Equal(t, arg.Industry, company2.Industry)
	require.Equal(t, arg.Website, company2.Website)
	require.Equal(t, arg.Address, company2.Address)
	require.Equal(t, arg.Description, company2.Description)
}

func TestDeleteCompany(t *testing.T) {
	// Create a random company
	company1 := createRandomCompany(t)

	// Delete the company
	err := testQueries.DeleteCompany(context.Background(), company1.CompanyID)
	require.NoError(t, err)

	// Try to retrieve the deleted company (should fail)
	_, err = testQueries.GetCompanyByID(context.Background(), company1.CompanyID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListCompaniesByUserID(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Create multiple companies for the user
	expectedCompanies := 10
	for i := 0; i < expectedCompanies; i++ {
		arg := CreateCompanyParams{
			UserID:      user.UserID,
			CompanyName: randomString(10),
			Industry:    sql.NullString{String: randomString(8), Valid: true},
		}
		_, err := testQueries.CreateCompany(context.Background(), arg)
		require.NoError(t, err)
	}

	// List companies with pagination
	arg := GetCompaniesByUserIDParams{
		UserID: user.UserID,
		Limit:  5,
		Offset: 0,
	}

	companies, err := testQueries.GetCompaniesByUserID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, companies, 5)

	for _, company := range companies {
		require.NotEmpty(t, company)
		require.Equal(t, user.UserID, company.UserID)
	}
}

func TestCreateCompanyWithDuplicateName(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Create initial company
	arg1 := CreateCompanyParams{
		UserID:      user.UserID,
		CompanyName: "UniqueCompanyName",
		Industry:    sql.NullString{String: randomString(8), Valid: true},
	}
	_, err := testQueries.CreateCompany(context.Background(), arg1)
	require.NoError(t, err)

	// Try to create another company with the same name (behavior depends on your specific database constraints)
	arg2 := CreateCompanyParams{
		UserID:      user.UserID,
		CompanyName: "UniqueCompanyName",
		Industry:    sql.NullString{String: randomString(8), Valid: true},
	}

	// This test assumes the database does NOT prevent duplicate company names
	// If your database has a unique constraint, modify this test accordingly
	company2, err := testQueries.CreateCompany(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, company2)
}

func TestListCompanies(t *testing.T) {
	// Create multiple companies for different users
	for i := 0; i < 10; i++ {
		createRandomCompany(t)
	}

	// List companies with pagination
	arg := ListCompaniesParams{
		Limit:  5,
		Offset: 0,
	}

	companies, err := testQueries.ListCompanies(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, companies, 5)

	for _, company := range companies {
		require.NotEmpty(t, company)
	}
}
