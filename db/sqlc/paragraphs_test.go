package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// createRandomParagraph creates a paragraph with random values for testing
func createRandomParagraph(t *testing.T) Paragraph {
	datasource := createRandomDatasource(t)

	arg := CreateParagraphParams{
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
		Content: randomString(100),
		MainIdea: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
		Classification: sql.NullString{
			String: []string{"introduction", "background", "analysis", "conclusion", "testimonial"}[randomInt(0, 4)],
			Valid:  true,
		},
		ConfidenceScore: sql.NullString{
			String: formatScore(float64(randomInt(50, 99)) / 100.0),
			Valid:  true,
		},
	}

	paragraph, err := testQueries.CreateParagraph(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paragraph)

	require.Equal(t, arg.DatasourceID, paragraph.DatasourceID)
	require.Equal(t, arg.Content, paragraph.Content)
	require.Equal(t, arg.MainIdea, paragraph.MainIdea)
	require.Equal(t, arg.Classification, paragraph.Classification)
	require.Equal(t, arg.ConfidenceScore, paragraph.ConfidenceScore)

	require.NotZero(t, paragraph.ParagraphID)

	return paragraph
}

// TestCreateParagraph tests the CreateParagraph function
func TestCreateParagraph(t *testing.T) {
	createRandomParagraph(t)
}

// TestGetParagraphByID tests the GetParagraphByID function
func TestGetParagraphByID(t *testing.T) {
	paragraph1 := createRandomParagraph(t)
	paragraph2, err := testQueries.GetParagraphByID(context.Background(), paragraph1.ParagraphID)
	require.NoError(t, err)
	require.NotEmpty(t, paragraph2)

	require.Equal(t, paragraph1.ParagraphID, paragraph2.ParagraphID)
	require.Equal(t, paragraph1.DatasourceID, paragraph2.DatasourceID)
	require.Equal(t, paragraph1.Content, paragraph2.Content)
	require.Equal(t, paragraph1.MainIdea, paragraph2.MainIdea)
	require.Equal(t, paragraph1.Classification, paragraph2.Classification)
	require.Equal(t, paragraph1.ConfidenceScore, paragraph2.ConfidenceScore)
}

// TestListParagraphsByDatasource tests the ListParagraphsByDatasource function
func TestListParagraphsByDatasource(t *testing.T) {
	datasource := createRandomDatasource(t)

	// Create several paragraphs from the same datasource
	for i := 0; i < 5; i++ {
		arg := CreateParagraphParams{
			DatasourceID: sql.NullInt32{
				Int32: datasource.DatasourceID,
				Valid: true,
			},
			Content: randomString(100),
			MainIdea: sql.NullString{
				String: randomString(30),
				Valid:  true,
			},
			Classification: sql.NullString{
				String: []string{"introduction", "background", "analysis", "conclusion", "testimonial"}[randomInt(0, 4)],
				Valid:  true,
			},
			ConfidenceScore: sql.NullString{
				String: formatScore(float64(randomInt(50, 99)) / 100.0),
				Valid:  true,
			},
		}

		_, err := testQueries.CreateParagraph(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListParagraphsByDatasourceParams{
		DatasourceID: sql.NullInt32{
			Int32: datasource.DatasourceID,
			Valid: true,
		},
		Limit:  10,
		Offset: 0,
	}

	paragraphs, err := testQueries.ListParagraphsByDatasource(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paragraphs)
	require.Len(t, paragraphs, 5)

	// Verify all paragraphs come from the same datasource
	for _, paragraph := range paragraphs {
		require.Equal(t, datasource.DatasourceID, paragraph.DatasourceID.Int32)
	}
}

// TestListParagraphsByClassification tests the ListParagraphsByClassification function
func TestListParagraphsByClassification(t *testing.T) {
	classification := "analysis"

	// Create several paragraphs with the same classification
	for i := 0; i < 5; i++ {
		arg := CreateParagraphParams{
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
			Content: randomString(100),
			MainIdea: sql.NullString{
				String: randomString(30),
				Valid:  true,
			},
			Classification: sql.NullString{
				String: classification,
				Valid:  true,
			},
			ConfidenceScore: sql.NullString{
				String: formatScore(0.5 + float64(i)*0.1), // Scores from 0.5 to 0.9
				Valid:  true,
			},
		}

		_, err := testQueries.CreateParagraph(context.Background(), arg)
		require.NoError(t, err)
	}

	// Create some paragraphs with different classifications
	for i := 0; i < 3; i++ {
		otherClassification := []string{"introduction", "background", "conclusion", "testimonial"}[randomInt(0, 3)]

		arg := CreateParagraphParams{
			DatasourceID: sql.NullInt32{
				Int32: createRandomDatasource(t).DatasourceID,
				Valid: true,
			},
			Content: randomString(100),
			MainIdea: sql.NullString{
				String: randomString(30),
				Valid:  true,
			},
			Classification: sql.NullString{
				String: otherClassification,
				Valid:  true,
			},
			ConfidenceScore: sql.NullString{
				String: formatScore(float64(randomInt(50, 99)) / 100.0),
				Valid:  true,
			},
		}

		_, err := testQueries.CreateParagraph(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListParagraphsByClassificationParams{
		Classification: sql.NullString{
			String: classification,
			Valid:  true,
		},
		Limit:  10,
		Offset: 0,
	}

	paragraphs, err := testQueries.ListParagraphsByClassification(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paragraphs)
	require.GreaterOrEqual(t, len(paragraphs), 5)

	// Verify all paragraphs have the same classification
	for _, paragraph := range paragraphs {
		require.Equal(t, classification, paragraph.Classification.String)
	}

	// Verify they're sorted by confidence score (descending)
	for i := 0; i < len(paragraphs)-1; i++ {
		score1 := parseScore(paragraphs[i].ConfidenceScore.String)
		score2 := parseScore(paragraphs[i+1].ConfidenceScore.String)
		require.GreaterOrEqual(t, score1, score2)
	}
}

// TestUpdateParagraph tests the UpdateParagraph function
func TestUpdateParagraph(t *testing.T) {
	paragraph1 := createRandomParagraph(t)

	arg := UpdateParagraphParams{
		ParagraphID: paragraph1.ParagraphID,
		Content:     randomString(100),
		MainIdea: sql.NullString{
			String: randomString(30),
			Valid:  true,
		},
		Classification: sql.NullString{
			String: []string{"introduction", "background", "analysis", "conclusion", "testimonial"}[randomInt(0, 4)],
			Valid:  true,
		},
		ConfidenceScore: sql.NullString{
			String: formatScore(float64(randomInt(50, 99)) / 100.0),
			Valid:  true,
		},
	}

	paragraph2, err := testQueries.UpdateParagraph(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, paragraph2)

	require.Equal(t, paragraph1.ParagraphID, paragraph2.ParagraphID)
	require.Equal(t, paragraph1.DatasourceID, paragraph2.DatasourceID)
	require.Equal(t, arg.Content, paragraph2.Content)
	require.Equal(t, arg.MainIdea, paragraph2.MainIdea)
	require.Equal(t, arg.Classification, paragraph2.Classification)
	require.Equal(t, arg.ConfidenceScore, paragraph2.ConfidenceScore)
}

// TestDeleteParagraph tests the DeleteParagraph function
func TestDeleteParagraph(t *testing.T) {
	paragraph1 := createRandomParagraph(t)
	err := testQueries.DeleteParagraph(context.Background(), paragraph1.ParagraphID)
	require.NoError(t, err)

	paragraph2, err := testQueries.GetParagraphByID(context.Background(), paragraph1.ParagraphID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, paragraph2)
}
