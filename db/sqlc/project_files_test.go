package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// createRandomProjectFile creates a project file with random values for testing
func createRandomProjectFile(t *testing.T) ProjectFile {
	project := createRandomProject(t)

	arg := UploadProjectFileParams{
		ProjectID: sql.NullInt32{
			Int32: project.ProjectID,
			Valid: true,
		},
		FileName: randomString(10) + ".pdf",
		FileType: "application/pdf",
		FilePath: "/storage/projects/" + randomString(20) + ".pdf",
		FileSize: sql.NullInt32{
			Int32: int32(randomInt(10000, 5000000)), // 10KB to 5MB
			Valid: true,
		},
	}

	file, err := testQueries.UploadProjectFile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, file)

	require.Equal(t, arg.ProjectID, file.ProjectID)
	require.Equal(t, arg.FileName, file.FileName)
	require.Equal(t, arg.FileType, file.FileType)
	require.Equal(t, arg.FilePath, file.FilePath)
	require.Equal(t, arg.FileSize, file.FileSize)

	require.NotZero(t, file.FileID)
	require.NotEmpty(t, file.UploadedAt)

	return file
}

// TestUploadProjectFile tests the UploadProjectFile function
func TestUploadProjectFile(t *testing.T) {
	createRandomProjectFile(t)
}

// TestGetFileByID tests the GetFileByID function
func TestGetFileByID(t *testing.T) {
	file1 := createRandomProjectFile(t)
	file2, err := testQueries.GetFileByID(context.Background(), file1.FileID)
	require.NoError(t, err)
	require.NotEmpty(t, file2)

	require.Equal(t, file1.FileID, file2.FileID)
	require.Equal(t, file1.ProjectID, file2.ProjectID)
	require.Equal(t, file1.FileName, file2.FileName)
	require.Equal(t, file1.FileType, file2.FileType)
	require.Equal(t, file1.FilePath, file2.FilePath)
	require.Equal(t, file1.FileSize, file2.FileSize)
	require.WithinDuration(t, file1.UploadedAt.Time, file2.UploadedAt.Time, 0)
}

// TestListFilesForProject tests the ListFilesForProject function
func TestListFilesForProject(t *testing.T) {
	project := createRandomProject(t)

	// Upload several files for the same project
	for i := 0; i < 5; i++ {
		fileType := []string{"application/pdf", "text/plain", "application/msword", "image/jpeg", "application/vnd.ms-excel"}[i]
		fileExt := []string{".pdf", ".txt", ".doc", ".jpg", ".xls"}[i]

		arg := UploadProjectFileParams{
			ProjectID: sql.NullInt32{
				Int32: project.ProjectID,
				Valid: true,
			},
			FileName: randomString(10) + fileExt,
			FileType: fileType,
			FilePath: "/storage/projects/" + randomString(20) + fileExt,
			FileSize: sql.NullInt32{
				Int32: int32(randomInt(10000, 5000000)), // 10KB to 5MB
				Valid: true,
			},
		}

		_, err := testQueries.UploadProjectFile(context.Background(), arg)
		require.NoError(t, err)
	}

	// Upload some files for other projects
	for i := 0; i < 3; i++ {
		createRandomProjectFile(t)
	}

	arg := ListFilesForProjectParams{
		ProjectID: sql.NullInt32{
			Int32: project.ProjectID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	files, err := testQueries.ListFilesForProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	require.Len(t, files, 5)

	// Verify all files belong to the same project
	for _, file := range files {
		require.Equal(t, project.ProjectID, file.ProjectID.Int32)
	}
}

// TestUpdateProjectFileInfo tests the UpdateProjectFileInfo function
func TestUpdateProjectFileInfo(t *testing.T) {
	file1 := createRandomProjectFile(t)

	arg := UpdateProjectFileInfoParams{
		FileID:   file1.FileID,
		FileName: randomString(10) + ".xlsx",
		FileType: "application/vnd.ms-excel", // Shorter MIME type that fits in VARCHAR(50)
		FilePath: "/storage/projects/updated/" + randomString(20) + ".xlsx",
		FileSize: sql.NullInt32{
			Int32: int32(randomInt(10000, 5000000)), // 10KB to 5MB
			Valid: true,
		},
	}

	file2, err := testQueries.UpdateProjectFileInfo(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, file2)

	require.Equal(t, file1.FileID, file2.FileID)
	require.Equal(t, file1.ProjectID, file2.ProjectID)
	require.Equal(t, arg.FileName, file2.FileName)
	require.Equal(t, arg.FileType, file2.FileType)
	require.Equal(t, arg.FilePath, file2.FilePath)
	require.Equal(t, arg.FileSize, file2.FileSize)
	require.Equal(t, file1.UploadedAt, file2.UploadedAt)
}

// TestDeleteProjectFile tests the DeleteProjectFile function
func TestDeleteProjectFile(t *testing.T) {
	file1 := createRandomProjectFile(t)
	err := testQueries.DeleteProjectFile(context.Background(), file1.FileID)
	require.NoError(t, err)

	file2, err := testQueries.GetFileByID(context.Background(), file1.FileID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, file2)
}
