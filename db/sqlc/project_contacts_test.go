package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAssociateContactWithProject tests the AssociateContactWithProject function
func TestAssociateContactWithProject(t *testing.T) {
	project := createRandomProject(t)
	contact := createRandomContact(t)

	arg := AssociateContactWithProjectParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	association, err := testQueries.AssociateContactWithProject(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, association)

	require.Equal(t, arg.ProjectID, association.ProjectID)
	require.Equal(t, arg.ContactID, association.ContactID)
	require.Equal(t, arg.AssociationNotes, association.AssociationNotes)
}

// TestGetProjectContactAssociation tests the GetProjectContactAssociation function
func TestGetProjectContactAssociation(t *testing.T) {
	project := createRandomProject(t)
	contact := createRandomContact(t)

	createArg := AssociateContactWithProjectParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	association1, err := testQueries.AssociateContactWithProject(context.Background(), createArg)
	require.NoError(t, err)
	require.NotEmpty(t, association1)

	getArg := GetProjectContactAssociationParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
	}

	association2, err := testQueries.GetProjectContactAssociation(context.Background(), getArg)
	require.NoError(t, err)
	require.NotEmpty(t, association2)

	require.Equal(t, association1.ProjectID, association2.ProjectID)
	require.Equal(t, association1.ContactID, association2.ContactID)
	require.Equal(t, association1.AssociationNotes, association2.AssociationNotes)
}

// TestListContactsForProject tests the ListContactsForProject function
func TestListContactsForProject(t *testing.T) {
	project := createRandomProject(t)

	// Associate several contacts with the project
	for i := 0; i < 5; i++ {
		contact := createRandomContact(t)

		arg := AssociateContactWithProjectParams{
			ProjectID: project.ProjectID,
			ContactID: contact.ContactID,
			AssociationNotes: sql.NullString{
				String: randomString(20),
				Valid:  true,
			},
		}

		_, err := testQueries.AssociateContactWithProject(context.Background(), arg)
		require.NoError(t, err)
	}

	listArg := ListContactsForProjectParams{
		ProjectID: project.ProjectID,
		Limit:     10,
		Offset:    0,
	}

	contacts, err := testQueries.ListContactsForProject(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, contacts)
	require.Len(t, contacts, 5)

	// Verify all retrieved contacts are associated with the project
	for _, contact := range contacts {
		require.Equal(t, project.ProjectID, contact.ProjectID)
	}
}

// TestListProjectsForContact tests the ListProjectsForContact function
func TestListProjectsForContact(t *testing.T) {
	contact := createRandomContact(t)

	// Associate the contact with several projects
	for i := 0; i < 5; i++ {
		project := createRandomProject(t)

		arg := AssociateContactWithProjectParams{
			ProjectID: project.ProjectID,
			ContactID: contact.ContactID,
			AssociationNotes: sql.NullString{
				String: randomString(20),
				Valid:  true,
			},
		}

		_, err := testQueries.AssociateContactWithProject(context.Background(), arg)
		require.NoError(t, err)
	}

	listArg := ListProjectsForContactParams{
		ContactID: contact.ContactID,
		Limit:     10,
		Offset:    0,
	}

	projects, err := testQueries.ListProjectsForContact(context.Background(), listArg)
	require.NoError(t, err)
	require.NotEmpty(t, projects)
	require.Len(t, projects, 5)

	// Verify all retrieved projects are associated with the contact
	for _, project := range projects {
		require.Equal(t, contact.ContactID, project.ContactID)
	}
}

// TestUpdateProjectContactAssociation tests the UpdateProjectContactAssociation function
func TestUpdateProjectContactAssociation(t *testing.T) {
	project := createRandomProject(t)
	contact := createRandomContact(t)

	createArg := AssociateContactWithProjectParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	association1, err := testQueries.AssociateContactWithProject(context.Background(), createArg)
	require.NoError(t, err)
	require.NotEmpty(t, association1)

	updateArg := UpdateProjectContactAssociationParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
		AssociationNotes: sql.NullString{
			String: randomString(20),
			Valid:  true,
		},
	}

	association2, err := testQueries.UpdateProjectContactAssociation(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, association2)

	require.Equal(t, association1.ProjectID, association2.ProjectID)
	require.Equal(t, association1.ContactID, association2.ContactID)
	require.Equal(t, updateArg.AssociationNotes, association2.AssociationNotes)
}

// TestRemoveProjectContactAssociation tests the RemoveProjectContactAssociation function
func TestRemoveProjectContactAssociation(t *testing.T) {
	project := createRandomProject(t)
	contact := createRandomContact(t)

	createArg := AssociateContactWithProjectParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
	}

	_, err := testQueries.AssociateContactWithProject(context.Background(), createArg)
	require.NoError(t, err)

	removeArg := RemoveProjectContactAssociationParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
	}

	err = testQueries.RemoveProjectContactAssociation(context.Background(), removeArg)
	require.NoError(t, err)

	getArg := GetProjectContactAssociationParams{
		ProjectID: project.ProjectID,
		ContactID: contact.ContactID,
	}

	_, err = testQueries.GetProjectContactAssociation(context.Background(), getArg)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}
