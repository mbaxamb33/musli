package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomProject(t *testing.T) Project {
	// First, create a random user to associate the project with
	user := createRandomUser(t)

	arg := CreateProjectParams{
		UserID:      user.UserID,
		ProjectName: randomString(10),
		MainIdea:    sql.NullString{String: "A test project main idea: " + randomString(20), Valid: true},
	}

	project, err := testQueries.CreateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project)

	require.Equal(t, arg.UserID, project.UserID)
	require.Equal(t, arg.ProjectName, project.ProjectName)
	require.Equal(t, arg.MainIdea, project.MainIdea)
	require.NotZero(t, project.ProjectID)
	require.NotZero(t, project.CreatedAt)
	require.NotZero(t, project.UpdatedAt)

	return project
}

func TestCreateProject(t *testing.T) {
	createRandomProject(t)
}

func TestGetProjectByID(t *testing.T) {
	// Create a random project
	project1 := createRandomProject(t)

	// Retrieve the project by ID
	project2, err := testQueries.GetProjectByID(context.Background(), project1.ProjectID)
	require.NoError(t, err)
	require.NotEmpty(t, project2)

	require.Equal(t, project1.ProjectID, project2.ProjectID)
	require.Equal(t, project1.UserID, project2.UserID)
	require.Equal(t, project1.ProjectName, project2.ProjectName)
	require.Equal(t, project1.MainIdea, project2.MainIdea)
	require.WithinDuration(t, project1.CreatedAt.Time, project2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, project1.UpdatedAt.Time, project2.UpdatedAt.Time, time.Second)
}

func TestUpdateProject(t *testing.T) {
	// Create a random project
	project1 := createRandomProject(t)

	// Prepare updated project details
	arg := UpdateProjectParams{
		ProjectID:   project1.ProjectID,
		ProjectName: randomString(12),
		MainIdea:    sql.NullString{String: "An updated test project main idea: " + randomString(25), Valid: true},
	}

	// Update the project
	project2, err := testQueries.UpdateProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, project2)

	// Verify updated details
	require.Equal(t, project1.ProjectID, project2.ProjectID)
	require.Equal(t, project1.UserID, project2.UserID)
	require.Equal(t, arg.ProjectName, project2.ProjectName)
	require.Equal(t, arg.MainIdea, project2.MainIdea)
	require.WithinDuration(t, project1.CreatedAt.Time, project2.CreatedAt.Time, time.Second)
	require.True(t, project2.UpdatedAt.Time.After(project1.UpdatedAt.Time), "Updated timestamp should be newer")
}

func TestDeleteProject(t *testing.T) {
	// Create a random project
	project1 := createRandomProject(t)

	// Delete the project
	err := testQueries.DeleteProject(context.Background(), project1.ProjectID)
	require.NoError(t, err)

	// Try to retrieve the deleted project (should fail)
	_, err = testQueries.GetProjectByID(context.Background(), project1.ProjectID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListProjectsByUserID(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Create multiple projects for the user
	expectedProjects := 10
	for i := 0; i < expectedProjects; i++ {
		arg := CreateProjectParams{
			UserID:      user.UserID,
			ProjectName: randomString(10),
			MainIdea:    sql.NullString{String: "Project " + randomString(15), Valid: true},
		}
		_, err := testQueries.CreateProject(context.Background(), arg)
		require.NoError(t, err)
	}

	// List projects with pagination
	arg := ListProjectsByUserIDParams{
		UserID: user.UserID,
		Limit:  5,
		Offset: 0,
	}

	projects, err := testQueries.ListProjectsByUserID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, projects, 5)

	for _, project := range projects {
		require.NotEmpty(t, project)
		require.Equal(t, user.UserID, project.UserID)
	}
}

func TestSearchProjectsByName(t *testing.T) {
	// Create a random user
	user := createRandomUser(t)

	// Create multiple projects
	projects := make([]Project, 5)
	for i := 0; i < 5; i++ {
		arg := CreateProjectParams{
			UserID:      user.UserID,
			ProjectName: "UniquePrefix" + randomString(8),
			MainIdea:    sql.NullString{String: "Project " + randomString(15), Valid: true},
		}
		projects[i], _ = testQueries.CreateProject(context.Background(), arg)
	}

	// Choose a project to search for
	searchProject := projects[2]
	searchTerm := searchProject.ProjectName[0:6] // Use first 6 characters

	// Search projects by name
	arg := SearchProjectsByNameParams{
		UserID:  user.UserID,
		Column2: sql.NullString{String: searchTerm, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	foundProjects, err := testQueries.SearchProjectsByName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, foundProjects)

	// Verify that the search results contain the searched project
	found := false
	for _, project := range foundProjects {
		require.Equal(t, user.UserID, project.UserID)
		if project.ProjectID == searchProject.ProjectID {
			found = true
			break
		}
	}
	require.True(t, found, "Search term did not return the expected project")
}

func TestProjectDatasourceAssociation(t *testing.T) {
	// Create a random project
	project := createRandomProject(t)

	// Create a random datasource
	datasourceArg := CreateDatasourceParams{
		SourceType: DatasourceTypePdf,
		FileName:   sql.NullString{String: randomString(10) + ".pdf", Valid: true},
		Link:       sql.NullString{String: "https://example.com/" + randomString(8), Valid: true},
		FileData:   []byte(randomString(50)),
	}
	datasource, err := testQueries.CreateDatasource(context.Background(), datasourceArg)
	require.NoError(t, err)

	// Associate datasource with project
	err = testQueries.AssociateDatasourceWithProject(context.Background(), AssociateDatasourceWithProjectParams{
		ProjectID:    project.ProjectID,
		DatasourceID: datasource.DatasourceID,
	})
	require.NoError(t, err)

	// Verify the association
	association, err := testQueries.GetProjectDatasourceAssociation(context.Background(), GetProjectDatasourceAssociationParams{
		ProjectID:    project.ProjectID,
		DatasourceID: datasource.DatasourceID,
	})
	require.NoError(t, err)
	require.Equal(t, project.ProjectID, association.ProjectID)
	require.Equal(t, datasource.DatasourceID, association.DatasourceID)

	// List datasources for the project
	datasourcesArg := ListDatasourcesByProjectParams{
		ProjectID: project.ProjectID,
		Limit:     5,
		Offset:    0,
	}
	projectDatasources, err := testQueries.ListDatasourcesByProject(context.Background(), datasourcesArg)
	require.NoError(t, err)
	require.NotEmpty(t, projectDatasources)

	// Verify the datasource is in the list
	found := false
	for _, ds := range projectDatasources {
		if ds.DatasourceID == datasource.DatasourceID {
			found = true
			break
		}
	}
	require.True(t, found, "Associated datasource not found in project's datasources")

	// Remove datasource from project
	err = testQueries.RemoveDatasourceFromProject(context.Background(), RemoveDatasourceFromProjectParams{
		ProjectID:    project.ProjectID,
		DatasourceID: datasource.DatasourceID,
	})
	require.NoError(t, err)

	// Verify datasource is removed
	_, err = testQueries.GetProjectDatasourceAssociation(context.Background(), GetProjectDatasourceAssociationParams{
		ProjectID:    project.ProjectID,
		DatasourceID: datasource.DatasourceID,
	})
	require.Error(t, err)
}
