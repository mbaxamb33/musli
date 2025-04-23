package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// // createRandomDatasource creates a datasource for testing
// func createRandomDatasource(t *testing.T) Datasource {
// 	arg := CreateDatasourceParams{
// 		SourceType: randomString(8),
// 		SourceID: sql.NullInt32{
// 			Int32: int32(randomInt(1, 1000)),
// 			Valid: true,
// 		},
// 	}

// 	datasource, err := testQueries.CreateDatasource(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, datasource)

// 	require.Equal(t, arg.SourceType, datasource.SourceType)
// 	require.Equal(t, arg.SourceID, datasource.SourceID)
// 	require.NotZero(t, datasource.DatasourceID)
// 	require.NotEmpty(t, datasource.ExtractionTimestamp)

// 	return datasource
// }

// createRandomCompanyNews creates company news with random values for testing
func createRandomCompanyNews(t *testing.T) CompanyNews {
	company := createRandomCompany(t)
	datasource := createRandomDatasource(t)

	pubDate := time.Now().AddDate(0, -1, 0) // One month ago

	arg := CreateCompanyNewsParams{
		CompanyID: sql.NullInt32{
			Int32: company.CompanyID,
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
		Sentiment: sql.NullString{
			String: []string{"positive", "negative", "neutral"}[randomInt(0, 2)],
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	news, err := testQueries.CreateCompanyNews(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, news)

	require.Equal(t, arg.CompanyID, news.CompanyID)
	require.Equal(t, arg.Title, news.Title)
	require.Equal(t, arg.PublicationDate.Time.Day(), news.PublicationDate.Time.Day())
	require.Equal(t, arg.PublicationDate.Time.Month(), news.PublicationDate.Time.Month())
	require.Equal(t, arg.PublicationDate.Time.Year(), news.PublicationDate.Time.Year())
	require.Equal(t, arg.Source, news.Source)
	require.Equal(t, arg.Url, news.Url)
	require.Equal(t, arg.Summary, news.Summary)
	require.Equal(t, arg.Sentiment, news.Sentiment)
	require.Equal(t, arg.DatasourceID, news.DatasourceID)

	require.NotZero(t, news.NewsID)

	return news
}

// TestCreateCompanyNews tests the CreateCompanyNews function
func TestCreateCompanyNews(t *testing.T) {
	createRandomCompanyNews(t)
}

// TestGetCompanyNewsByID tests the GetCompanyNewsByID function
func TestGetCompanyNewsByID(t *testing.T) {
	news1 := createRandomCompanyNews(t)
	news2, err := testQueries.GetCompanyNewsByID(context.Background(), news1.NewsID)
	require.NoError(t, err)
	require.NotEmpty(t, news2)

	require.Equal(t, news1.NewsID, news2.NewsID)
	require.Equal(t, news1.CompanyID, news2.CompanyID)
	require.Equal(t, news1.Title, news2.Title)
	require.Equal(t, news1.PublicationDate.Time.Format("2006-01-02"),
		news2.PublicationDate.Time.Format("2006-01-02"))
	require.Equal(t, news1.Source, news2.Source)
	require.Equal(t, news1.Url, news2.Url)
	require.Equal(t, news1.Summary, news2.Summary)
	require.Equal(t, news1.Sentiment, news2.Sentiment)
	require.Equal(t, news1.DatasourceID, news2.DatasourceID)
}

// TestListCompanyNewsByCompany tests the ListCompanyNewsByCompany function
func TestListCompanyNewsByCompany(t *testing.T) {
	company := createRandomCompany(t)

	// Create several news items for the same company
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateCompanyNewsParams{
			CompanyID: sql.NullInt32{
				Int32: company.CompanyID,
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
			Sentiment: sql.NullString{
				String: []string{"positive", "negative", "neutral"}[randomInt(0, 2)],
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateCompanyNews(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some news for other companies
	for i := 0; i < 3; i++ {
		createRandomCompanyNews(t)
	}

	arg := ListCompanyNewsByCompanyParams{
		CompanyID: sql.NullInt32{
			Int32: company.CompanyID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListCompanyNewsByCompany(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.Len(t, newsList, 5)

	// Verify all news items belong to the same company
	for _, news := range newsList {
		require.Equal(t, company.CompanyID, news.CompanyID.Int32)
	}

	// Verify they're sorted by publication date (descending)
	for i := 0; i < len(newsList)-1; i++ {
		require.True(t, newsList[i].PublicationDate.Time.After(newsList[i+1].PublicationDate.Time) ||
			newsList[i].PublicationDate.Time.Equal(newsList[i+1].PublicationDate.Time))
	}
}

// TestListCompanyNewsByDatasource tests the ListCompanyNewsByDatasource function
func TestListCompanyNewsByDatasource(t *testing.T) {
	datasource := createRandomDatasource(t)

	// Create several news items from the same datasource
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateCompanyNewsParams{
			CompanyID: sql.NullInt32{
				Int32: createRandomCompany(t).CompanyID,
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
			Sentiment: sql.NullString{
				String: []string{"positive", "negative", "neutral"}[randomInt(0, 2)],
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: datasource.DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateCompanyNews(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListCompanyNewsByDatasourceParams{
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListCompanyNewsByDatasource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.Len(t, newsList, 5)

	// Verify all news items come from the same datasource
	for _, news := range newsList {
		require.Equal(t, datasource.DatasourceID, news.DatasourceID.Int32)
	}
}

// TestListCompanyNewsBySentiment tests the ListCompanyNewsBySentiment function
func TestListCompanyNewsBySentiment(t *testing.T) {
	sentiment := "positive"

	// Create several news items with the same sentiment
	for i := 0; i < 5; i++ {
		pubDate := time.Now().AddDate(0, 0, -i) // Different days

		arg := CreateCompanyNewsParams{
			CompanyID: sql.NullInt32{
				Int32: createRandomCompany(t).CompanyID,
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
			Sentiment: sql.NullString{
				String: sentiment,
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateCompanyNews(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some news with other sentiments
	for i := 0; i < 3; i++ {
		otherSentiment := []string{"negative", "neutral"}[randomInt(0, 1)]

		arg := CreateCompanyNewsParams{
			CompanyID: sql.NullInt32{
				Int32: createRandomCompany(t).CompanyID,
				Valid: true,
			},
			Title: randomString(15),
			PublicationDate: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Sentiment: sql.NullString{
				String: otherSentiment,
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateCompanyNews(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListCompanyNewsBySentimentParams{
		Sentiment: sql.NullString{
			String: sentiment,
			Valid:  true,
		},
		Limit:  10,
		Offset: 0,
	}

	newsList, err := testQueries.ListCompanyNewsBySentiment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newsList)
	require.GreaterOrEqual(t, len(newsList), 5)

	// Verify all news items have the same sentiment
	for _, news := range newsList {
		require.Equal(t, sentiment, news.Sentiment.String)
	}
}

// TestUpdateCompanyNews tests the UpdateCompanyNews function
func TestUpdateCompanyNews(t *testing.T) {
	news1 := createRandomCompanyNews(t)

	pubDate := time.Now().AddDate(0, -2, 0) // Two months ago

	arg := UpdateCompanyNewsParams{
		NewsID: news1.NewsID,
		Title:  randomString(15),
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
		Sentiment: sql.NullString{
			String: []string{"positive", "negative", "neutral"}[randomInt(0, 2)],
			Valid:  true,
		},
	}

	news2, err := testQueries.UpdateCompanyNews(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, news2)

	require.Equal(t, news1.NewsID, news2.NewsID)
	require.Equal(t, news1.CompanyID, news2.CompanyID)
	require.Equal(t, arg.Title, news2.Title)
	require.Equal(t, arg.PublicationDate.Time.Day(), news2.PublicationDate.Time.Day())
	require.Equal(t, arg.PublicationDate.Time.Month(), news2.PublicationDate.Time.Month())
	require.Equal(t, arg.PublicationDate.Time.Year(), news2.PublicationDate.Time.Year())
	require.Equal(t, arg.Source, news2.Source)
	require.Equal(t, arg.Url, news2.Url)
	require.Equal(t, arg.Summary, news2.Summary)
	require.Equal(t, arg.Sentiment, news2.Sentiment)
	require.Equal(t, news1.DatasourceID, news2.DatasourceID)
}

// TestDeleteCompanyNews tests the DeleteCompanyNews function
func TestDeleteCompanyNews(t *testing.T) {
	news1 := createRandomCompanyNews(t)
	err := testQueries.DeleteCompanyNews(context.Background(), news1.NewsID)
	require.NoError(t, err)

	news2, err := testQueries.GetCompanyNewsByID(context.Background(), news1.NewsID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, news2)
}
