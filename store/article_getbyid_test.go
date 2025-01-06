package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)

// TestArticleStoreGetByID is the unit test for GetByID function in store package
func TestArticleStoreGetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		wantErr bool
	}{
		{
			name:    "Successful retrieval of an article by ID",
			id:      1,
			wantErr: false,
		},
		{
			name:    "Attempt to retrieve an article with a non-existing ID",
			id:      100,
			wantErr: true,
		},
		{
			name:    "Database error during article retrieval",
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			gormDB, _ := gorm.Open("mysql", mockDB)

			// TODO: Mock your database operations based on test scenarios here
			switch tt.name {
			case "Successful retrieval of an article by ID":
				// mock your successful operation
			case "Attempt to retrieve an article with a non-existing ID":
				// mock your operation with non-existing ID
			case "Database error during article retrieval":
				// mock your db error
			}

			store := &ArticleStore{db: gormDB}
			_, err = store.GetByID(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %v", err)
				}
			}

			mock.ExpectationsWereMet()
		})
	}
}
