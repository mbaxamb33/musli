package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomContact(t *testing.T) Contact {
	// First, create a random company to associate the contact with
	company := createRandomCompany(t)

	arg := CreateContactParams{
		CompanyID: company.CompanyID,
		FirstName: randomString(6),
		LastName:  randomString(8),
		Position:  sql.NullString{String: randomString(10), Valid: true},
		Email:     sql.NullString{String: randomEmail(), Valid: true},
		Phone:     sql.NullString{String: "+" + randomString(11), Valid: true},
		Notes:     sql.NullString{String: "Test contact notes", Valid: true},
	}

	contact, err := testQueries.CreateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	require.Equal(t, arg.CompanyID, contact.CompanyID)
	require.Equal(t, arg.FirstName, contact.FirstName)
	require.Equal(t, arg.LastName, contact.LastName)
	require.Equal(t, arg.Position, contact.Position)
	require.Equal(t, arg.Email, contact.Email)
	require.Equal(t, arg.Phone, contact.Phone)
	require.Equal(t, arg.Notes, contact.Notes)
	require.NotZero(t, contact.ContactID)
	require.NotZero(t, contact.CreatedAt)

	return contact
}

func TestCreateContact(t *testing.T) {
	createRandomContact(t)
}

func TestGetContactByID(t *testing.T) {
	// Create a random contact
	contact1 := createRandomContact(t)

	// Retrieve the contact by ID
	contact2, err := testQueries.GetContactByID(context.Background(), contact1.ContactID)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ContactID, contact2.ContactID)
	require.Equal(t, contact1.CompanyID, contact2.CompanyID)
	require.Equal(t, contact1.FirstName, contact2.FirstName)
	require.Equal(t, contact1.LastName, contact2.LastName)
	require.Equal(t, contact1.Position, contact2.Position)
	require.Equal(t, contact1.Email, contact2.Email)
	require.Equal(t, contact1.Phone, contact2.Phone)
	require.Equal(t, contact1.Notes, contact2.Notes)
	require.WithinDuration(t, contact1.CreatedAt.Time, contact2.CreatedAt.Time, time.Second)
}

func TestUpdateContact(t *testing.T) {
	// Create a random contact
	contact1 := createRandomContact(t)

	// Prepare updated contact details
	arg := UpdateContactParams{
		ContactID: contact1.ContactID,
		FirstName: randomString(7),
		LastName:  randomString(9),
		Position:  sql.NullString{String: randomString(12), Valid: true},
		Email:     sql.NullString{String: randomEmail(), Valid: true},
		Phone:     sql.NullString{String: "+" + randomString(12), Valid: true},
		Notes:     sql.NullString{String: "Updated test contact notes", Valid: true},
	}

	// Update the contact
	contact2, err := testQueries.UpdateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	// Verify updated details
	require.Equal(t, contact1.ContactID, contact2.ContactID)
	require.Equal(t, contact1.CompanyID, contact2.CompanyID)
	require.Equal(t, arg.FirstName, contact2.FirstName)
	require.Equal(t, arg.LastName, contact2.LastName)
	require.Equal(t, arg.Position, contact2.Position)
	require.Equal(t, arg.Email, contact2.Email)
	require.Equal(t, arg.Phone, contact2.Phone)
	require.Equal(t, arg.Notes, contact2.Notes)
}

func TestDeleteContact(t *testing.T) {
	// Create a random contact
	contact1 := createRandomContact(t)

	// Delete the contact
	err := testQueries.DeleteContact(context.Background(), contact1.ContactID)
	require.NoError(t, err)

	// Try to retrieve the deleted contact (should fail)
	_, err = testQueries.GetContactByID(context.Background(), contact1.ContactID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
}

func TestListContactsByCompanyID(t *testing.T) {
	// Create a random company
	company := createRandomCompany(t)

	// Create multiple contacts for the company
	expectedContacts := 10
	for i := 0; i < expectedContacts; i++ {
		arg := CreateContactParams{
			CompanyID: company.CompanyID,
			FirstName: randomString(6),
			LastName:  randomString(8),
			Position:  sql.NullString{String: randomString(10), Valid: true},
		}
		_, err := testQueries.CreateContact(context.Background(), arg)
		require.NoError(t, err)
	}

	// List contacts with pagination
	arg := ListContactsByCompanyIDParams{
		CompanyID: company.CompanyID,
		Limit:     5,
		Offset:    0,
	}

	contacts, err := testQueries.ListContactsByCompanyID(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, contacts, 5)

	for _, contact := range contacts {
		require.NotEmpty(t, contact)
		require.Equal(t, company.CompanyID, contact.CompanyID)
	}
}

func TestSearchContactsByName(t *testing.T) {
	// Create multiple contacts
	contacts := make([]Contact, 5)
	for i := 0; i < 5; i++ {
		contacts[i] = createRandomContact(t)
	}

	// Choose a contact to search for (use part of first or last name)
	searchContact := contacts[2]
	searchTerm := searchContact.FirstName[:3]

	// Search contacts by name
	arg := SearchContactsByNameParams{
		Column1: sql.NullString{String: searchTerm, Valid: true},
		Limit:   5,
		Offset:  0,
	}

	foundContacts, err := testQueries.SearchContactsByName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, foundContacts)

	// Verify that the search results contain the searched contact
	found := false
	for _, contact := range foundContacts {
		if contact.ContactID == searchContact.ContactID {
			found = true
			break
		}
	}
	require.True(t, found, "Search term did not return the expected contact")
}

func TestSearchContactsByCompanyAndName(t *testing.T) {
	// Create a random company
	company := createRandomCompany(t)

	// Create multiple contacts for the company
	contacts := make([]Contact, 5)
	for i := 0; i < 5; i++ {
		arg := CreateContactParams{
			CompanyID: company.CompanyID,
			FirstName: randomString(6),
			LastName:  randomString(8),
			Position:  sql.NullString{String: randomString(10), Valid: true},
		}
		contacts[i], _ = testQueries.CreateContact(context.Background(), arg)
	}

	// Choose a contact to search for
	searchContact := contacts[2]
	searchTerm := searchContact.FirstName[:3]

	// Search contacts by company and name
	arg := SearchContactsByCompanyAndNameParams{
		CompanyID: company.CompanyID,
		Column2:   sql.NullString{String: searchTerm, Valid: true},
		Limit:     5,
		Offset:    0,
	}

	foundContacts, err := testQueries.SearchContactsByCompanyAndName(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, foundContacts)

	// Verify search results
	found := false
	for _, contact := range foundContacts {
		require.Equal(t, company.CompanyID, contact.CompanyID)
		if contact.ContactID == searchContact.ContactID {
			found = true
		}
	}
	require.True(t, found, "Search term did not return the expected contact")
}
