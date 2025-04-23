package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAssociateCompanyWithProject tests the AssociateCompanyWithProject function
func TestAssociateCompanyWithProject(t *testing.T) {
	project := createRandomProject(t)
	company := createRandomCompany(t)

	arg := AssociateCompanyWithProjectParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		MatchingScore: sql.NullString{
			String: "0.85",
			Valid:  true,
		},
		ApproachStrategy: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	association, err := testQueries.AssociateCompanyWithProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, association)

	require.Equal(t, arg.ProjectID, association.ProjectID)
	require.Equal(t, arg.CompanyID, association.CompanyID)
	require.Equal(t, arg.AssociationNotes, association.AssociationNotes)
	require.Equal(t, arg.MatchingScore, association.MatchingScore)
	require.Equal(t, arg.ApproachStrategy, association.ApproachStrategy)
}

// TestGetProjectCompanyAssociation tests the GetProjectCompanyAssociation function
func TestGetProjectCompanyAssociation(t *testing.T) {
	project := createRandomProject(t)
	company := createRandomCompany(t)

	createArg := AssociateCompanyWithProjectParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		MatchingScore: sql.NullString{
			String: "0.75",
			Valid:  true,
		},
		ApproachStrategy: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	association1, err := testQueries.AssociateCompanyWithProject(context.Background(), createArg)
	require.NoError(t, err)
	require.NotEmpty(t, association1)

	getArg := GetProjectCompanyAssociationParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
	}

	association2, err := testQueries.GetProjectCompanyAssociation(context.Background(), getArg)
	require.NoError(t, err)
	require.NotEmpty(t, association2)

	require.Equal(t, association1.ProjectID, association2.ProjectID)
	require.Equal(t, association1.CompanyID, association2.CompanyID)
	require.Equal(t, association1.AssociationNotes, association2.AssociationNotes)
	require.Equal(t, association1.MatchingScore, association2.MatchingScore)
	require.Equal(t, association1.ApproachStrategy, association2.ApproachStrategy)
}

// TestListCompaniesForProject tests the ListCompaniesForProject function
func TestListCompaniesForProject(t *testing.T) {
	project := createRandomProject(t)

	// Associate several companies with the project
	for i := 0; i < 5; i++ {
		company := createRandomCompany(t)
		score := 0.5 + float64(i)*0.1 // Scores from 0.5 to 0.9

		arg := AssociateCompanyWithProjectParams{
			ProjectID: project.ProjectID,
			CompanyID: company.CompanyID,
			MatchingScore: sql.NullString{
				String: formatScore(score),
				Valid:  true,
			},
		}

		_, err := testQueries.AssociateCompanyWithProject(context.Background(), arg)
		require.NoError(t, err)
	}

	listArg := ListCompaniesForProjectParams{
		ProjectID: project.ProjectID,
		Limit:     10,
		Offset:    0,
	}

	companies, err := testQueries.ListCompaniesForProject(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, companies)
	require.Len(t, companies, 5)

	// Verify companies are sorted by matching score (descending)
	for i := 0; i < len(companies)-1; i++ {
		score1 := parseScore(companies[i].MatchingScore.String)
		score2 := parseScore(companies[i+1].MatchingScore.String)
		require.GreaterOrEqual(t, score1, score2)
	}
}

// TestListProjectsForCompany tests the ListProjectsForCompany function
func TestListProjectsForCompany(t *testing.T) {
	company := createRandomCompany(t)

	// Associate the company with several projects
	for i := 0; i < 5; i++ {
		project := createRandomProject(t)
		score := 0.5 + float64(i)*0.1 // Scores from 0.5 to 0.9

		arg := AssociateCompanyWithProjectParams{
			ProjectID: project.ProjectID,
			CompanyID: company.CompanyID,
			MatchingScore: sql.NullString{
				String: formatScore(score),
				Valid:  true,
			},
		}

		_, err := testQueries.AssociateCompanyWithProject(context.Background(), arg)
		require.NoError(t, err)
	}

	listArg := ListProjectsForCompanyParams{
		CompanyID: company.CompanyID,
		Limit:     10,
		Offset:    0,
	}

	projects, err := testQueries.ListProjectsForCompany(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, projects)
	require.Len(t, projects, 5)

	// Verify projects are sorted by matching score (descending)
	for i := 0; i < len(projects)-1; i++ {
		score1 := parseScore(projects[i].MatchingScore.String)
		score2 := parseScore(projects[i+1].MatchingScore.String)
		require.GreaterOrEqual(t, score1, score2)
	}
}

// TestUpdateProjectCompanyAssociation tests the UpdateProjectCompanyAssociation function
func TestUpdateProjectCompanyAssociation(t *testing.T) {
	project := createRandomProject(t)
	company := createRandomCompany(t)

	createArg := AssociateCompanyWithProjectParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		MatchingScore: sql.NullString{
			String: "0.65",
			Valid:  true,
		},
		ApproachStrategy: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	association1, err := testQueries.AssociateCompanyWithProject(context.Background(), createArg)
	require.NoError(t, err)
	require.NotEmpty(t, association1)

	updateArg := UpdateProjectCompanyAssociationParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
		MatchingScore: sql.NullString{
			String: "0.90",
			Valid:  true,
		},
		ApproachStrategy: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	association2, err := testQueries.UpdateProjectCompanyAssociation(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, association2)

	require.Equal(t, association1.ProjectID, association2.ProjectID)
	require.Equal(t, association1.CompanyID, association2.CompanyID)
	require.Equal(t, updateArg.AssociationNotes, association2.AssociationNotes)
	require.Equal(t, updateArg.MatchingScore, association2.MatchingScore)
	require.Equal(t, updateArg.ApproachStrategy, association2.ApproachStrategy)
}

// TestRemoveProjectCompanyAssociation tests the RemoveProjectCompanyAssociation function
func TestRemoveProjectCompanyAssociation(t *testing.T) {
	project := createRandomProject(t)
	company := createRandomCompany(t)

	createArg := AssociateCompanyWithProjectParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
	}

	_, err := testQueries.AssociateCompanyWithProject(context.Background(), createArg)
	require.NoError(t, err)

	removeArg := RemoveProjectCompanyAssociationParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
	}

	err = testQueries.RemoveProjectCompanyAssociation(context.Background(), removeArg)
	require.NoError(t, err)

	getArg := GetProjectCompanyAssociationParams{
		ProjectID: project.ProjectID,
		CompanyID: company.CompanyID,
	}

	_, err = testQueries.GetProjectCompanyAssociation(context.Background(), getArg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

// Helper functions for score formatting and parsing
func formatScore(score float64) string {
	return fmt.Sprintf("%.2f", score)
}

func parseScore(scoreStr string) float64 {
	var score float64
	_, err := fmt.Sscanf(scoreStr, "%f", &score)
	if err != nil {
		return 0
	}
	return score
}
