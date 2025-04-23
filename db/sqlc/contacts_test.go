package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomContact creates a contact with random values for testing
func createRandomContact(t *testing.T) Contact {
	company := createRandomCompany(t)

	arg := CreateContactParams{
		FirstName: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		LastName: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		Email: sql.NullString{
			String: randomEmail(),
			Valid:  true,
		},
		Phone: sql.NullString{
			String: fmt.Sprintf("+1-555-%d", randomInt(1000000, 9999999)),
			Valid:  true,
		},
		LinkedinProfile: sql.NullString{
			String: fmt.Sprintf("https://linkedin.com/in/%s", randomString(10)),
			Valid:  true,
		},
		JobTitle: sql.NullString{
			String: randomString(12),
			Valid:  true,
		},
		CompanyID: sql.NullInt32{
			Int32: company.CompanyID,
			Valid: true,
		},
		Location: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		Bio: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	contact, err := testQueries.CreateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	require.Equal(t, arg.FirstName, contact.FirstName)
	require.Equal(t, arg.LastName, contact.LastName)
	require.Equal(t, arg.Email, contact.Email)
	require.Equal(t, arg.Phone, contact.Phone)
	require.Equal(t, arg.LinkedinProfile, contact.LinkedinProfile)
	require.Equal(t, arg.JobTitle, contact.JobTitle)
	require.Equal(t, arg.CompanyID, contact.CompanyID)
	require.Equal(t, arg.Location, contact.Location)
	require.Equal(t, arg.Bio, contact.Bio)

	require.NotZero(t, contact.ContactID)
	require.NotEmpty(t, contact.ScrapeTimestamp)

	return contact
}

// TestCreateContact tests the CreateContact function
func TestCreateContact(t *testing.T) {
	createRandomContact(t)
}

// TestGetContactByID tests the GetContactByID function
func TestGetContactByID(t *testing.T) {
	contact1 := createRandomContact(t)
	contact2, err := testQueries.GetContactByID(context.Background(), contact1.ContactID)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ContactID, contact2.ContactID)
	require.Equal(t, contact1.FirstName, contact2.FirstName)
	require.Equal(t, contact1.LastName, contact2.LastName)
	require.Equal(t, contact1.Email, contact2.Email)
	require.Equal(t, contact1.Phone, contact2.Phone)
	require.Equal(t, contact1.LinkedinProfile, contact2.LinkedinProfile)
	require.Equal(t, contact1.JobTitle, contact2.JobTitle)
	require.Equal(t, contact1.CompanyID, contact2.CompanyID)
	require.Equal(t, contact1.Location, contact2.Location)
	require.Equal(t, contact1.Bio, contact2.Bio)
	require.WithinDuration(t, contact1.ScrapeTimestamp.Time, contact2.ScrapeTimestamp.Time, time.Second)
}

// TestListContactsByCompany tests the ListContactsByCompany function
func TestListContactsByCompany(t *testing.T) {
	company := createRandomCompany(t)

	// Create several contacts for the same company
	var createdContacts []Contact
	for i := 0; i < 5; i++ {
		arg := CreateContactParams{
			FirstName: sql.NullString{
				String: randomString(8),
				Valid:  true,
			},
			LastName: sql.NullString{
				String: randomString(8),
				Valid:  true,
			},
			Email: sql.NullString{
				String: randomEmail(),
				Valid:  true,
			},
			Phone: sql.NullString{
				String: fmt.Sprintf("+1-555-%d", randomInt(1000000, 9999999)),
				Valid:  true,
			},
			LinkedinProfile: sql.NullString{
				String: fmt.Sprintf("https://linkedin.com/in/%s", randomString(10)),
				Valid:  true,
			},
			JobTitle: sql.NullString{
				String: randomString(12),
				Valid:  true,
			},
			CompanyID: sql.NullInt32{
				Int32: company.CompanyID,
				Valid: true,
			},
			Location: sql.NullString{
				String: randomString(10),
				Valid:  true,
			},
			Bio: sql.NullString{
				String: randomString(30),
				Valid:  true,
			},
		}

		contact, err := testQueries.CreateContact(context.Background(), arg)
		require.NoError(t, err)
		createdContacts = append(createdContacts, contact)
	}

	// Create a few contacts for different companies
	for i := 0; i < 3; i++ {
		createRandomContact(t)
	}

	arg := ListContactsByCompanyParams{
		CompanyID: sql.NullInt32{
			Int32: company.CompanyID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	contacts, err := testQueries.ListContactsByCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contacts)
	require.Len(t, contacts, 5)

	// Verify all contacts belong to the same company
	for _, contact := range contacts {
		require.Equal(t, company.CompanyID, contact.CompanyID.Int32)
	}
}

// TestSearchContactsByName tests the SearchContactsByName function
func TestSearchContactsByName(t *testing.T) {
	// Create a contact with a specific first name for searching
	specificFirstName := "TestSpecific" + randomString(4)

	arg := CreateContactParams{
		FirstName: sql.NullString{
			String: specificFirstName,
			Valid:  true,
		},
		LastName: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		Email: sql.NullString{
			String: randomEmail(),
			Valid:  true,
		},
		CompanyID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
	}

	contact, err := testQueries.CreateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	// Create some random contacts
	for i := 0; i < 5; i++ {
		createRandomContact(t)
	}

	// Search by the specific first name
	searchArg := SearchContactsByNameParams{
		Lower:  "%TestSpecific%",
		Limit:  10,
		Offset: 0,
	}

	searchResults, err := testQueries.SearchContactsByName(context.Background(), searchArg)
	require.NoError(t, err)
	require.NotEmpty(t, searchResults)
	require.GreaterOrEqual(t, len(searchResults), 1)

	found := false
	for _, result := range searchResults {
		if result.FirstName.String == specificFirstName {
			found = true
			break
		}
	}
	require.True(t, found, "Could not find contact with the specific first name")
}

// TestUpdateContact tests the UpdateContact function
func TestUpdateContact(t *testing.T) {
	contact1 := createRandomContact(t)

	arg := UpdateContactParams{
		ContactID: contact1.ContactID,
		FirstName: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		LastName: sql.NullString{
			String: randomString(8),
			Valid:  true,
		},
		Email: sql.NullString{
			String: randomEmail(),
			Valid:  true,
		},
		Phone: sql.NullString{
			String: fmt.Sprintf("+1-555-%d", randomInt(1000000, 9999999)),
			Valid:  true,
		},
		LinkedinProfile: sql.NullString{
			String: fmt.Sprintf("https://linkedin.com/in/%s", randomString(10)),
			Valid:  true,
		},
		JobTitle: sql.NullString{
			String: randomString(12),
			Valid:  true,
		},
		CompanyID: contact1.CompanyID,
		Location: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		Bio: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
	}

	contact2, err := testQueries.UpdateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ContactID, contact2.ContactID)
	require.Equal(t, arg.FirstName, contact2.FirstName)
	require.Equal(t, arg.LastName, contact2.LastName)
	require.Equal(t, arg.Email, contact2.Email)
	require.Equal(t, arg.Phone, contact2.Phone)
	require.Equal(t, arg.LinkedinProfile, contact2.LinkedinProfile)
	require.Equal(t, arg.JobTitle, contact2.JobTitle)
	require.Equal(t, arg.CompanyID, contact2.CompanyID)
	require.Equal(t, arg.Location, contact2.Location)
	require.Equal(t, arg.Bio, contact2.Bio)
	require.WithinDuration(t, contact1.ScrapeTimestamp.Time, contact2.ScrapeTimestamp.Time, time.Second*5)
}

// TestDeleteContact tests the DeleteContact function
func TestDeleteContact(t *testing.T) {
	contact1 := createRandomContact(t)
	err := testQueries.DeleteContact(context.Background(), contact1.ContactID)
	require.NoError(t, err)

	contact2, err := testQueries.GetContactByID(context.Background(), contact1.ContactID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, contact2)
}
