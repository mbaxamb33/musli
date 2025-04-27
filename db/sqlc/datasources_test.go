package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomDatasource(t *testing.T) Datasource {
	// Create a datasource with various types of sources
	datasourceTypes := []DatasourceType{
		DatasourceTypePdf,
		DatasourceTypeWebsite,
		DatasourceTypeExcel,
		DatasourceTypeMp3,
		DatasourceTypeWordDocument,
		DatasourceTypePowerpoint,
		DatasourceTypePlainText,
	}

	// Randomly select a datasource type
	sourceType := datasourceTypes[randomInt(0, int64(len(datasourceTypes)-1))]

	arg := CreateDatasourceParams{
		SourceType: sourceType,
		Link:       sql.NullString{String: "https://example.com/" + randomString(10), Valid: true},
		FileName:   sql.NullString{String: randomString(8) + getFileExtension(sourceType), Valid: true},
		FileData:   []byte(randomString(100)), // Some sample file content
	}

	datasource, err := testQueries.CreateDatasource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, datasource)

	require.Equal(t, arg.SourceType, datasource.SourceType)
	require.Equal(t, arg.Link, datasource.Link)
	require.Equal(t, arg.FileName, datasource.FileName)
	require.Equal(t, arg.FileData, datasource.FileData)
	require.NotZero(t, datasource.DatasourceID)
	require.NotZero(t, datasource.CreatedAt)

	return datasource
}

// Helper function to get file extension based on datasource type
func getFileExtension(sourceType DatasourceType) string {
	switch sourceType {
	case DatasourceTypePdf:
		return ".pdf"
	case DatasourceTypeWebsite:
		return ".html"
	case DatasourceTypeExcel:
		return ".xlsx"
	case DatasourceTypeMp3:
		return ".mp3"
	case DatasourceTypeWordDocument:
		return ".docx"
	case DatasourceTypePowerpoint:
		return ".pptx"
	case DatasourceTypePlainText:
		return ".txt"
	default:
		return ".unknown"
	}
}

func TestCreateDatasource(t *testing.T) {
	createRandomDatasource(t)
}

func TestGetDatasourceByID(t *testing.T) {
	// Create a random datasource
	datasource1 := createRandomDatasource(t)

	// Retrieve the datasource by ID
	datasource2, err := testQueries.GetDatasourceByID(context.Background(), datasource1.DatasourceID)
	require.NoError(t, err)
	require.NotEmpty(t, datasource2)

	require.Equal(t, datasource1.DatasourceID, datasource2.DatasourceID)
	require.Equal(t, datasource1.SourceType, datasource2.SourceType)
	require.Equal(t, datasource1.Link, datasource2.Link)
	require.Equal(t, datasource1.FileName, datasource2.FileName)
	require.WithinDuration(t, datasource1.CreatedAt.Time, datasource2.CreatedAt.Time, time.Second)
}

func TestDeleteDatasource(t *testing.T) {
	// Create a random datasource
	datasource1 := createRandomDatasource(t)

	// Delete the datasource
	err := testQueries.DeleteDatasource(context.Background(), datasource1.DatasourceID)
	require.NoError(t, err)

	// Try to retrieve the deleted datasource (should fail)
	_, err = testQueries.GetDatasourceByID(context.Background(), datasource1.DatasourceID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListDatasources(t *testing.T) {
	// Create multiple datasources
	expectedDatasources := 10
	for i := 0; i < expectedDatasources; i++ {
		createRandomDatasource(t)
	}

	// List datasources with pagination
	arg := ListDatasourcesParams{
		Limit:  5,
		Offset: 0,
	}

	datasources, err := testQueries.ListDatasources(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, datasources, 5)

	for _, datasource := range datasources {
		require.NotEmpty(t, datasource)
	}
}

func TestListDatasourcesByType(t *testing.T) {
	// Choose a specific datasource type to test
	sourceType := DatasourceTypePdf

	// Create multiple PDF datasources
	expectedDatasources := 10
	for i := 0; i < expectedDatasources; i++ {
		arg := CreateDatasourceParams{
			SourceType: sourceType,
			Link:       sql.NullString{String: "https://example.com/" + randomString(10), Valid: true},
			FileName:   sql.NullString{String: randomString(8) + ".pdf", Valid: true},
			FileData:   []byte(randomString(100)),
		}
		_, err := testQueries.CreateDatasource(context.Background(), arg)
		require.NoError(t, err)
	}

	// List datasources by type with pagination
	arg := ListDatasourcesByTypeParams{
		SourceType: sourceType,
		Limit:      5,
		Offset:     0,
	}

	datasources, err := testQueries.ListDatasourcesByType(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, datasources, 5)

	for _, datasource := range datasources {
		require.NotEmpty(t, datasource)
		require.Equal(t, sourceType, datasource.SourceType)
	}
}

func TestDatasourceProjectAssociation(t *testing.T) {
	// Create a random datasource
	datasource := createRandomDatasource(t)

	// Create a random project
	project := createRandomProject(t)

	// Associate datasource with project
	err := testQueries.AssociateDatasourceWithProject(context.Background(), AssociateDatasourceWithProjectParams{
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

	// List projects for the datasource
	projectsArg := ListProjectsByDatasourceParams{
		DatasourceID: datasource.DatasourceID,
		Limit:        5,
		Offset:       0,
	}
	datasourceProjects, err := testQueries.ListProjectsByDatasource(context.Background(), projectsArg)
	require.NoError(t, err)
	require.NotEmpty(t, datasourceProjects)

	// Verify the project is in the list
	found := false
	for _, proj := range datasourceProjects {
		if proj.ProjectID == project.ProjectID {
			found = true
			break
		}
	}
	require.True(t, found, "Associated project not found in datasource's projects")

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

func TestCreateDatasourceWithAllTypes(t *testing.T) {
	datasourceTypes := []DatasourceType{
		DatasourceTypePdf,
		DatasourceTypeWebsite,
		DatasourceTypeExcel,
		DatasourceTypeMp3,
		DatasourceTypeWordDocument,
		DatasourceTypePowerpoint,
		DatasourceTypePlainText,
	}

	for _, sourceType := range datasourceTypes {
		t.Run(string(sourceType), func(t *testing.T) {
			arg := CreateDatasourceParams{
				SourceType: sourceType,
				Link:       sql.NullString{String: "https://example.com/" + randomString(10), Valid: true},
				FileName:   sql.NullString{String: randomString(8) + getFileExtension(sourceType), Valid: true},
				FileData:   []byte(randomString(100)),
			}

			datasource, err := testQueries.CreateDatasource(context.Background(), arg)
			require.NoError(t, err)
			require.NotEmpty(t, datasource)

			require.Equal(t, sourceType, datasource.SourceType)
			require.Equal(t, arg.Link, datasource.Link)
			require.Equal(t, arg.FileName, datasource.FileName)
			require.Equal(t, arg.FileData, datasource.FileData)
		})
	}
}
