package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestDBConnection ensures that the database connection is working properly
func TestDBConnection(t *testing.T) {
	// Check that testDB is not nil
	require.NotNil(t, testDB, "testDB should not be nil")

	// Check that testQueries is not nil
	require.NotNil(t, testQueries, "testQueries should not be nil")

	// Test the database connection with a ping
	err := testDB.Ping()
	require.NoError(t, err, "Database ping should succeed")
}
