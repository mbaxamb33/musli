package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomContactNews creates contact news with random values for testing
func createRandomContactNews(t *testing.T) ContactNews {
	contact := createRandomContact(t)
	datasource := createRandomDatasource(t)

	pubDate := time.Now().AddDate(0, -1, 0) // One month ago

	arg := CreateContactNewsParams{
		ContactID: sql.NullInt32{
			Int32: contact.ContactID,
			Valid: true,
		},
		Title: randomString(15),
		PublicationDate: sql.NullTime{
			Time:  pubDate,
			Valid: true,
		},
		Source: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		Url: sql.NullString{
			String: "https://" + randomString(8) + ".com/news/" + randomString(10),
			Valid:  true,
		},
		Summary: sql.NullString{
			String: randomString(50),
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	news, err := testQueries.CreateContactNews(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, news)

	require.Equal(t, arg.ContactID, news.ContactID)
	require.Equal(t, arg.Title, news.Title)
	require.Equal(t, arg.PublicationDate.Time.Day(), news.PublicationDate.Time.Day())
	require.Equal(t, arg.PublicationDate.Time.Month(), news.PublicationDate.Time.Month())
	require.Equal(t, arg.PublicationDate.Time.Year(), news.PublicationDate.Time.Year())
	require.Equal(t, arg.Source, news.Source)
	require.Equal(t, arg.Url, news.Url)
	require.Equal(t, arg.Summary, news.Summary)
	require.Equal(t, arg.DatasourceID, news.DatasourceID)

	require.NotZero(t, news.MentionID)

	return news
}

// TestCreateContactNews tests the CreateContactNews function
func TestCreateContactNews(t *testing.T) {
	createRandomContactNews(t)
}

// TestGetContactNewsByID tests the GetContactNewsByID function
func TestGetContactNewsByID(t *testing.T) {
	news1 := createRandomContactNews(t)
	news2, err := testQueries.GetContactNewsByID(context.Background(), news1.MentionID)
	require.NoError(t, err)
	require.NotEmpty(t, news2)

	require.Equal(t, news1.MentionID, news2.MentionID)
	require.Equal(t, news1.ContactID, news2.ContactID)
	require.Equal(t, news1.Title, news2.Title)
	require.Equal(t, news1.PublicationDate.Time.Format("2006-01-02"),
		news2.PublicationDate.Time.Format("2006-01-02"))
	require.Equal(t, news1.Source, news2.Source)
	require.Equal(t, news1.Url, news2.Url)
	require.Equal(t, news1.Summary, news2.Summary)
	require.Equal(t, news1.DatasourceID, news2.DatasourceID)
}

// TestListContactNewsByContact tests the ListContactNewsByContact function
func TestListContactNewsByContact(t *testing.T) {
	contact := createRandomContact(t)

	// Create several news items for the same contact
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateContactNewsParams{
			ContactID: sql.NullInt32{
				Int32: contact.ContactID,
				Valid: true,
			},
			Title: randomString(15),
			PublicationDate: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			Source: sql.NullString{
				String: randomString(10),
				Valid:  true,
			},
			Url: sql.NullString{
				String: "https://" + randomString(8) + ".com/news/" + randomString(10),
				Valid:  true,
			},
			Summary: sql.NullString{
				String: randomString(50),
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateContactNews(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some news for other contacts
	for i := 0; i < 3; i++ {
		createRandomContactNews(t)
	}

	arg := ListContactNewsByContactParams{
		ContactID: sql.NullInt32{
			Int32: contact.ContactID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListContactNewsByContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.Len(t, newsList, 5)

	// Verify all news items belong to the same contact
	for _, news := range newsList {
		require.Equal(t, contact.ContactID, news.ContactID.Int32)
	}

	// Verify they're sorted by publication date (descending)
	for i := 0; i < len(newsList)-1; i++ {
		require.True(t, newsList[i].PublicationDate.Time.After(newsList[i+1].PublicationDate.Time) ||
			newsList[i].PublicationDate.Time.Equal(newsList[i+1].PublicationDate.Time))
	}
}

// TestListContactNewsByDatasource tests the ListContactNewsByDatasource function
func TestListContactNewsByDatasource(t *testing.T) {
	datasource := createRandomDatasource(t)

	// Create several news items from the same datasource
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateContactNewsParams{
			ContactID: sql.NullInt32{
				Int32: createRandomContact(t).ContactID,
				Valid: true,
			},
			Title: randomString(15),
			PublicationDate: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			Source: sql.NullString{
				String: randomString(10),
				Valid:  true,
			},
			Url: sql.NullString{
				String: "https://" + randomString(8) + ".com/news/" + randomString(10),
				Valid:  true,
			},
			Summary: sql.NullString{
				String: randomString(50),
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: datasource.DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateContactNews(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListContactNewsByDatasourceParams{
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListContactNewsByDatasource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.Len(t, newsList, 5)

	// Verify all news items come from the same datasource
	for _, news := range newsList {
		require.Equal(t, datasource.DatasourceID, news.DatasourceID.Int32)
	}
}

// TestListContactNewsBySource tests the ListContactNewsBySource function
func TestListContactNewsBySource(t *testing.T) {
	sourceIdentifier := "TestSource_" + randomString(5)

	// Create several news items with the same source
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateContactNewsParams{
			ContactID: sql.NullInt32{
				Int32: createRandomContact(t).ContactID,
				Valid: true,
			},
			Title: randomString(15),
			PublicationDate: sql.NullTime{
				Time:  pubDate,
				Valid: true,
			},
			Source: sql.NullString{
				String: sourceIdentifier,
				Valid:  true,
			},
			Url: sql.NullString{
				String: "https://" + randomString(8) + ".com/news/" + randomString(10),
				Valid:  true,
			},
			Summary: sql.NullString{
				String: randomString(50),
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateContactNews(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some news with other sources
	for i := 0; i < 3; i++ {
		createRandomContactNews(t)
	}

	arg := ListContactNewsBySourceParams{
		Source: sql.NullString{
			String: sourceIdentifier,
			Valid:  true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListContactNewsBySource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.Len(t, newsList, 5)

	// Verify all news items have the same source
	for _, news := range newsList {
		require.Equal(t, sourceIdentifier, news.Source.String)
	}
}

// TestUpdateContactNews tests the UpdateContactNews function
func TestUpdateContactNews(t *testing.T) {
	news1 := createRandomContactNews(t)

	pubDate := time.Now().AddDate(0, -2, 0) // Two months ago

	arg := UpdateContactNewsParams{
		MentionID: news1.MentionID,
		Title:     randomString(15),
		PublicationDate: sql.NullTime{
			Time:  pubDate,
			Valid: true,
		},
		Source: sql.NullString{
			String: randomString(10),
			Valid:  true,
		},
		Url: sql.NullString{
			String: "https://" + randomString(8) + ".com/news/" + randomString(10),
			Valid:  true,
		},
		Summary: sql.NullString{
			String: randomString(50),
			Valid:  true,
		},
	}

	news2, err := testQueries.UpdateContactNews(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, news2)

	require.Equal(t, news1.MentionID, news2.MentionID)
	require.Equal(t, news1.ContactID, news2.ContactID)
	require.Equal(t, arg.Title, news2.Title)
	require.Equal(t, arg.PublicationDate.Time.Day(), news2.PublicationDate.Time.Day())
	require.Equal(t, arg.PublicationDate.Time.Month(), news2.PublicationDate.Time.Month())
	require.Equal(t, arg.PublicationDate.Time.Year(), news2.PublicationDate.Time.Year())
	require.Equal(t, arg.Source, news2.Source)
	require.Equal(t, arg.Url, news2.Url)
	require.Equal(t, arg.Summary, news2.Summary)
	require.Equal(t, news1.DatasourceID, news2.DatasourceID)
}

// TestDeleteContactNews tests the DeleteContactNews function
func TestDeleteContactNews(t *testing.T) {
	news1 := createRandomContactNews(t)
	err := testQueries.DeleteContactNews(context.Background(), news1.MentionID)
	require.NoError(t, err)

	news2, err := testQueries.GetContactNewsByID(context.Background(), news1.MentionID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, news2)
}
