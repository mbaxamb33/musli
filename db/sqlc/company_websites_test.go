package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomCompanyWebsite creates a company website with random values for testing
func createRandomCompanyWebsite(t *testing.T) CompanyWebsite {
	company := createRandomCompany(t)
	datasource := createRandomDatasource(t)

	arg := CreateCompanyWebsiteParams{
		CompanyID: company.CompanyID,
		BaseUrl:   "https://www." + randomString(8) + ".com",
		SiteTitle: sql.NullString{
			String: company.Name + " Website",
			Valid:  true,
		},
		ScrapeFrequencyDays: sql.NullInt32{
			Int32: int32(randomInt(1, 30)),
			Valid: true,
		},
		IsActive: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	website, err := testQueries.CreateCompanyWebsite(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, website)

	require.Equal(t, arg.CompanyID, website.CompanyID)
	require.Equal(t, arg.BaseUrl, website.BaseUrl)
	require.Equal(t, arg.SiteTitle, website.SiteTitle)
	require.Equal(t, arg.ScrapeFrequencyDays, website.ScrapeFrequencyDays)
	require.Equal(t, arg.IsActive, website.IsActive)
	require.Equal(t, arg.DatasourceID, website.DatasourceID)

	require.NotZero(t, website.WebsiteID)

	return website
}

// TestCreateCompanyWebsite tests the CreateCompanyWebsite function
func TestCreateCompanyWebsite(t *testing.T) {
	createRandomCompanyWebsite(t)
}

// TestGetCompanyWebsiteByID tests the GetCompanyWebsiteByID function
func TestGetCompanyWebsiteByID(t *testing.T) {
	website1 := createRandomCompanyWebsite(t)
	website2, err := testQueries.GetCompanyWebsiteByID(context.Background(), website1.WebsiteID)
	require.NoError(t, err)
	require.NotEmpty(t, website2)

	require.Equal(t, website1.WebsiteID, website2.WebsiteID)
	require.Equal(t, website1.CompanyID, website2.CompanyID)
	require.Equal(t, website1.BaseUrl, website2.BaseUrl)
	require.Equal(t, website1.SiteTitle, website2.SiteTitle)
	require.Equal(t, website1.ScrapeFrequencyDays, website2.ScrapeFrequencyDays)
	require.Equal(t, website1.IsActive, website2.IsActive)
	require.Equal(t, website1.DatasourceID, website2.DatasourceID)

	// LastScrapedAt is initially NULL, so we don't need to check it
}

// TestGetCompanyWebsitesByCompanyID tests the GetCompanyWebsitesByCompanyID function
func TestGetCompanyWebsitesByCompanyID(t *testing.T) {
	company := createRandomCompany(t)

	// Create several websites for the same company
	for i := 0; i < 5; i++ {
		arg := CreateCompanyWebsiteParams{
			CompanyID: company.CompanyID,
			BaseUrl:   "https://www." + randomString(8) + ".com",
			SiteTitle: sql.NullString{
				String: company.Name + " Website " + randomString(5),
				Valid:  true,
			},
			ScrapeFrequencyDays: sql.NullInt32{
				Int32: int32(randomInt(1, 30)),
				Valid: true,
			},
			IsActive: sql.NullBool{
				Bool:  randomBool(),
				Valid: true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateCompanyWebsite(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create websites for other companies
	for i := 0; i < 3; i++ {
		createRandomCompanyWebsite(t)
	}

	arg := GetCompanyWebsitesByCompanyIDParams{
		CompanyID: company.CompanyID,
		Limit:     10,
		Offset:    0,
	}

	websites, err := testQueries.GetCompanyWebsitesByCompanyID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, websites)
	require.Len(t, websites, 5)

	// Verify all websites belong to the same company
	for _, website := range websites {
		require.Equal(t, company.CompanyID, website.CompanyID)
	}
}

// TestUpdateCompanyWebsite tests the UpdateCompanyWebsite function
func TestUpdateCompanyWebsite(t *testing.T) {
	website1 := createRandomCompanyWebsite(t)
	datasource := createRandomDatasource(t)

	// Use UTC time to avoid timezone differences when comparing with database times
	lastScraped := time.Now().UTC().Add(-48 * time.Hour) // 2 days ago, in UTC

	arg := UpdateCompanyWebsiteParams{
		WebsiteID: website1.WebsiteID,
		BaseUrl:   website1.BaseUrl,
		SiteTitle: sql.NullString{
			String: "Updated Website Title " + randomString(5),
			Valid:  true,
		},
		LastScrapedAt: sql.NullTime{
			Time:  lastScraped,
			Valid: true,
		},
		ScrapeFrequencyDays: sql.NullInt32{
			Int32: 7, // Weekly
			Valid: true,
		},
		IsActive: sql.NullBool{
			Bool:  !website1.IsActive.Bool, // Toggle active status
			Valid: true,
		},
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	website2, err := testQueries.UpdateCompanyWebsite(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, website2)

	require.Equal(t, website1.WebsiteID, website2.WebsiteID)
	require.Equal(t, website1.CompanyID, website2.CompanyID)
	require.Equal(t, arg.BaseUrl, website2.BaseUrl)
	require.Equal(t, arg.SiteTitle, website2.SiteTitle)
	require.Equal(t, arg.ScrapeFrequencyDays, website2.ScrapeFrequencyDays)
	require.Equal(t, arg.IsActive, website2.IsActive)
	require.Equal(t, arg.DatasourceID, website2.DatasourceID)

	// Compare both times in UTC to avoid timezone issues
	require.WithinDuration(t, lastScraped.UTC(), website2.LastScrapedAt.Time.UTC(), time.Second)
}

// TestUpdateLastScrapedAt tests the UpdateLastScrapedAt function
func TestUpdateLastScrapedAt(t *testing.T) {
	website1 := createRandomCompanyWebsite(t)

	// Update the last scraped timestamp
	website2, err := testQueries.UpdateLastScrapedAt(context.Background(), website1.WebsiteID)
	require.NoError(t, err)
	require.NotEmpty(t, website2)

	require.Equal(t, website1.WebsiteID, website2.WebsiteID)
	require.Equal(t, website1.CompanyID, website2.CompanyID)
	require.Equal(t, website1.BaseUrl, website2.BaseUrl)
	require.Equal(t, website1.SiteTitle, website2.SiteTitle)
	require.Equal(t, website1.ScrapeFrequencyDays, website2.ScrapeFrequencyDays)
	require.Equal(t, website1.IsActive, website2.IsActive)
	require.Equal(t, website1.DatasourceID, website2.DatasourceID)

	// LastScrapedAt should now be set
	require.True(t, website2.LastScrapedAt.Valid)

	// Use a larger duration tolerance to account for possible time zone differences
	require.WithinDuration(t, time.Now().UTC(), website2.LastScrapedAt.Time.UTC(), 5*time.Second)
}

// TestListCompanyWebsitesForScraping tests the ListCompanyWebsitesForScraping function
func TestListCompanyWebsitesForScraping(t *testing.T) {
	// Create websites that need scraping (active with null LastScrapedAt)
	for i := 0; i < 3; i++ {
		website := createRandomCompanyWebsite(t)

		// Ensure it's active
		updateArg := UpdateCompanyWebsiteParams{
			WebsiteID: website.WebsiteID,
			BaseUrl:   website.BaseUrl,
			SiteTitle: website.SiteTitle,
			LastScrapedAt: sql.NullTime{
				Valid: false,
			},
			ScrapeFrequencyDays: website.ScrapeFrequencyDays,
			IsActive: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			DatasourceID: website.DatasourceID,
		}

		_, err := testQueries.UpdateCompanyWebsite(context.Background(), updateArg)
		require.NoError(t, err)
	}

	// Create websites that were recently scraped (shouldn't need scraping)
	for i := 0; i < 2; i++ {
		website := createRandomCompanyWebsite(t)

		// Set LastScrapedAt to recent time
		updateArg := UpdateCompanyWebsiteParams{
			WebsiteID: website.WebsiteID,
			BaseUrl:   website.BaseUrl,
			SiteTitle: website.SiteTitle,
			LastScrapedAt: sql.NullTime{
				Time:  time.Now().UTC(), // Use UTC time
				Valid: true,
			},
			ScrapeFrequencyDays: sql.NullInt32{
				Int32: 7, // Weekly
				Valid: true,
			},
			IsActive: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			DatasourceID: website.DatasourceID,
		}

		_, err := testQueries.UpdateCompanyWebsite(context.Background(), updateArg)
		require.NoError(t, err)
	}

	// Create websites that need scraping (active with old LastScrapedAt)
	for i := 0; i < 2; i++ {
		website := createRandomCompanyWebsite(t)

		// Set LastScrapedAt to old time exceeding frequency
		updateArg := UpdateCompanyWebsiteParams{
			WebsiteID: website.WebsiteID,
			BaseUrl:   website.BaseUrl,
			SiteTitle: website.SiteTitle,
			LastScrapedAt: sql.NullTime{
				Time:  time.Now().UTC().Add(-8 * 24 * time.Hour), // 8 days ago, in UTC
				Valid: true,
			},
			ScrapeFrequencyDays: sql.NullInt32{
				Int32: 7, // Weekly
				Valid: true,
			},
			IsActive: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			DatasourceID: website.DatasourceID,
		}

		_, err := testQueries.UpdateCompanyWebsite(context.Background(), updateArg)
		require.NoError(t, err)
	}

	// Create inactive websites (shouldn't need scraping regardless of timestamp)
	for i := 0; i < 2; i++ {
		website := createRandomCompanyWebsite(t)

		updateArg := UpdateCompanyWebsiteParams{
			WebsiteID: website.WebsiteID,
			BaseUrl:   website.BaseUrl,
			SiteTitle: website.SiteTitle,
			LastScrapedAt: sql.NullTime{
				Time:  time.Now().UTC().Add(-30 * 24 * time.Hour), // 30 days ago, in UTC
				Valid: true,
			},
			ScrapeFrequencyDays: sql.NullInt32{
				Int32: 7, // Weekly
				Valid: true,
			},
			IsActive: sql.NullBool{
				Bool:  false, // Inactive
				Valid: true,
			},
			DatasourceID: website.DatasourceID,
		}

		_, err := testQueries.UpdateCompanyWebsite(context.Background(), updateArg)
		require.NoError(t, err)
	}

	// Get websites for scraping
	websites, err := testQueries.ListCompanyWebsitesForScraping(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, websites)

	// Should have at least 5 websites (3 with null LastScrapedAt + 2 with old timestamps)
	require.GreaterOrEqual(t, len(websites), 5)

	// All returned websites should be active
	for _, website := range websites {
		require.True(t, website.IsActive.Valid && website.IsActive.Bool, "Only active websites should be returned")

		// If LastScrapedAt is set, it should be older than the frequency
		if website.LastScrapedAt.Valid {
			// Convert both times to UTC to avoid timezone issues
			scrapedTime := website.LastScrapedAt.Time.UTC()
			dueForScraping := time.Now().UTC().After(scrapedTime.Add(
				time.Duration(website.ScrapeFrequencyDays.Int32) * 24 * time.Hour))
			require.True(t, dueForScraping, "Website should be due for scraping")
		}
	}
}

// TestDeleteCompanyWebsite tests the DeleteCompanyWebsite function
func TestDeleteCompanyWebsite(t *testing.T) {
	website1 := createRandomCompanyWebsite(t)
	err := testQueries.DeleteCompanyWebsite(context.Background(), website1.WebsiteID)
	require.NoError(t, err)

	website2, err := testQueries.GetCompanyWebsiteByID(context.Background(), website1.WebsiteID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, website2)
}
