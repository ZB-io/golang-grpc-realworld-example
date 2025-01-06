package store

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserStoreIsFollowing(t *testing.T) {

	// scenarios to test
	tests := []struct {
		name           string
		a              *model.User
		b              *model.User
		mockDbBehavior func(db sqlmock.Sqlmock, a *model.User, b *model.User)
		expectedResult bool
		expectError    bool
	}{
		{
			name: "Test when both users are not nil and there is a following relationship between them",
			a:    &model.User{Model: gorm.Model{ID: 1}},
			b:    &model.User{Model: gorm.Model{ID: 2}},
			mockDbBehavior: func(db sqlmock.Sqlmock, a *model.User, b *model.User) {
				db.ExpectQuery("^SELECT count\\(\\*\\) FROM \"follows\"  WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)$").
					WithArgs(a.Model.ID, b.Model.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedResult: true,
			expectError:    false,
		},
		{
			name: "Test when both users are not nil and there is no following relationship between them",
			a:    &model.User{Model: gorm.Model{ID: 1}},
			b:    &model.User{Model: gorm.Model{ID: 2}},
			mockDbBehavior: func(db sqlmock.Sqlmock, a *model.User, b *model.User) {
				db.ExpectQuery("^SELECT count\\(\\*\\) FROM \"follows\"  WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)$").
					WithArgs(a.Model.ID, b.Model.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name: "Test when one or both users are nil",
			a:    &model.User{Model: gorm.Model{ID: 1}},
			b:    nil,
			mockDbBehavior: func(db sqlmock.Sqlmock, a *model.User, b *model.User) {
			},
			expectedResult: false,
			expectError:    false,
		},
		{
			name: "Test when there is a database error",
			a:    &model.User{Model: gorm.Model{ID: 1}},
			b:    &model.User{Model: gorm.Model{ID: 2}},
			mockDbBehavior: func(db sqlmock.Sqlmock, a *model.User, b *model.User) {
				db.ExpectQuery("^SELECT count\\(\\*\\) FROM \"follows\"  WHERE \\(from_user_id = \\$1 AND to_user_id = \\$2\\)$").
					WithArgs(a.Model.ID, b.Model.ID).
					WillReturnError(errors.New("some db error"))
			},
			expectedResult: false,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormDb, _ := gorm.Open("postgres", db)
			tt.mockDbBehavior(mock, tt.a, tt.b)

			userStore := &UserStore{db: gormDb}
			result, err := userStore.IsFollowing(tt.a, tt.b)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
