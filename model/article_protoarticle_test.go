package model

import (
	"testing"
	"time"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/proto"
	"github.com/stretchr/testify/assert"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type User struct {
	gorm.Model
	Username         string    `gorm:"unique_index;not null"`
	Email            string    `gorm:"unique_index;not null"`
	Password         string    `gorm:"not null"`
	Bio              string    `gorm:"not null"`
	Image            string    `gorm:"not null"`
	Follows          []User    `gorm:"many2many:follows;jointable_foreignkey:from_user_id;association_jointable_foreignkey:to_user_id"`
	FavoriteArticles []Article `gorm:"many2many:favorite_articles;"`
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestProtoArticle(t *testing.T) {
	type testCase struct {
		name          string
		article       model.Article
		favorited     bool
		expectedProto proto.Article
	}

	createdTime, _ := time.Parse(model.ISO8601, "2023-10-06T12:00:00+0000Z")
	updatedTime, _ := time.Parse(model.ISO8601, "2023-10-07T12:00:00+0000Z")

	var mockUser model.User = model.User{
		Username: "testuser",
	}

	tests := []testCase{
		{
			name: "Convert Article with All Fields Populated",
			article: model.Article{
				Model:          gorm.Model{ID: 123, CreatedAt: createdTime, UpdatedAt: updatedTime},
				Title:          "Test Title",
				Description:    "This is a test description",
				Body:           "Test body content",
				Tags:           []model.Tag{{Name: "test1"}, {Name: "test2"}},
				Author:         mockUser,
				UserID:         1,
				FavoritesCount: 10,
			},
			favorited: true,
			expectedProto: proto.Article{
				Slug:           "123",
				Title:          "Test Title",
				Description:    "This is a test description",
				Body:           "Test body content",
				TagList:        []string{"test1", "test2"},
				CreatedAt:      "2023-10-06T12:00:00+0000Z",
				UpdatedAt:      "2023-10-07T12:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 10,
				Author:         nil,
			},
		},
		{
			name: "Convert Article with No Tags",
			article: model.Article{
				Model:          gorm.Model{ID: 456, CreatedAt: createdTime, UpdatedAt: updatedTime},
				Title:          "Another Test Title",
				Description:    "Another description",
				Body:           "Another body content",
				Tags:           []model.Tag{},
				Author:         mockUser,
				UserID:         1,
				FavoritesCount: 25,
			},
			favorited: false,
			expectedProto: proto.Article{
				Slug:           "456",
				Title:          "Another Test Title",
				Description:    "Another description",
				Body:           "Another body content",
				TagList:        []string{},
				CreatedAt:      "2023-10-06T12:00:00+0000Z",
				UpdatedAt:      "2023-10-07T12:00:00+0000Z",
				Favorited:      false,
				FavoritesCount: 25,
				Author:         nil,
			},
		},
		{
			name: "Convert Article Edge Case with Default Values",
			article: model.Article{
				Model:          gorm.Model{ID: 789, CreatedAt: createdTime, UpdatedAt: updatedTime},
				Title:          "",
				Description:    "",
				Body:           "",
				Tags:           []model.Tag{},
				Author:         mockUser,
				UserID:         1,
				FavoritesCount: 0,
			},
			favorited: true,
			expectedProto: proto.Article{
				Slug:           "789",
				Title:          "",
				Description:    "",
				Body:           "",
				TagList:        []string{},
				CreatedAt:      "2023-10-06T12:00:00+0000Z",
				UpdatedAt:      "2023-10-07T12:00:00+0000Z",
				Favorited:      true,
				FavoritesCount: 0,
				Author:         nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualProto := tc.article.ProtoArticle(tc.favorited)

			assert.Equal(t, tc.expectedProto.Slug, actualProto.Slug)
			assert.Equal(t, tc.expectedProto.Title, actualProto.Title)
			assert.Equal(t, tc.expectedProto.Description, actualProto.Description)
			assert.Equal(t, tc.expectedProto.Body, actualProto.Body)
			assert.ElementsMatch(t, tc.expectedProto.TagList, actualProto.TagList)
			assert.Equal(t, tc.expectedProto.CreatedAt, actualProto.CreatedAt)
			assert.Equal(t, tc.expectedProto.UpdatedAt, actualProto.UpdatedAt)
			assert.Equal(t, tc.expectedProto.Favorited, actualProto.Favorited)
			assert.Equal(t, tc.expectedProto.FavoritesCount, actualProto.FavoritesCount)

		})
	}
}
