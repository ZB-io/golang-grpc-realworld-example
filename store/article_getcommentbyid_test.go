package store

import (
    "errors"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "github.com/raahii/golang-grpc-realworld-example/model"
)

/* Disabled redeclared type and method
type ArticleStore struct {
 db *gorm.DB
}

func (s *ArticleStore) GetCommentByID(id uint) (*model.Comment, error) {
 var m model.Comment
 err := s.db.Find(&m, id).Error
 if err != nil {
  return nil, err
 }
 return &m, nil
}
*/

func TestArticleStoreGetCommentByID(t *testing.T) {
    knownID := uint(404)
    comment := &model.Comment{ Body: "Test Comment", UserID: knownID }

    tests := []struct {
        name      string
        setupMock func(mock sqlmock.Sqlmock)
        ID        uint
        wantError bool
    }{
        {
            "Successful retrieval of comment by ID",
            func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"ID", "Body", "UserID"}).
                    AddRow(comment.ID, comment.Body, comment.UserID)
                mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").WillReturnRows(rows)
            },
            knownID,
            false,
        },
        {
            "Retrieval of non-existent comment by ID",
            func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"ID", "Body", "UserID"})
                mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").WillReturnRows(rows)
            },
            uint(999),
            true,
        },
        {
            "Database retrieval error",
            func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery("^SELECT (.+) FROM `comments` WHERE (.+)").WillReturnError(errors.New("database error"))
            },
            uint(500),
            true,
        },
    }

    for _, tt := range tests {
        tt := tt 
        t.Run(tt.name, func(t *testing.T) { 
            db, mock, _ := sqlmock.New()
            gDb, _ := gorm.Open("sqlite3", db)
            defer db.Close()

            tt.setupMock(mock)

            s := &ArticleStore{db: gDb}

            _, err := s.GetCommentByID(tt.ID)

            if tt.wantError {
                if err == nil {
                    t.Errorf("Error: Expected an error but got none")
                } 
            } else {
                if err != nil {
                    t.Errorf("Error: Did not expect an error but got one, %v", err.Error())
                }
            }
        })
    }
}
