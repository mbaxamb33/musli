package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProjectCreate tests creating a new project in the database
func TestProjectCreate(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// First create a user to be the project owner
	userID := createTestUser(t, db)

	// Generate test data
	projectName := "TestProject_" + randomString(8)
	description := sql.NullString{String: "A test project for database testing", Valid: true}

	// Create dates
	now := time.Now()
	startDate := sql.NullTime{Time: now, Valid: true}
	endDate := sql.NullTime{Time: now.AddDate(0, 3, 0), Valid: true} // 3 months later
	status := sql.NullString{String: "active", Valid: true}

	// Execute
	query := `
	INSERT INTO projects (user_id, name, description, start_date, end_date, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING project_id, user_id, name, description, start_date, end_date, status, created_at, updated_at
	`

	row := db.QueryRowContext(
		context.Background(),
		query,
		userID, projectName, description, startDate, endDate, status,
	)

	// Verify
	var project struct {
		ProjectID   int32
		UserID      int32
		Name        string
		Description sql.NullString
		StartDate   sql.NullTime
		EndDate     sql.NullTime
		Status      sql.NullString
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
	}

	err = row.Scan(
		&project.ProjectID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.StartDate,
		&project.EndDate,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	require.NoError(t, err)

	// Assert
	assert.NotZero(t, project.ProjectID)
	assert.Equal(t, userID, project.UserID)
	assert.Equal(t, projectName, project.Name)
	assert.Equal(t, description, project.Description)
	// Don't compare times directly
	assert.True(t, project.StartDate.Valid)
	assert.True(t, project.EndDate.Valid)
	assert.Equal(t, status, project.Status)
	assert.True(t, project.CreatedAt.Valid)
	assert.True(t, project.UpdatedAt.Valid)
	assert.NotZero(t, project.CreatedAt.Time)
	assert.NotZero(t, project.UpdatedAt.Time)

	// Log success
	t.Logf("Successfully created project: ID=%d, Name=%s", project.ProjectID, project.Name)

	// Clean up
	cleanupTestProject(t, db, project.ProjectID)
	cleanupTestUser(t, db, userID)
}

// TestProjectGet tests retrieving a project by ID
func TestProjectGet(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a test user and project
	userID := createTestUser(t, db)
	projectID := createTestProject(t, db, userID)

	// Get the project
	query := "SELECT * FROM projects WHERE project_id = $1"
	var project struct {
		ProjectID   int32
		UserID      int32
		Name        string
		Description sql.NullString
		StartDate   sql.NullTime
		EndDate     sql.NullTime
		Status      sql.NullString
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
	}

	err = db.QueryRowContext(context.Background(), query, projectID).Scan(
		&project.ProjectID,
		&project.UserID,
		&project.Name,
		&project.Description,
		&project.StartDate,
		&project.EndDate,
		&project.Status,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	require.NoError(t, err)

	// Verify
	assert.Equal(t, projectID, project.ProjectID)
	assert.Equal(t, userID, project.UserID)
	assert.NotEmpty(t, project.Name)
	assert.True(t, project.Description.Valid)
	assert.NotEmpty(t, project.Description.String)
	assert.True(t, project.StartDate.Valid)
	assert.True(t, project.EndDate.Valid)
	assert.True(t, project.Status.Valid)
	assert.Equal(t, "active", project.Status.String)

	// Clean up
	cleanupTestProject(t, db, projectID)
	cleanupTestUser(t, db, userID)
}

// TestProjectUpdate tests updating a project's information
func TestProjectUpdate(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a test user and project
	userID := createTestUser(t, db)
	projectID := createTestProject(t, db, userID)

	// Get original project data
	var originalUpdatedAt time.Time
	err = db.QueryRowContext(
		context.Background(),
		"SELECT updated_at FROM projects WHERE project_id = $1",
		projectID,
	).Scan(&originalUpdatedAt)
	require.NoError(t, err)

	// Wait a moment to ensure updated_at timestamp will be different
	time.Sleep(10 * time.Millisecond)

	// Update project
	newName := "Updated Project Name " + randomString(5)
	newDescription := sql.NullString{String: "Updated project description", Valid: true}
	newStatus := sql.NullString{String: "completed", Valid: true}

	// Create new end date (2 months from now)
	now := time.Now()
	newEndDate := sql.NullTime{
		Time:  now.AddDate(0, 2, 0),
		Valid: true,
	}

	query := `
	UPDATE projects
	SET 
		name = $1,
		description = $2,
		end_date = $3,
		status = $4,
		updated_at = CURRENT_TIMESTAMP
	WHERE project_id = $5
	RETURNING project_id, user_id, name, description, start_date, end_date, status, created_at, updated_at
	`

	var updatedProject struct {
		ProjectID   int32
		UserID      int32
		Name        string
		Description sql.NullString
		StartDate   sql.NullTime
		EndDate     sql.NullTime
		Status      sql.NullString
		CreatedAt   sql.NullTime
		UpdatedAt   sql.NullTime
	}

	err = db.QueryRowContext(
		context.Background(),
		query,
		newName, newDescription, newEndDate, newStatus, projectID,
	).Scan(
		&updatedProject.ProjectID,
		&updatedProject.UserID,
		&updatedProject.Name,
		&updatedProject.Description,
		&updatedProject.StartDate,
		&updatedProject.EndDate,
		&updatedProject.Status,
		&updatedProject.CreatedAt,
		&updatedProject.UpdatedAt,
	)
	require.NoError(t, err)

	// Verify fields were updated
	assert.Equal(t, projectID, updatedProject.ProjectID)
	assert.Equal(t, userID, updatedProject.UserID)
	assert.Equal(t, newName, updatedProject.Name)
	assert.Equal(t, newDescription, updatedProject.Description)
	// Don't compare end date directly
	assert.True(t, updatedProject.EndDate.Valid)
	assert.Equal(t, newStatus, updatedProject.Status)
	assert.True(t, updatedProject.UpdatedAt.Time.After(originalUpdatedAt))

	// Clean up
	cleanupTestProject(t, db, projectID)
	cleanupTestUser(t, db, userID)
}

// TestProjectStatusUpdate tests updating just the status of a project
func TestProjectStatusUpdate(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a test user and project
	userID := createTestUser(t, db)
	projectID := createTestProject(t, db, userID)

	// Update just the project status
	newStatus := sql.NullString{String: "on_hold", Valid: true}

	query := `
	UPDATE projects
	SET 
		status = $1,
		updated_at = CURRENT_TIMESTAMP
	WHERE project_id = $2
	RETURNING project_id, status, updated_at
	`

	var updatedProject struct {
		ProjectID int32
		Status    sql.NullString
		UpdatedAt sql.NullTime
	}

	err = db.QueryRowContext(
		context.Background(),
		query,
		newStatus, projectID,
	).Scan(
		&updatedProject.ProjectID,
		&updatedProject.Status,
		&updatedProject.UpdatedAt,
	)
	require.NoError(t, err)

	// Verify status was updated
	assert.Equal(t, projectID, updatedProject.ProjectID)
	assert.Equal(t, newStatus, updatedProject.Status)

	// Clean up
	cleanupTestProject(t, db, projectID)
	cleanupTestUser(t, db, userID)
}

// TestProjectList tests listing and pagination of projects
func TestProjectList(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a test user
	userID := createTestUser(t, db)

	// Create multiple test projects
	projectCount := 5
	projectIDs := make([]int32, projectCount)

	for i := 0; i < projectCount; i++ {
		projectIDs[i] = createTestProject(t, db, userID)
	}

	// Test listing all projects
	rows, err := db.QueryContext(
		context.Background(),
		"SELECT project_id, user_id, name FROM projects WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	require.NoError(t, err)
	defer rows.Close()

	var projects []struct {
		ProjectID int32
		UserID    int32
		Name      string
	}

	for rows.Next() {
		var project struct {
			ProjectID int32
			UserID    int32
			Name      string
		}
		err := rows.Scan(&project.ProjectID, &project.UserID, &project.Name)
		require.NoError(t, err)
		projects = append(projects, project)
	}
	require.NoError(t, rows.Err())

	// Verify we have at least the number of projects we created
	assert.GreaterOrEqual(t, len(projects), projectCount)

	// Verify all returned projects belong to our test user
	for _, project := range projects {
		assert.Equal(t, userID, project.UserID)
	}

	// Test pagination
	limit := 3
	offset := 0

	paginatedRows, err := db.QueryContext(
		context.Background(),
		"SELECT project_id FROM projects WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset,
	)
	require.NoError(t, err)
	defer paginatedRows.Close()

	var paginatedProjects []int32
	for paginatedRows.Next() {
		var id int32
		err := paginatedRows.Scan(&id)
		require.NoError(t, err)
		paginatedProjects = append(paginatedProjects, id)
	}
	require.NoError(t, paginatedRows.Err())

	// Verify pagination returned the correct number of projects
	assert.Equal(t, limit, len(paginatedProjects))

	// Clean up
	for _, id := range projectIDs {
		cleanupTestProject(t, db, id)
	}
	cleanupTestUser(t, db, userID)
}

// TestProjectResourceCount tests getting projects with their resource counts
func TestProjectResourceCount(t *testing.T) {
	// Setup
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/musli?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a test user and project
	userID := createTestUser(t, db)
	projectID := createTestProject(t, db, userID)

	// Create some test resources
	resourceCount := 3
	var resourceIDs []int32

	for i := 0; i < resourceCount; i++ {
		var resourceID int32
		resourceName := fmt.Sprintf("TestResource_%d_%s", i, randomString(5))

		err := db.QueryRowContext(
			context.Background(),
			`INSERT INTO resources (name, description) 
			VALUES ($1, $2) RETURNING resource_id`,
			resourceName, sql.NullString{String: "Test resource description", Valid: true},
		).Scan(&resourceID)

		require.NoError(t, err)
		resourceIDs = append(resourceIDs, resourceID)

		// Link resources to project
		_, err = db.ExecContext(
			context.Background(),
			`INSERT INTO project_resources (project_id, resource_id, quantity) 
			VALUES ($1, $2, $3)`,
			projectID, resourceID, "1.0",
		)
		require.NoError(t, err)
	}

	// Query projects with resource count
	query := `
	SELECT 
		p.project_id, 
		COUNT(pr.project_resource_id) AS resource_count
	FROM projects p
	LEFT JOIN project_resources pr ON p.project_id = pr.project_id
	WHERE p.project_id = $1
	GROUP BY p.project_id
	`

	var resultProject struct {
		ProjectID     int32
		ResourceCount int64
	}

	err = db.QueryRowContext(
		context.Background(),
		query,
		projectID,
	).Scan(
		&resultProject.ProjectID,
		&resultProject.ResourceCount,
	)
	require.NoError(t, err)

	// Verify resource count
	assert.Equal(t, projectID, resultProject.ProjectID)
	assert.Equal(t, int64(resourceCount), resultProject.ResourceCount)

	// Clean up
	for _, resourceID := range resourceIDs {
		_, err = db.ExecContext(
			context.Background(),
			"DELETE FROM project_resources WHERE project_id = $1 AND resource_id = $2",
			projectID, resourceID,
		)
		require.NoError(t, err)

		_, err = db.ExecContext(
			context.Background(),
			"DELETE FROM resources WHERE resource_id = $1",
			resourceID,
		)
		require.NoError(t, err)
	}

	cleanupTestProject(t, db, projectID)
	cleanupTestUser(t, db, userID)
}

// Helper function to create a test user
func createTestUser(t *testing.T, db *sql.DB) int32 {
	username := "testuser_" + randomString(8)
	email := username + "@test.com"
	passwordHash := "hashed_password_" + randomString(10)

	var userID int32
	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3) RETURNING user_id`,
		username, email, passwordHash,
	).Scan(&userID)

	require.NoError(t, err)
	require.NotZero(t, userID)

	return userID
}

// Helper function to create a test project
func createTestProject(t *testing.T, db *sql.DB, userID int32) int32 {
	projectName := "TestProject_" + randomString(8)
	description := sql.NullString{String: "Test project description", Valid: true}

	// Create dates
	now := time.Now()
	startDate := sql.NullTime{Time: now, Valid: true}
	endDate := sql.NullTime{Time: now.AddDate(0, 3, 0), Valid: true} // 3 months later
	status := sql.NullString{String: "active", Valid: true}

	var projectID int32
	err := db.QueryRowContext(
		context.Background(),
		`INSERT INTO projects (user_id, name, description, start_date, end_date, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING project_id`,
		userID, projectName, description, startDate, endDate, status,
	).Scan(&projectID)

	require.NoError(t, err)
	require.NotZero(t, projectID)

	return projectID
}

// Helper function to clean up a test project
func cleanupTestProject(t *testing.T, db *sql.DB, projectID int32) {
	_, err := db.ExecContext(
		context.Background(),
		"DELETE FROM projects WHERE project_id = $1",
		projectID,
	)
	require.NoError(t, err)
}

// Helper function to clean up a test user
func cleanupTestUser(t *testing.T, db *sql.DB, userID int32) {
	_, err := db.ExecContext(
		context.Background(),
		"DELETE FROM users WHERE user_id = $1",
		userID,
	)
	require.NoError(t, err)
}
