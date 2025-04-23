package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomUser creates a user for project testing
func createRandomUserForProject(t *testing.T) User {
	return createRandomUser(t)
}

// createRandomProject creates a project with random values for testing
func createRandomProject(t *testing.T) Project {
	user := createRandomUserForProject(t)

	arg := CreateProjectParams{
		UserID: sql.NullInt32{
			Int32: user.UserID,
			Valid: true,
		},
		ProjectName: randomString(10),
		Description: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	project, err := testQueries.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	require.Equal(t, arg.UserID, project.UserID)
	require.Equal(t, arg.ProjectName, project.ProjectName)
	require.Equal(t, arg.Description, project.Description)

	require.NotZero(t, project.ProjectID)
	require.NotEmpty(t, project.CreatedAt)
	require.NotEmpty(t, project.LastUpdatedAt)

	return project
}

// TestCreateProject tests the CreateProject function
func TestCreateProject(t *testing.T) {
	createRandomProject(t)
}

// TestGetProjectByID tests the GetProjectByID function
func TestGetProjectByID(t *testing.T) {
	project1 := createRandomProject(t)
	project2, err := testQueries.GetProjectByID(context.Background(), project1.ProjectID)
	require.NoError(t, err)
	require.NotEmpty(t, project2)

	require.Equal(t, project1.ProjectID, project2.ProjectID)
	require.Equal(t, project1.UserID, project2.UserID)
	require.Equal(t, project1.ProjectName, project2.ProjectName)
	require.Equal(t, project1.Description, project2.Description)
	require.WithinDuration(t, project1.CreatedAt.Time, project2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, project1.LastUpdatedAt.Time, project2.LastUpdatedAt.Time, time.Second)
}

// TestListProjectsByUser tests the ListProjectsByUser function
func TestListProjectsByUser(t *testing.T) {
	user := createRandomUserForProject(t)

	// Create several projects for the same user
	for i := 0; i < 5; i++ {
		arg := CreateProjectParams{
			UserID: sql.NullInt32{
				Int32: user.UserID,
				Valid: true,
			},
			ProjectName: randomString(10),
			Description: sql.NullString{
				String: randomString(20),
				Valid:  true,
			},
		}

		project, err := testQueries.CreateProject(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, project)
	}

	// Also create some projects for different users
	for i := 0; i < 3; i++ {
		createRandomProject(t)
	}

	arg := ListProjectsByUserParams{
		UserID: sql.NullInt32{
			Int32: user.UserID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	projects, err := testQueries.ListProjectsByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, projects)
	require.Len(t, projects, 5)

	// Verify all projects belong to the same user
	for _, project := range projects {
		require.Equal(t, user.UserID, project.UserID.Int32)
	}
}

// TestUpdateProject tests the UpdateProject function
func TestUpdateProject(t *testing.T) {
	project1 := createRandomProject(t)

	arg := UpdateProjectParams{
		ProjectID:   project1.ProjectID,
		ProjectName: randomString(10),
		Description: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	project2, err := testQueries.UpdateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project2)

	require.Equal(t, project1.ProjectID, project2.ProjectID)
	require.Equal(t, project1.UserID, project2.UserID)
	require.Equal(t, arg.ProjectName, project2.ProjectName)
	require.Equal(t, arg.Description, project2.Description)
	require.WithinDuration(t, project1.CreatedAt.Time, project2.CreatedAt.Time, time.Second)

	// Last updated time should be more recent
	require.True(t, project2.LastUpdatedAt.Time.After(project1.LastUpdatedAt.Time) ||
		project2.LastUpdatedAt.Time.Equal(project1.LastUpdatedAt.Time))
}

// TestDeleteProject tests the DeleteProject function
func TestDeleteProject(t *testing.T) {
	project1 := createRandomProject(t)
	err := testQueries.DeleteProject(context.Background(), project1.ProjectID)
	require.NoError(t, err)

	project2, err := testQueries.GetProjectByID(context.Background(), project1.ProjectID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, project2)
}
