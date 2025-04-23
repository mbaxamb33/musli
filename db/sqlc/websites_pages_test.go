package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRandomWebsitePage creates a website page with random values for testing
func createRandomWebsitePage(t *testing.T) WebsitePage {
	// First create a company website for the page to belong to
	companyWebsite := createRandomCompanyWebsite(t)

	arg := CreateWebsitePageParams{
		WebsiteID: companyWebsite.WebsiteID,
		Url:       "https://" + randomString(8) + ".com/" + randomString(6),
		Path:      "/" + randomString(6),
		Title: sql.NullString{
			String: "Page Title " + randomString(10),
			Valid:  true,
		},
		ParentPageID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		Depth: 1,
		ExtractStatus: sql.NullString{
			String: "pending",
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: createRandomDatasource(t).DatasourceID,
			Valid: true,
		},
	}

	page, err := testQueries.CreateWebsitePage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, page)

	require.Equal(t, arg.WebsiteID, page.WebsiteID)
	require.Equal(t, arg.Url, page.Url)
	require.Equal(t, arg.Path, page.Path)
	require.Equal(t, arg.Title, page.Title)
	require.Equal(t, arg.ParentPageID, page.ParentPageID)
	require.Equal(t, arg.Depth, page.Depth)
	require.Equal(t, arg.ExtractStatus, page.ExtractStatus)
	require.Equal(t, arg.DatasourceID, page.DatasourceID)

	require.NotZero(t, page.PageID)

	return page
}

// TestCreateWebsitePage tests the CreateWebsitePage function
func TestCreateWebsitePage(t *testing.T) {
	createRandomWebsitePage(t)
}

// TestGetWebsitePageByID tests the GetWebsitePageByID function
func TestGetWebsitePageByID(t *testing.T) {
	page1 := createRandomWebsitePage(t)
	page2, err := testQueries.GetWebsitePageByID(context.Background(), page1.PageID)
	require.NoError(t, err)
	require.NotEmpty(t, page2)

	require.Equal(t, page1.PageID, page2.PageID)
	require.Equal(t, page1.WebsiteID, page2.WebsiteID)
	require.Equal(t, page1.Url, page2.Url)
	require.Equal(t, page1.Path, page2.Path)
	require.Equal(t, page1.Title, page2.Title)
	require.Equal(t, page1.ParentPageID, page2.ParentPageID)
	require.Equal(t, page1.Depth, page2.Depth)
	require.Equal(t, page1.ExtractStatus, page2.ExtractStatus)
	require.Equal(t, page1.DatasourceID, page2.DatasourceID)
	if page1.LastExtractedAt.Valid && page2.LastExtractedAt.Valid {
		require.WithinDuration(t, page1.LastExtractedAt.Time, page2.LastExtractedAt.Time, time.Second)
	}
}

// TestGetWebsitePageByURL tests the GetWebsitePageByURL function
func TestGetWebsitePageByURL(t *testing.T) {
	page1 := createRandomWebsitePage(t)

	arg := GetWebsitePageByURLParams{
		WebsiteID: page1.WebsiteID,
		Url:       page1.Url,
	}

	page2, err := testQueries.GetWebsitePageByURL(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, page2)

	require.Equal(t, page1.PageID, page2.PageID)
	require.Equal(t, page1.WebsiteID, page2.WebsiteID)
	require.Equal(t, page1.Url, page2.Url)
	require.Equal(t, page1.Path, page2.Path)
	require.Equal(t, page1.Title, page2.Title)
	require.Equal(t, page1.ParentPageID, page2.ParentPageID)
	require.Equal(t, page1.Depth, page2.Depth)
	require.Equal(t, page1.ExtractStatus, page2.ExtractStatus)
	require.Equal(t, page1.DatasourceID, page2.DatasourceID)
}

// TestListWebsitePagesByWebsiteID tests the ListWebsitePagesByWebsiteID function
func TestListWebsitePagesByWebsiteID(t *testing.T) {
	companyWebsite := createRandomCompanyWebsite(t)

	// Create several pages for the same website
	for i := 0; i < 5; i++ {
		arg := CreateWebsitePageParams{
			WebsiteID: companyWebsite.WebsiteID,
			Url:       "https://" + randomString(8) + ".com/" + randomString(6) + "/" + randomString(4),
			Path:      "/" + randomString(6) + "/" + randomString(4),
			Title: sql.NullString{
				String: "Page Title " + randomString(10),
				Valid:  true,
			},
			ParentPageID: sql.NullInt32{
				Int32: 0,
				Valid: false,
			},
			Depth: int32(i + 1),
			ExtractStatus: sql.NullString{
				String: "pending",
				Valid:  true,
			},
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
		}

		_, err := testQueries.CreateWebsitePage(context.Background(), arg)
		require.NoError(t, err)
	}

	// Also create some pages for a different website
	for i := 0; i < 3; i++ {
		createRandomWebsitePage(t)
	}

	arg := ListWebsitePagesByWebsiteIDParams{
		WebsiteID: companyWebsite.WebsiteID,
		Limit:     10,
		Offset:    0,
	}

	pages, err := testQueries.ListWebsitePagesByWebsiteID(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, pages)
	require.Len(t, pages, 5)

	// Verify all pages belong to the same website
	for _, page := range pages {
		require.Equal(t, companyWebsite.WebsiteID, page.WebsiteID)
	}
}

// TestUpdateWebsitePage tests the UpdateWebsitePage function
func TestUpdateWebsitePage(t *testing.T) {
	page1 := createRandomWebsitePage(t)
	datasource := createRandomDatasource(t)

	extractTime := time.Now().Add(-24 * time.Hour) // 1 day ago

	arg := UpdateWebsitePageParams{
		WebsiteID: page1.WebsiteID,
		Url:       page1.Url,
		Title: sql.NullString{
			String: "Updated Page Title " + randomString(10),
			Valid:  true,
		},
		ParentPageID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		LastExtractedAt: sql.NullTime{
			Time:  extractTime,
			Valid: true,
		},
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
	}

	page2, err := testQueries.UpdateWebsitePage(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, page2)

	require.Equal(t, page1.PageID, page2.PageID)
	require.Equal(t, page1.WebsiteID, page2.WebsiteID)
	require.Equal(t, page1.Url, page2.Url)
	require.Equal(t, page1.Path, page2.Path)
	require.Equal(t, arg.Title, page2.Title)
	require.Equal(t, arg.ParentPageID, page2.ParentPageID)
	require.Equal(t, page1.Depth, page2.Depth)
	require.Equal(t, arg.ExtractStatus, page2.ExtractStatus)
	require.Equal(t, arg.DatasourceID, page2.DatasourceID)
	require.WithinDuration(t, extractTime, page2.LastExtractedAt.Time, time.Second)
}

// TestUpdateExtractStatus tests the UpdateExtractStatus function
func TestUpdateExtractStatus(t *testing.T) {
	page1 := createRandomWebsitePage(t)

	arg := UpdateExtractStatusParams{
		WebsiteID: page1.WebsiteID,
		PageID:    page1.PageID,
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
	}

	page2, err := testQueries.UpdateExtractStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, page2)

	require.Equal(t, page1.PageID, page2.PageID)
	require.Equal(t, page1.WebsiteID, page2.WebsiteID)
	require.Equal(t, page1.Url, page2.Url)
	require.Equal(t, page1.Path, page2.Path)
	require.Equal(t, page1.Title, page2.Title)
	require.Equal(t, page1.ParentPageID, page2.ParentPageID)
	require.Equal(t, page1.Depth, page2.Depth)
	require.Equal(t, arg.ExtractStatus, page2.ExtractStatus)
	require.NotNil(t, page2.LastExtractedAt)
	require.True(t, page2.LastExtractedAt.Valid)
}

// TestGetPagesForExtraction tests the GetPagesForExtraction function
func TestGetPagesForExtraction(t *testing.T) {
	// Create a page with pending status
	pendingPage := createRandomWebsitePage(t)

	// Update a different page to have completed status with old extraction time
	completedPage := createRandomWebsitePage(t)
	oldTime := time.Now().Add(-31 * 24 * time.Hour) // 31 days ago

	updateArg := UpdateWebsitePageParams{
		WebsiteID:    completedPage.WebsiteID,
		Url:          completedPage.Url,
		Title:        completedPage.Title,
		ParentPageID: completedPage.ParentPageID,
		LastExtractedAt: sql.NullTime{
			Time:  oldTime,
			Valid: true,
		},
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: completedPage.DatasourceID,
	}

	_, err := testQueries.UpdateWebsitePage(context.Background(), updateArg)
	require.NoError(t, err)

	// Create another page with recent extraction time that shouldn't be returned
	recentPage := createRandomWebsitePage(t)
	recentTime := time.Now().Add(-1 * 24 * time.Hour) // 1 day ago

	updateRecentArg := UpdateWebsitePageParams{
		WebsiteID:    recentPage.WebsiteID,
		Url:          recentPage.Url,
		Title:        recentPage.Title,
		ParentPageID: recentPage.ParentPageID,
		LastExtractedAt: sql.NullTime{
			Time:  recentTime,
			Valid: true,
		},
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: recentPage.DatasourceID,
	}

	_, err = testQueries.UpdateWebsitePage(context.Background(), updateRecentArg)
	require.NoError(t, err)

	// Get pages for extraction
	pages, err := testQueries.GetPagesForExtraction(context.Background(), 10)
	require.NoError(t, err)
	require.NotEmpty(t, pages)

	// Should have at least the pending page and the old completed page
	require.GreaterOrEqual(t, len(pages), 2)

	// Check if our specific pages are in the results
	foundPending := false
	foundCompleted := false
	foundRecent := false

	for _, page := range pages {
		if page.PageID == pendingPage.PageID {
			foundPending = true
		}
		if page.PageID == completedPage.PageID {
			foundCompleted = true
		}
		if page.PageID == recentPage.PageID {
			foundRecent = true
		}
	}

	require.True(t, foundPending, "Pending page should be returned")
	require.True(t, foundCompleted, "Old completed page should be returned")
	require.False(t, foundRecent, "Recent completed page should not be returned")
}

// TestDeleteWebsitePage tests the DeleteWebsitePage function
func TestDeleteWebsitePage(t *testing.T) {
	page1 := createRandomWebsitePage(t)
	err := testQueries.DeleteWebsitePage(context.Background(), page1.PageID)
	require.NoError(t, err)

	page2, err := testQueries.GetWebsitePageByID(context.Background(), page1.PageID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, page2)
}

// TestListWebsitePageTree tests the ListWebsitePageTree function
func TestListWebsitePageTree(t *testing.T) {
	companyWebsite := createRandomCompanyWebsite(t)

	// Create a root page
	rootArg := CreateWebsitePageParams{
		WebsiteID: companyWebsite.WebsiteID,
		Url:       "https://" + randomString(8) + ".com/",
		Path:      "/",
		Title: sql.NullString{
			String: "Root Page",
			Valid:  true,
		},
		ParentPageID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
		Depth: 0,
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: createRandomDatasource(t).DatasourceID,
			Valid: true,
		},
	}

	rootPage, err := testQueries.CreateWebsitePage(context.Background(), rootArg)
	require.NoError(t, err)

	// Create a child page
	childArg := CreateWebsitePageParams{
		WebsiteID: companyWebsite.WebsiteID,
		Url:       "https://" + randomString(8) + ".com/child",
		Path:      "/child",
		Title: sql.NullString{
			String: "Child Page",
			Valid:  true,
		},
		ParentPageID: sql.NullInt32{
			Int32: rootPage.PageID,
			Valid: true,
		},
		Depth: 1,
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: createRandomDatasource(t).DatasourceID,
			Valid: true,
		},
	}

	childPage, err := testQueries.CreateWebsitePage(context.Background(), childArg)
	require.NoError(t, err)

	// Create a grandchild page
	grandchildArg := CreateWebsitePageParams{
		WebsiteID: companyWebsite.WebsiteID,
		Url:       "https://" + randomString(8) + ".com/child/grandchild",
		Path:      "/child/grandchild",
		Title: sql.NullString{
			String: "Grandchild Page",
			Valid:  true,
		},
		ParentPageID: sql.NullInt32{
			Int32: childPage.PageID,
			Valid: true,
		},
		Depth: 2,
		ExtractStatus: sql.NullString{
			String: "completed",
			Valid:  true,
		},
		DatasourceID: sql.NullInt32{
			Int32: createRandomDatasource(t).DatasourceID,
			Valid: true,
		},
	}

	_, err = testQueries.CreateWebsitePage(context.Background(), grandchildArg)
	require.NoError(t, err)

	// Get the page tree
	arg := ListWebsitePageTreeParams{
		WebsiteID: companyWebsite.WebsiteID,
		Limit:     10,
		Offset:    0,
	}

	pageTree, err := testQueries.ListWebsitePageTree(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, pageTree)
	require.Len(t, pageTree, 3) // Root, child, and grandchild

	// The first page should be the root
	require.Equal(t, rootPage.PageID, pageTree[0].PageID)
	require.False(t, pageTree[0].ParentPageID.Valid)

	// Find the child and grandchild pages in the results
	var foundChild, foundGrandchild bool

	for _, page := range pageTree {
		if page.PageID == childPage.PageID {
			foundChild = true
			// Check parent relationship
			require.True(t, page.ParentPageID.Valid)
			require.Equal(t, rootPage.PageID, page.ParentPageID.Int32)
		}

		if page.PageID == childPage.PageID {
			foundGrandchild = true
		}
	}

	require.True(t, foundChild, "Child page should be in the tree")
	require.True(t, foundGrandchild, "Grandchild page should be in the tree")
}
