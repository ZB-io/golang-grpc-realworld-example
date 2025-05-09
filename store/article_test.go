package store

import (
	errors "errors"
	fmt "fmt"
	debug "runtime/debug"
	sync "sync"
	testing "testing"

	"github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
)

/*
ROOST_METHOD_HASH=ArticleStore_CreateComment_b16d4a71d4
ROOST_METHOD_SIG_HASH=ArticleStore_CreateComment_7475736b06

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error // CreateComment creates a comment of the article
*/
func TestArticleStoreCreateComment(t *testing.T) {
	type args struct {
		m *model.Comment
	}

	tests := []struct {
		name        string
		args        args
		setupMockDB func() (*gorm.DB, sqlmock.Sqlmock, error)
		wantErr     bool
	}{
		{
			name: "Successful Comment Creation",
			args: args{
				m: &model.Comment{ArticleID: 1, UserID: 1, Body: "Test comment"},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("INSERT INTO \"comments\"").WillReturnResult(sqlmock.NewResult(1, 1))
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: false,
		},
		{
			name: "Invalid Comment Creation",
			args: args{
				m: &model.Comment{},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
		{
			name: "Database Connection Error",
			args: args{
				m: &model.Comment{ArticleID: 1, UserID: 1, Body: "Test comment"},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("INSERT INTO \"comments\"").WillReturnError(errors.New("database connection error"))
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
		{
			name: "Unique Constraint Violation",
			args: args{
				m: &model.Comment{ArticleID: 1, UserID: 1, Body: "Test comment"},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				mock.ExpectExec("INSERT INTO \"comments\"").WillReturnError(errors.New("unique constraint violation"))
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
		{
			name: "Nil Comment Object",
			args: args{
				m: nil,
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
		{
			name: "Empty Comment Object",
			args: args{
				m: &model.Comment{},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
		{
			name: "Comment with Invalid Fields",
			args: args{
				m: &model.Comment{ArticleID: 0, UserID: 0, Body: "Invalid comment"},
			},
			setupMockDB: func() (*gorm.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, nil, err
				}
				gdb, err := gorm.Open("postgres", db)
				return gdb, mock, err
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := tt.setupMockDB()
			if err != nil {
				t.Fatalf("failed to setup mock DB: %v", err)
			}
			defer db.Close()

			s := &ArticleStore{db: db}

			if err := s.CreateComment(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

/*
ROOST_METHOD_HASH=ArticleStore_Delete_8daad9ff19
ROOST_METHOD_SIG_HASH=ArticleStore_Delete_0e09651031

FUNCTION_DEF=func (s *ArticleStore) Delete(m *model.Article) error // Delete deletes an article
*/
func TestArticleStoreDelete(t *testing.T) {

	tests := []struct {
		name        string
		setup       func() (*ArticleStore, error)
		article     *model.Article
		expectedErr error
	}{
		{
			name: "Successful Deletion",
			setup: func() (*ArticleStore, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\"").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: nil,
		},
		{
			name: "Article Not Found",
			setup: func() (*ArticleStore, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\"").
					WithArgs(1).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Connection Error",
			setup: func() (*ArticleStore, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				db.Close()
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: errors.New("sql: database is closed"),
		},
		{
			name: "Invalid Article Data",
			setup: func() (*ArticleStore, error) {
				db, _, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     nil,
			expectedErr: errors.New("invalid input"),
		},
		{
			name: "Concurrent Deletion",
			setup: func() (*ArticleStore, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				for i := 1; i <= 5; i++ {
					mock.ExpectBegin()
					mock.ExpectExec("DELETE FROM \"articles\"").
						WithArgs(i).
						WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				}
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: nil,
		},
		{
			name: "Deletion of Last Article",
			setup: func() (*ArticleStore, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\"").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: nil,
		},
		{
			name: "Deletion with Foreign Key Constraints",
			setup: func() (*ArticleStore, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					return nil, err
				}
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM \"articles\"").
					WithArgs(1).
					WillReturnError(errors.New("foreign key constraint"))
				mock.ExpectRollback()
				gormDB, err := gorm.Open("postgres", db)
				if err != nil {
					return nil, err
				}
				return &ArticleStore{db: gormDB}, nil
			},
			article:     &model.Article{Model: gorm.Model{ID: 1}},
			expectedErr: errors.New("foreign key constraint"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			store, err := tt.setup()
			if err != nil {
				t.Fatalf("Failed to set up test: %v", err)
			}

			if tt.name == "Invalid Article Data" && tt.article == nil {
				err := store.Delete(tt.article)
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("Delete() error = %v, expectedErr %v", err, tt.expectedErr)
				}
				return
			}

			if tt.name == "Concurrent Deletion" {
				var wg sync.WaitGroup
				for i := 1; i <= 5; i++ {
					wg.Add(1)
					go func(id uint) {
						defer wg.Done()
						err := store.Delete(&model.Article{Model: gorm.Model{ID: id}})
						if err != nil {
							t.Errorf("Delete() error = %v", err)
						}
					}(uint(i))
				}
				wg.Wait()
				return
			}

			err = store.Delete(tt.article)
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Delete() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			t.Logf("Test %s passed", tt.name)
		})
	}
}

/*
ROOST_METHOD_HASH=ArticleStore_IsFavorited_799826fee5
ROOST_METHOD_SIG_HASH=ArticleStore_IsFavorited_f6d5e67492

FUNCTION_DEF=func (s *ArticleStore) IsFavorited(a *model.Article, u *model.User) (bool, error) // IsFavorited returns whether the article is favorited by the user
*/
func TestArticleStoreIsFavorited(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("sqlite3", db)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	store := &ArticleStore{db: gormDB}

	tests := []struct {
		name      string
		article   *model.Article
		user      *model.User
		mockSetup func()
		want      bool
		wantErr   bool
	}{
		{
			name:      "Article and User are nil",
			article:   nil,
			user:      nil,
			mockSetup: func() {},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Article is nil, User is valid",
			article:   nil,
			user:      &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func() {},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Article is valid, User is nil",
			article:   &model.Article{Model: gorm.Model{ID: 1}},
			user:      nil,
			mockSetup: func() {},
			want:      false,
			wantErr:   false,
		},
		{
			name:    "Article is favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Article is not favorited by the user",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "Database error occurs",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("database error"))
			},
			want:    false,
			wantErr: true,
		},
		{
			name:    "Invalid database schema",
			article: &model.Article{Model: gorm.Model{ID: 1}},
			user:    &model.User{Model: gorm.Model{ID: 1}},
			mockSetup: func() {
				mock.ExpectQuery("SELECT count").
					WithArgs(1, 1).
					WillReturnError(fmt.Errorf("no such table: favorite_articles"))
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.mockSetup()

			got, err := store.IsFavorited(tt.article, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.IsFavorited() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArticleStore.IsFavorited() = %v, want %v", got, tt.want)
			}

			if tt.wantErr && err == nil {
				t.Errorf("Expected error, but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Did not expect error, but got %v", err)
			}

			t.Logf("Test %q passed successfully", tt.name)
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
