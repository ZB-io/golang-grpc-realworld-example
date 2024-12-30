package store

import (
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestUpdate(t *testing.T) {
	tests := []struct {
		name        string
		setupDB     func(*gorm.DB)
		input       *model.Article
		expectedErr bool
		validate    func(*testing.T, *gorm.DB, *model.Article)
	}{
		{
			name: "Successfully Update an Existing Article",
			setupDB: func(db *gorm.DB) {
				article := &model.Article{
					Model:       gorm.Model{ID: 1},
					Title:       "Original Title",
					Description: "Original Description",
					Body:        "Original Body",
					UserID:      1,
				}
				require.NoError(t, db.Create(article).Error)
			},
			input: &model.Article{
				Model:       gorm.Model{ID: 1},
				Title:       "Updated Title",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			expectedErr: false,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var updatedArticle model.Article
				err := db.First(&updatedArticle, input.ID).Error
				require.NoError(t, err)
				assert.Equal(t, input.Title, updatedArticle.Title)
				assert.Equal(t, input.Description, updatedArticle.Description)
				assert.Equal(t, input.Body, updatedArticle.Body)
			},
		},
		{
			name:    "Attempt to Update a Non-existent Article",
			setupDB: func(db *gorm.DB) {},
			input: &model.Article{
				Model: gorm.Model{ID: 999},
				Title: "Non-existent Article",
			},
			expectedErr: true,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var article model.Article
				err := db.First(&article, input.ID).Error
				assert.Error(t, err)
				assert.True(t, gorm.IsRecordNotFoundError(err))
			},
		},
		{
			name: "Update Article with Invalid Data",
			setupDB: func(db *gorm.DB) {
				article := &model.Article{
					Model:       gorm.Model{ID: 2},
					Title:       "Original Title",
					Description: "Original Description",
					Body:        "Original Body",
					UserID:      1,
				}
				require.NoError(t, db.Create(article).Error)
			},
			input: &model.Article{
				Model:       gorm.Model{ID: 2},
				Title:       "",
				Description: "Updated Description",
				Body:        "Updated Body",
				UserID:      1,
			},
			expectedErr: true,
			validate: func(t *testing.T, db *gorm.DB, input *model.Article) {
				var originalArticle model.Article
				err := db.First(&originalArticle, input.ID).Error
				require.NoError(t, err)
				assert.NotEqual(t, input.Title, originalArticle.Title)
				assert.NotEqual(t, input.Description, originalArticle.Description)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open("sqlite3", ":memory:")
			require.NoError(t, err)
			defer db.Close()

			require.NoError(t, db.AutoMigrate(&model.Article{}, &model.Tag{}).Error)

			tt.setupDB(db)

			store := &ArticleStore{db: db}
			err = store.Update(tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, db, tt.input)
		})
	}
}
