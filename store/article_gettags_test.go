package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"log"
	"testing"
)

// Instead of redeclaring ArticleStore and GetTags method in the test file, we should use them from the original package where they were declared.
// 
// type ArticleStore struct {
//     db *gorm.DB
// }
// func (s *ArticleStore) GetTags() ([]model.Tag, error) {
//     var tags []model.Tag
//     if err := s.db.Find(&tags).Error; err != nil {
// 	    return tags, err
//     }
//     return tags, nil
// }

type TestSet struct {
	Name          string
	DBError       error
	DBTags        []model.Tag
	ExpectedError error
	ExpectedTags  []model.Tag
}

func TestArticleStoreGetTags(t *testing.T) {
	testSets := []TestSet{
		{
			Name: "Successful retrieval of all tags",
			DBTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
			},
			ExpectedTags: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
			},
		},
		{
			Name:          "Failure due to an error from the database transaction",
			DBError:       errors.New("database error"),
			ExpectedError: errors.New("database error"),
		},
		{
			Name:         "Successful retrieval of an empty tag list",
			ExpectedTags: []model.Tag{},
		},
	}

	for _, test := range testSets {
		t.Run(test.Name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				log.Fatalf("Failed to open the database: %v", err)
			}

			// Here, use the GetTags method from the original ArticleStore struct, not from the redeclared one in the test file.
			// s := &ArticleStore{gdb}
            s := &store.ArticleStore{gdb}

			resultTags, resultErr := s.GetTags()

			if resultErr != test.DBError {
				t.Errorf("Received unexpected error %v", resultErr)
				t.FailNow()
			}

			if len(test.ExpectedTags) != len(resultTags) {
				t.Errorf("Received incorrect number of tags. Expecting: %v, received: %v", len(test.ExpectedTags), len(resultTags))
			}

			for i := range test.ExpectedTags {
				if test.ExpectedTags[i].Name != resultTags[i].Name {
					t.Errorf("Failed on %v. Expecting: %v, received: %v", test.Name, test.ExpectedTags[i].Name, resultTags[i].Name)
					t.FailNow()
				}
			}
		})
	}
}
