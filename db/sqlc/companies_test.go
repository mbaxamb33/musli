package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomCompany creates a company with random values for testing
func createRandomCompany(t *testing.T) Company {
	arg := CreateCompanyParams{
		Name: randomString(10),
		Website: sql.NullString{
			String: fmt.Sprintf("https://www.%s.com", randomString(8)),
			Valid:  true,
		},
		Industry: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		Description: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		HeadquartersLocation: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		FoundedYear: sql.NullInt32{
			Int32: int32(randomInt(1950, 2023)),
			Valid: true,
		},
		IsPublic: sql.NullBool{
			Bool:  randomBool(),
			Valid: true,
		},
		TickerSymbol: sql.NullString{
			String: randomString(4),
			Valid:  true,
		},
	}

	company, err := testQueries.CreateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company)

	require.Equal(t, arg.Name, company.Name)
	require.Equal(t, arg.Website, company.Website)
	require.Equal(t, arg.Industry, company.Industry)
	require.Equal(t, arg.Description, company.Description)
	require.Equal(t, arg.HeadquartersLocation, company.HeadquartersLocation)
	require.Equal(t, arg.FoundedYear, company.FoundedYear)
	require.Equal(t, arg.IsPublic, company.IsPublic)
	require.Equal(t, arg.TickerSymbol, company.TickerSymbol)

	require.NotZero(t, company.CompanyID)
	require.NotEmpty(t, company.ScrapeTimestamp)

	return company
}

// TestCreateCompany tests the CreateCompany function
func TestCreateCompany(t *testing.T) {
	createRandomCompany(t)
}

// TestGetCompanyByID tests the GetCompanyByID function
func TestGetCompanyByID(t *testing.T) {
	company1 := createRandomCompany(t)
	company2, err := testQueries.GetCompanyByID(context.Background(), company1.CompanyID)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, company1.Name, company2.Name)
	require.Equal(t, company1.Website, company2.Website)
	require.Equal(t, company1.Industry, company2.Industry)
	require.Equal(t, company1.Description, company2.Description)
	require.Equal(t, company1.HeadquartersLocation, company2.HeadquartersLocation)
	require.Equal(t, company1.FoundedYear, company2.FoundedYear)
	require.Equal(t, company1.IsPublic, company2.IsPublic)
	require.Equal(t, company1.TickerSymbol, company2.TickerSymbol)
	require.WithinDuration(t, company1.ScrapeTimestamp.Time, company2.ScrapeTimestamp.Time, time.Second)
}

// TestGetCompanyByName tests the GetCompanyByName function
func TestGetCompanyByName(t *testing.T) {
	company1 := createRandomCompany(t)
	company2, err := testQueries.GetCompanyByName(context.Background(), company1.Name)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, company1.Name, company2.Name)
	require.Equal(t, company1.Website, company2.Website)
	require.Equal(t, company1.Industry, company2.Industry)
	require.Equal(t, company1.Description, company2.Description)
	require.Equal(t, company1.HeadquartersLocation, company2.HeadquartersLocation)
	require.Equal(t, company1.FoundedYear, company2.FoundedYear)
	require.Equal(t, company1.IsPublic, company2.IsPublic)
	require.Equal(t, company1.TickerSymbol, company2.TickerSymbol)
	require.WithinDuration(t, company1.ScrapeTimestamp.Time, company2.ScrapeTimestamp.Time, time.Second)
}

// TestUpdateCompany tests the UpdateCompany function
func TestUpdateCompany(t *testing.T) {
	company1 := createRandomCompany(t)

	arg := UpdateCompanyParams{
		CompanyID: company1.CompanyID,
		Name:      randomString(10),
		Website: sql.NullString{
			String: fmt.Sprintf("https://www.%s.com", randomString(8)),
			Valid:  true,
		},
		Industry: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		Description: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		HeadquartersLocation: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		FoundedYear: sql.NullInt32{
			Int32: int32(randomInt(1950, 2023)),
			Valid: true,
		},
		IsPublic: sql.NullBool{
			Bool:  !company1.IsPublic.Bool, // Toggle the value
			Valid: true,
		},
		TickerSymbol: sql.NullString{
			String: randomString(4),
			Valid:  true,
		},
	}

	company2, err := testQueries.UpdateCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, company2)

	require.Equal(t, company1.CompanyID, company2.CompanyID)
	require.Equal(t, arg.Name, company2.Name)
	require.Equal(t, arg.Website, company2.Website)
	require.Equal(t, arg.Industry, company2.Industry)
	require.Equal(t, arg.Description, company2.Description)
	require.Equal(t, arg.HeadquartersLocation, company2.HeadquartersLocation)
	require.Equal(t, arg.FoundedYear, company2.FoundedYear)
	require.Equal(t, arg.IsPublic, company2.IsPublic)
	require.Equal(t, arg.TickerSymbol, company2.TickerSymbol)
	require.WithinDuration(t, company1.ScrapeTimestamp.Time, company2.ScrapeTimestamp.Time, time.Second*5)
}

// TestListCompanies tests the ListCompanies function
func TestListCompanies(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCompany(t)
	}

	arg := ListCompaniesParams{
		Limit:  5,
		Offset: 0,
	}

	companies, err := testQueries.ListCompanies(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, companies)
	require.Len(t, companies, 5)

	for _, company := range companies {
		require.NotEmpty(t, company)
		require.NotZero(t, company.CompanyID)
	}

	// Test pagination
	arg2 := ListCompaniesParams{
		Limit:  5,
		Offset: 5,
	}

	companies2, err := testQueries.ListCompanies(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, companies2)
	require.Len(t, companies2, 5)

	// Verify different sets of companies
	companiesMap := make(map[int32]bool)
	for _, company := range companies {
		companiesMap[company.CompanyID] = true
	}

	for _, company := range companies2 {
		_, exists := companiesMap[company.CompanyID]
		require.False(t, exists, "Company appears in both result sets")
	}
}

// TestDeleteCompany tests the DeleteCompany function
func TestDeleteCompany(t *testing.T) {
	company1 := createRandomCompany(t)
	err := testQueries.DeleteCompany(context.Background(), company1.CompanyID)
	require.NoError(t, err)

	company2, err := testQueries.GetCompanyByID(context.Background(), company1.CompanyID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, company2)
}
