package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomSalesProcess(t *testing.T) SalesProcess {
	// First, create a random user and contact to associate the sales process with
	user := createRandomUser(t)
	// company := createRandomCompany(t)
	contact := createRandomContact(t)

	arg := CreateSalesProcessParams{
		UserID:               user.UserID,
		ContactID:            contact.ContactID,
		OverallMatchingScore: sql.NullString{String: "0.75", Valid: true},
		Status:               sql.NullString{String: "initial", Valid: true},
	}

	salesProcess, err := testQueries.CreateSalesProcess(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, salesProcess)

	require.Equal(t, arg.UserID, salesProcess.UserID)
	require.Equal(t, arg.ContactID, salesProcess.ContactID)
	require.Equal(t, arg.OverallMatchingScore, salesProcess.OverallMatchingScore)
	require.Equal(t, arg.Status, salesProcess.Status)
	require.NotZero(t, salesProcess.SalesProcessID)
	require.NotZero(t, salesProcess.CreatedAt)
	require.NotZero(t, salesProcess.UpdatedAt)

	return salesProcess
}

func TestCreateSalesProcess(t *testing.T) {
	createRandomSalesProcess(t)
}

func TestGetSalesProcessByID(t *testing.T) {
	// Create a random sales process
	salesProcess1 := createRandomSalesProcess(t)

	// Retrieve the sales process by ID
	salesProcess2, err := testQueries.GetSalesProcessByID(context.Background(), salesProcess1.SalesProcessID)
	require.NoError(t, err)
	require.NotEmpty(t, salesProcess2)

	require.Equal(t, salesProcess1.SalesProcessID, salesProcess2.SalesProcessID)
	require.Equal(t, salesProcess1.UserID, salesProcess2.UserID)
	require.Equal(t, salesProcess1.ContactID, salesProcess2.ContactID)
	require.Equal(t, salesProcess1.OverallMatchingScore, salesProcess2.OverallMatchingScore)
	require.Equal(t, salesProcess1.Status, salesProcess2.Status)
	require.WithinDuration(t, salesProcess1.CreatedAt.Time, salesProcess2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, salesProcess1.UpdatedAt.Time, salesProcess2.UpdatedAt.Time, time.Second)
}

func TestUpdateSalesProcess(t *testing.T) {
	// Create a random sales process
	salesProcess1 := createRandomSalesProcess(t)

	// Prepare updated sales process details
	arg := UpdateSalesProcessParams{
		SalesProcessID:       salesProcess1.SalesProcessID,
		OverallMatchingScore: sql.NullString{String: "0.85", Valid: true},
		Status:               sql.NullString{String: "in_progress", Valid: true},
	}

	// Update the sales process
	salesProcess2, err := testQueries.UpdateSalesProcess(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, salesProcess2)

	// Verify updated details
	require.Equal(t, salesProcess1.SalesProcessID, salesProcess2.SalesProcessID)
	require.Equal(t, salesProcess1.UserID, salesProcess2.UserID)
	require.Equal(t, salesProcess1.ContactID, salesProcess2.ContactID)
	require.Equal(t, arg.OverallMatchingScore, salesProcess2.OverallMatchingScore)
	require.Equal(t, arg.Status, salesProcess2.Status)
	require.True(t, salesProcess2.UpdatedAt.Time.After(salesProcess1.UpdatedAt.Time), "Updated timestamp should be newer")
}

func TestDeleteSalesProcess(t *testing.T) {
	// Create a random sales process
	salesProcess1 := createRandomSalesProcess(t)

	// Delete the sales process
	err := testQueries.DeleteSalesProcess(context.Background(), salesProcess1.SalesProcessID)
	require.NoError(t, err)

	// Try to retrieve the deleted sales process (should fail)
	_, err = testQueries.GetSalesProcessByID(context.Background(), salesProcess1.SalesProcessID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListSalesProcessesByUser(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Create multiple sales processes for the user
	expectedSalesProcesses := 10
	for i := 0; i < expectedSalesProcesses; i++ {
		// Create a new contact for each sales process
		contact := createRandomContact(t)

		arg := CreateSalesProcessParams{
			UserID:    user.UserID,
			ContactID: contact.ContactID,
			Status:    sql.NullString{String: "initial", Valid: true},
		}
		_, err := testQueries.CreateSalesProcess(context.Background(), arg)
		require.NoError(t, err)
	}

	// List sales processes with pagination
	arg := ListSalesProcessesByUserParams{
		UserID: user.UserID,
		Limit:  5,
		Offset: 0,
	}

	salesProcesses, err := testQueries.ListSalesProcessesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, salesProcesses, 5)

	for _, sp := range salesProcesses {
		require.NotEmpty(t, sp)
		require.Equal(t, user.UserID, sp.UserID)
	}
}

func TestListSalesProcessesByContact(t *testing.T) {
	// Create a random contact
	contact := createRandomContact(t)

	// Create multiple sales processes for the contact
	expectedSalesProcesses := 10
	for i := 0; i < expectedSalesProcesses; i++ {
		// Create a new user for each sales process
		user := createRandomUser(t)

		arg := CreateSalesProcessParams{
			UserID:    user.UserID,
			ContactID: contact.ContactID,
			Status:    sql.NullString{String: "initial", Valid: true},
		}
		_, err := testQueries.CreateSalesProcess(context.Background(), arg)
		require.NoError(t, err)
	}

	// List sales processes with pagination
	arg := ListSalesProcessesByContactParams{
		ContactID: contact.ContactID,
		Limit:     5,
		Offset:    0,
	}

	salesProcesses, err := testQueries.ListSalesProcessesByContact(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, salesProcesses, 5)

	for _, sp := range salesProcesses {
		require.NotEmpty(t, sp)
		require.Equal(t, contact.ContactID, sp.ContactID)
	}
}

func TestListSalesProcessesByStatus(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Define test statuses
	testStatuses := []string{"initial", "in_progress", "negotiation", "closed_won", "closed_lost"}

	// Create multiple sales processes with different statuses
	for _, status := range testStatuses {
		// Create a contact for each sales process
		contact := createRandomContact(t)

		arg := CreateSalesProcessParams{
			UserID:    user.UserID,
			ContactID: contact.ContactID,
			Status:    sql.NullString{String: status, Valid: true},
		}
		_, err := testQueries.CreateSalesProcess(context.Background(), arg)
		require.NoError(t, err)
	}

	// Choose a specific status to list
	testStatus := "in_progress"

	// List sales processes by status with pagination
	arg := ListSalesProcessesByStatusParams{
		UserID: user.UserID,
		Status: sql.NullString{String: testStatus, Valid: true},
		Limit:  5,
		Offset: 0,
	}

	salesProcesses, err := testQueries.ListSalesProcessesByStatus(context.Background(), arg)
	require.NoError(t, err)

	for _, sp := range salesProcesses {
		require.NotEmpty(t, sp)
		require.Equal(t, user.UserID, sp.UserID)
		require.Equal(t, testStatus, sp.Status.String)
	}
}

func TestSalesProcessProjectAssociation(t *testing.T) {
	// Create a random sales process
	salesProcess := createRandomSalesProcess(t)

	// Create a random project
	project := createRandomProject(t)

	// Link project to sales process
	err := testQueries.LinkProjectToSalesProcess(context.Background(), LinkProjectToSalesProcessParams{
		SalesProcessID: salesProcess.SalesProcessID,
		ProjectID:      project.ProjectID,
	})
	require.NoError(t, err)

	// Get projects for sales process
	projectsArg := GetProjectsForSalesProcessParams{
		SalesProcessID: salesProcess.SalesProcessID,
		Limit:          5,
		Offset:         0,
	}
	associatedProjects, err := testQueries.GetProjectsForSalesProcess(context.Background(), projectsArg)
	require.NoError(t, err)
	require.NotEmpty(t, associatedProjects)

	// Verify the project is associated
	found := false
	for _, proj := range associatedProjects {
		if proj.ProjectID == project.ProjectID {
			found = true
			break
		}
	}
	require.True(t, found, "Project not found in sales process's projects")

	// Get sales processes for project
	salesProcessesArg := GetSalesProcessesForProjectParams{
		ProjectID: project.ProjectID,
		Limit:     5,
		Offset:    0,
	}
	associatedSalesProcesses, err := testQueries.GetSalesProcessesForProject(context.Background(), salesProcessesArg)
	require.NoError(t, err)
	require.NotEmpty(t, associatedSalesProcesses)

	// Verify the sales process is associated
	foundSP := false
	for _, sp := range associatedSalesProcesses {
		if sp.SalesProcessID == salesProcess.SalesProcessID {
			foundSP = true
			break
		}
	}
	require.True(t, foundSP, "Sales process not found in project's sales processes")

	// Unlink project from sales process
	err = testQueries.UnlinkProjectFromSalesProcess(context.Background(), UnlinkProjectFromSalesProcessParams{
		SalesProcessID: salesProcess.SalesProcessID,
		ProjectID:      project.ProjectID,
	})
	require.NoError(t, err)

	// Re-check projects for sales process (should be empty)
	emptyProjects, err := testQueries.GetProjectsForSalesProcess(context.Background(), projectsArg)
	require.NoError(t, err)
	require.Len(t, emptyProjects, 0)
}
