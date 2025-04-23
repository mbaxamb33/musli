package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Valid source types for the check_source_type constraint
// These match exactly what's defined in the database constraint:
// CHECK (source_type IN ('website', 'page', 'api', 'file', 'manual', 'other'))
var validSourceTypes = []string{
	"website",
	"page",
	"api",
	"file",
	"manual",
	"other",
}

// createRandomDatasource creates a datasource for testing with valid source_type
func createRandomDatasource(t *testing.T) Datasource {
	// Select a valid source type from the predefined list
	sourceType := validSourceTypes[randomInt(0, int64(len(validSourceTypes)-1))]

	arg := CreateDatasourceParams{
		SourceType: sourceType,
		SourceID: sql.NullInt32{
			Int32: int32(randomInt(1, 1000)),
			Valid: true,
		},
	}

	datasource, err := testQueries.CreateDatasource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, datasource)

	require.Equal(t, arg.SourceType, datasource.SourceType)
	require.Equal(t, arg.SourceID, datasource.SourceID)
	require.NotZero(t, datasource.DatasourceID)
	require.NotEmpty(t, datasource.ExtractionTimestamp)

	return datasource
}

// TestCreateDatasource tests the CreateDatasource function
func TestCreateDatasource(t *testing.T) {
	// Already tested in createRandomDatasource which is used in other tests
	datasource := createRandomDatasource(t)
	require.NotZero(t, datasource.DatasourceID)
}

// TestGetDatasourceByID tests the GetDatasourceByID function
func TestGetDatasourceByID(t *testing.T) {
	datasource1 := createRandomDatasource(t)
	datasource2, err := testQueries.GetDatasourceByID(context.Background(), datasource1.DatasourceID)
	require.NoError(t, err)
	require.NotEmpty(t, datasource2)

	require.Equal(t, datasource1.DatasourceID, datasource2.DatasourceID)
	require.Equal(t, datasource1.SourceType, datasource2.SourceType)
	require.Equal(t, datasource1.SourceID, datasource2.SourceID)
	require.WithinDuration(t, datasource1.ExtractionTimestamp.Time, datasource2.ExtractionTimestamp.Time, time.Second)
}

// TestListDatasources tests the ListDatasources function
func TestListDatasources(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomDatasource(t)
	}

	arg := ListDatasourcesParams{
		Limit:  5,
		Offset: 0,
	}

	datasources, err := testQueries.ListDatasources(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, datasources)
	require.Len(t, datasources, 5)

	for _, datasource := range datasources {
		require.NotEmpty(t, datasource)
		require.NotZero(t, datasource.DatasourceID)
	}

	// Test pagination
	arg2 := ListDatasourcesParams{
		Limit:  5,
		Offset: 5,
	}

	datasources2, err := testQueries.ListDatasources(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, datasources2)
	require.Len(t, datasources2, 5)

	// Make sure the two sets of datasources are different
	datasourceMap := make(map[int32]bool)
	for _, d := range datasources {
		datasourceMap[d.DatasourceID] = true
	}

	for _, d := range datasources2 {
		_, exists := datasourceMap[d.DatasourceID]
		require.False(t, exists, "Datasource appears in both result sets")
	}
}

// TestListDatasourcesByType tests the ListDatasourcesByType function
func TestListDatasourcesByType(t *testing.T) {
	// Choose a valid source type that's less likely to be in the database already
	// We'll use "file" since it's not as commonly used as "website" or "page"
	sourceType := "file"

	// Delete any existing datasources with this type to ensure clean test
	// This is just for test reliability and not needed in production
	cleanupExistingDatasources(t, sourceType)

	// Create several datasources with the same type
	for i := 0; i < 5; i++ {
		arg := CreateDatasourceParams{
			SourceType: sourceType,
			SourceID: sql.NullInt32{
				Int32: int32(randomInt(1, 1000)),
				Valid: true,
			},
		}

		_, err := testQueries.CreateDatasource(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some datasources with different types
	for i := 0; i < 3; i++ {
		createRandomDatasource(t)
	}

	arg := ListDatasourcesByTypeParams{
		SourceType: sourceType,
		Limit:      10,
		Offset:     0,
	}

	datasources, err := testQueries.ListDatasourcesByType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, datasources)
	require.Len(t, datasources, 5)

	// Verify all datasources have the same type
	for _, datasource := range datasources {
		require.Equal(t, sourceType, datasource.SourceType)
	}
}

// Helper function to clean up existing datasources of a specific type before testing
func cleanupExistingDatasources(t *testing.T, sourceType string) {
	// First get the existing datasources of this type
	arg := ListDatasourcesByTypeParams{
		SourceType: sourceType,
		Limit:      100, // Arbitrary high number to get all
		Offset:     0,
	}

	datasources, err := testQueries.ListDatasourcesByType(context.Background(), arg)
	if err != nil {
		// If there's an error, just return - this is just a cleanup helper
		return
	}

	// Delete each datasource
	for _, ds := range datasources {
		err := testQueries.DeleteDatasource(context.Background(), ds.DatasourceID)
		if err != nil {
			// If we can't delete, just continue - this is just a cleanup helper
			continue
		}
	}
}

// TestDeleteDatasource tests the DeleteDatasource function
func TestDeleteDatasource(t *testing.T) {
	datasource1 := createRandomDatasource(t)
	err := testQueries.DeleteDatasource(context.Background(), datasource1.DatasourceID)
	require.NoError(t, err)

	datasource2, err := testQueries.GetDatasourceByID(context.Background(), datasource1.DatasourceID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, datasource2)
}
