package store

import (
	debug "runtime/debug"
	testing "testing"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	time "time"
	assert "github.com/stretchr/testify/assert"
	sql "database/sql"
)








/*
ROOST_METHOD_HASH=NewArticleStore_85784abca5
ROOST_METHOD_SIG_HASH=NewArticleStore_436ae9c986

FUNCTION_DEF=func NewArticleStore(db *gorm.DB) *ArticleStore // NewArticleStore returns a new ArticleStore


*/
func TestNewArticleStore(t *testing.T) {
	tests := []struct {
		name    string
		db      *gorm.DB
		wantErr bool
	}{
		{
			name:    "Successfully create new article store",
			db:      nil,
			wantErr: false,
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

			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock DB: %v", err)
			}
			defer mockDB.Close()

			gormDB, err := gorm.Open("mysql", mockDB)
			if err != nil {
				t.Fatalf("Failed to create GORM DB: %v", err)
			}
			defer gormDB.Close()

			t.Log("Creating new article store with mock DB")
			store := NewArticleStore(gormDB)

			if store == nil {
				t.Error("Expected non-nil ArticleStore, got nil")
			}

			if store.db != gormDB {
				t.Error("ArticleStore DB connection doesn't match provided connection")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			t.Log("Successfully validated article store creation")
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Create_1273475ade
ROOST_METHOD_SIG_HASH=ArticleStore_Create_a27282cad5

FUNCTION_DEF=func (s *ArticleStore) Create(m *model.Article) error // Create creates an article


*/
func TestArticleStoreCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		article *model.Article
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful Article Creation",
			article: &model.Article{
				Title:       "Test Article",
				Slug:        "test-article",
				Description: "Test Description",
				Body:        "Test Body Content",
				AuthorID:    1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `articles`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			tt.mockFn(mock)

			err := store.Create(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if !tt.wantErr {
				t.Log("Successfully created article:", tt.article.Title)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_CreateComment_b16d4a71d4
ROOST_METHOD_SIG_HASH=ArticleStore_CreateComment_7475736b06

FUNCTION_DEF=func (s *ArticleStore) CreateComment(m *model.Comment) error // CreateComment creates a comment of the article


*/
func TestArticleStoreCreateComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name    string
		comment *model.Comment
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successful comment creation",
			comment: &model.Comment{
				Body:      "Test comment",
				ArticleID: 1,
				UserID:    1,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `comments`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
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

			t.Logf("Testing scenario: %s", tt.name)

			tt.mockFn(mock)

			err := store.CreateComment(tt.comment)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			if err == nil {
				t.Log("Successfully created comment")
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

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Article Deletion",
			article: &model.Article{
				ID:      1,
				Title:   "Test Article",
				Content: "Test Content",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `articles`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
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

			tt.setupMock(mock)

			err := store.Delete(tt.article)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case '%s' completed successfully", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_DeleteComment_effbcb38aa
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteComment_d3c99623e4

FUNCTION_DEF=func (s *ArticleStore) DeleteComment(m *model.Comment) error // DeleteComment deletes an comment


*/
func TestArticleStoreDeleteComment(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		comment     *model.Comment
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful comment deletion",
			comment: &model.Comment{
				ID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Failed comment deletion - DB error",
			comment: &model.Comment{
				ID: 2,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `comments`").
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			expectError: true,
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

			t.Logf("Running test case: %s", tt.name)

			tt.setupMock(mock)

			err := store.DeleteComment(tt.comment)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetTags_45f5cdc4bb
ROOST_METHOD_SIG_HASH=ArticleStore_GetTags_fb0aefcdd2

FUNCTION_DEF=func (s *ArticleStore) GetTags() ([ // GetTags creates a article tag
]model.Tag, error) 

*/
func TestArticleStoreGetTags(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedTags  []model.Tag
		expectedError error
	}{
		{
			name: "Successfully retrieve tags",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "tag"}).
					AddRow(1, "golang").
					AddRow(2, "testing").
					AddRow(3, "programming")
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
			},
			expectedTags: []model.Tag{
				{ID: 1, Tag: "golang"},
				{ID: 2, Tag: "testing"},
				{ID: 3, Tag: "programming"},
			},
			expectedError: nil,
		},
		{
			name: "Empty tags list",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "tag"})
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnRows(rows)
			},
			expectedTags:  []model.Tag{},
			expectedError: nil,
		},
		{
			name: "Database error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `tags`").WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedTags:  []model.Tag{},
			expectedError: gorm.ErrRecordNotFound,
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

			tt.mockSetup(mock)

			tags, err := store.GetTags()

			if err != tt.expectedError {
				t.Errorf("GetTags() error = %v, expected error %v", err, tt.expectedError)
				return
			}

			if len(tags) != len(tt.expectedTags) {
				t.Errorf("GetTags() returned %d tags, expected %d tags", len(tags), len(tt.expectedTags))
				return
			}

			for i, tag := range tags {
				if tag.ID != tt.expectedTags[i].ID || tag.Tag != tt.expectedTags[i].Tag {
					t.Errorf("Tag at index %d mismatch: got %+v, expected %+v", i, tag, tt.expectedTags[i])
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}

			t.Logf("Successfully completed test case: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetByID_6fe18728fc
ROOST_METHOD_SIG_HASH=ArticleStore_GetByID_bb488e542f

FUNCTION_DEF=func (s *ArticleStore) GetByID(id uint) (*model.Article, error) // GetByID finds an article from id


*/
func TestArticleStoreGetById(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name          string
		articleID     uint
		mockSetup     func(sqlmock.Sqlmock)
		expectedError bool
		validate      func(*testing.T, *model.Article, error)
	}{
		{
			name:      "Successfully retrieve existing article",
			articleID: 1,
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description"}).
						AddRow(1, "Test Article", "Test Description"))

				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
						AddRow(1, "tag1").
						AddRow(2, "tag2"))

				mock.ExpectQuery("SELECT").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email"}).
						AddRow(1, "testuser", "test@example.com"))
			},
			expectedError: false,
			validate: func(t *testing.T, article *model.Article, err error) {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
						t.Fail()
					}
				}()

				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				if article == nil {
					t.Error("Expected article to not be nil")
					return
				}

				if article.ID != 1 {
					t.Errorf("Expected article ID 1, got %d", article.ID)
				}

				if article.Tags == nil {
					t.Error("Expected Tags to be preloaded")
				}

				if article.Author == nil {
					t.Error("Expected Author to be preloaded")
				}

				t.Log("Successfully validated article retrieval with preloaded relationships")
			},
		},
		{
			name:      "Article not found",
			articleID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").
					WithArgs(999).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: true,
			validate: func(t *testing.T, article *model.Article, err error) {
				if err == nil {
					t.Error("Expected error for non-existent article")
					return
				}
				if article != nil {
					t.Error("Expected nil article for non-existent ID")
				}
				t.Log("Successfully validated not found scenario")
			},
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

			tt.mockSetup(mock)

			article, err := store.GetByID(tt.articleID)

			tt.validate(t, article, err)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_Update_3cddacb803
ROOST_METHOD_SIG_HASH=ArticleStore_Update_e245edd177

FUNCTION_DEF=func (s *ArticleStore) Update(m *model.Article) error // Update updates an article


*/
func TestArticleStoreUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		article     *model.Article
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful Article Update",
			article: &model.Article{
				ID:     1,
				Title:  "Updated Title",
				Body:   "Updated content body",
				UserID: 123,
				Slug:   "updated-title",
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `articles`").
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						1,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
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

			tt.setupMock(mock)

			err := store.Update(tt.article)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				t.Logf("Got expected error: %v", err)
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				t.Log("Article updated successfully")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetComments_d7c78dda64
ROOST_METHOD_SIG_HASH=ArticleStore_GetComments_af08ddd59e

FUNCTION_DEF=func (s *ArticleStore) GetComments(m *model.Article) ([ // GetComments gets coments of the article
]model.Comment, error) 

*/
func TestArticleStoreGetComments(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open gorm connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name         string
		article      *model.Article
		setupMock    func(sqlmock.Sqlmock)
		wantComments []model.Comment
		wantErr      bool
	}{
		{
			name: "Successfully retrieve comments for article",
			article: &model.Article{
				ID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery("SELECT (.+) FROM `comments`").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "article_id", "body", "author_id"}).
						AddRow(1, 1, "Test comment 1", 10).
						AddRow(2, 1, "Test comment 2", 11))

				mock.ExpectQuery("SELECT (.+) FROM `users`").
					WithArgs(10, 11).
					WillReturnRows(sqlmock.NewRows([]string{"id", "username"}).
						AddRow(10, "author1").
						AddRow(11, "author2"))
			},
			wantComments: []model.Comment{
				{
					ID:        1,
					ArticleID: 1,
					Body:      "Test comment 1",
					Author: model.User{
						ID:       10,
						Username: "author1",
					},
				},
				{
					ID:        2,
					ArticleID: 1,
					Body:      "Test comment 2",
					Author: model.User{
						ID:       11,
						Username: "author2",
					},
				},
			},
			wantErr: false,
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

			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			gotComments, err := store.GetComments(tt.article)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetComments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(gotComments) != len(tt.wantComments) {
					t.Errorf("ArticleStore.GetComments() returned %d comments, want %d",
						len(gotComments), len(tt.wantComments))
					return
				}

				for i, comment := range gotComments {
					if comment.ArticleID != tt.article.ID {
						t.Errorf("Comment[%d].ArticleID = %v, want %v", i, comment.ArticleID, tt.article.ID)
					}
					if comment.Author.ID == 0 {
						t.Error("Author information not properly preloaded")
					}
				}

				t.Logf("Successfully verified %d comments for article ID %d",
					len(gotComments), tt.article.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
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
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expected    bool
		expectError bool
	}{
		{
			name: "Valid Article and User with Favorite Relationship",
			article: &model.Article{
				ID: 1,
			},
			user: &model.User{
				ID: 1,
			},
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `favorite_articles` WHERE \\(article_id = \\? AND user_id = \\?\\)").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:    true,
			expectError: false,
		},
		{
			name:    "Nil Article",
			article: nil,
			user: &model.User{
				ID: 1,
			},
			setupMock: func(m sqlmock.Sqlmock) {

			},
			expected:    false,
			expectError: false,
		},
		{
			name: "Nil User",
			article: &model.Article{
				ID: 1,
			},
			user: nil,
			setupMock: func(m sqlmock.Sqlmock) {

			},
			expected:    false,
			expectError: false,
		},
		{
			name: "DB Error",
			article: &model.Article{
				ID: 1,
			},
			user: &model.User{
				ID: 1,
			},
			setupMock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery("SELECT COUNT\\(\\*\\) FROM `favorite_articles` WHERE \\(article_id = \\? AND user_id = \\?\\)").
					WithArgs(1, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expected:    false,
			expectError: true,
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

			t.Logf("Running test case: %s", tt.name)

			tt.setupMock(mock)

			result, err := store.IsFavorited(tt.article, tt.user)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if result != tt.expected {
				t.Errorf("Expected result %v but got %v", tt.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}

			t.Logf("Test case completed successfully: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_GetFeedArticles_a37e1934b6
ROOST_METHOD_SIG_HASH=ArticleStore_GetFeedArticles_f5f09c020e

FUNCTION_DEF=func (s *ArticleStore) GetFeedArticles(userIDs [ // GetFeedArticles returns following users' articles
]uint, limit, offset int64) ([]model.Article, error) 

*/
func TestArticleStoreGetFeedArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name      string
		userIDs   []uint
		limit     int64
		offset    int64
		mockSetup func(sqlmock.Sqlmock)
		want      []model.Article
		wantErr   bool
	}{
		{
			name:    "Successfully retrieve articles for multiple users",
			userIDs: []uint{1, 2},
			limit:   10,
			offset:  0,
			mockSetup: func(mock sqlmock.Sqlmock) {
				expectedSQL := "SELECT * FROM `articles` WHERE user_id in (?,?) LIMIT 10 OFFSET 0"
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "title", "description", "body", "created_at", "updated_at",
				}).AddRow(
					1, 1, "Test Article 1", "Description 1", "Body 1", time.Now(), time.Now(),
				).AddRow(
					2, 2, "Test Article 2", "Description 2", "Body 2", time.Now(), time.Now(),
				)

				mock.ExpectQuery(expectedSQL).
					WithArgs(1, 2).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{
					"id", "username", "email",
				}).AddRow(
					1, "user1", "user1@example.com",
				).AddRow(
					2, "user2", "user2@example.com",
				)

				mock.ExpectQuery("SELECT * FROM `users`").
					WillReturnRows(authorRows)
			},
			want: []model.Article{
				{
					ID:          1,
					UserID:      1,
					Title:       "Test Article 1",
					Description: "Description 1",
					Body:        "Body 1",
					Author: &model.User{
						ID:       1,
						Username: "user1",
						Email:    "user1@example.com",
					},
				},
				{
					ID:          2,
					UserID:      2,
					Title:       "Test Article 2",
					Description: "Description 2",
					Body:        "Body 2",
					Author: &model.User{
						ID:       2,
						Username: "user2",
						Email:    "user2@example.com",
					},
				},
			},
			wantErr: false,
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

			tt.mockSetup(mock)

			got, err := store.GetFeedArticles(tt.userIDs, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tt.want), len(got), "Expected %d articles, got %d", len(tt.want), len(got))

			for i, article := range got {
				assert.Equal(t, tt.want[i].ID, article.ID)
				assert.Equal(t, tt.want[i].UserID, article.UserID)
				assert.Equal(t, tt.want[i].Title, article.Title)
				assert.Equal(t, tt.want[i].Description, article.Description)
				assert.Equal(t, tt.want[i].Body, article.Body)

				assert.NotNil(t, article.Author)
				assert.Equal(t, tt.want[i].Author.ID, article.Author.ID)
				assert.Equal(t, tt.want[i].Author.Username, article.Author.Username)
				assert.Equal(t, tt.want[i].Author.Email, article.Author.Email)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			t.Logf("Successfully tested scenario: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_AddFavorite_9460fca478
ROOST_METHOD_SIG_HASH=ArticleStore_AddFavorite_c13a109f91

FUNCTION_DEF=func (s *ArticleStore) AddFavorite(a *model.Article, u *model.User) error // AddFavorite favorite an article


*/
func TestArticleStoreAddFavorite(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		article     *model.Article
		user        *model.User
		setupMock   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name: "Successful favorite addition",
			article: &model.Article{
				ID:             1,
				FavoritesCount: 0,
			},
			user: &model.User{
				ID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("UPDATE").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "Failed association append",
			article: &model.Article{
				ID:             1,
				FavoritesCount: 0,
			},
			user: &model.User{
				ID: 1,
			},
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()

				mock.ExpectExec("INSERT INTO").
					WithArgs(1, 1).
					WillReturnError(sql.ErrConnDone)

				mock.ExpectRollback()
			},
			expectError: true,
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

			tt.setupMock(mock)

			initialCount := tt.article.FavoritesCount
			err := store.AddFavorite(tt.article, tt.user)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}

				if tt.article.FavoritesCount != initialCount {
					t.Errorf("FavoritesCount should not change on error, got %d, want %d",
						tt.article.FavoritesCount, initialCount)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if tt.article.FavoritesCount != initialCount+1 {
					t.Errorf("FavoritesCount not incremented, got %d, want %d",
						tt.article.FavoritesCount, initialCount+1)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=ArticleStore_DeleteFavorite_29c18a04a8
ROOST_METHOD_SIG_HASH=ArticleStore_DeleteFavorite_53deb5e792

FUNCTION_DEF=func (s *ArticleStore) DeleteFavorite(a *model.Article, u *model.User) error // DeleteFavorite unfavorite an article


*/
func TestArticleStoreDeleteFavorite(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	t.Run("Successful Deletion of Article Favorite", func(t *testing.T) {

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		testArticle := &model.Article{
			ID:             1,
			FavoritesCount: 2,
		}
		testUser := &model.User{
			ID: 1,
		}

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM `article_favorites`").
			WithArgs(testArticle.ID, testUser.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("UPDATE `articles`").
			WithArgs(1, testArticle.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		err := store.DeleteFavorite(testArticle, testUser)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if testArticle.FavoritesCount != 1 {
			t.Errorf("Expected favorites count to be 1, got %d", testArticle.FavoritesCount)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}

		t.Log("Successfully tested article favorite deletion")
	})

	t.Run("Failed Association Delete", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
				t.Fail()
			}
		}()

		testArticle := &model.Article{
			ID:             1,
			FavoritesCount: 2,
		}
		testUser := &model.User{
			ID: 1,
		}

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `article_favorites`").
			WithArgs(testArticle.ID, testUser.ID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		err := store.DeleteFavorite(testArticle, testUser)

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if testArticle.FavoritesCount != 2 {
			t.Errorf("FavoritesCount should not have changed, expected 2, got %d", testArticle.FavoritesCount)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("Unfulfilled expectations: %v", err)
		}

		t.Log("Successfully tested failed association delete")
	})

}


/*
ROOST_METHOD_HASH=ArticleStore_GetArticles_101b7250e8
ROOST_METHOD_SIG_HASH=ArticleStore_GetArticles_91bc0a6760

FUNCTION_DEF=func (s *ArticleStore) GetArticles(tagName, username string, favoritedBy *model.User, limit, offset int64) ([ // GetArticles get global articles
]model.Article, error) 

*/
func TestArticleStoreGetArticles(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open("mysql", db)
	if err != nil {
		t.Fatalf("Failed to open GORM connection: %v", err)
	}
	defer gormDB.Close()

	store := &ArticleStore{
		db: gormDB,
	}

	tests := []struct {
		name        string
		tagName     string
		username    string
		favoritedBy *model.User
		limit       int64
		offset      int64
		setupMock   func(sqlmock.Sqlmock)
		wantErr     bool
		wantLen     int
	}{
		{
			name:        "Get All Articles Without Filters",
			tagName:     "",
			username:    "",
			favoritedBy: nil,
			limit:       10,
			offset:      0,
			setupMock: func(mock sqlmock.Sqlmock) {

				rows := sqlmock.NewRows([]string{"id", "title", "description", "body", "user_id"}).
					AddRow(1, "Test Article 1", "Description 1", "Body 1", 1).
					AddRow(2, "Test Article 2", "Description 2", "Body 2", 2)

				mock.ExpectQuery(`SELECT \* FROM "articles"`).
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "username", "email"}).
					AddRow(1, "user1", "user1@example.com").
					AddRow(2, "user2", "user2@example.com")

				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(authorRows)
			},
			wantErr: false,
			wantLen: 2,
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

			tt.setupMock(mock)

			got, err := store.GetArticles(tt.tagName, tt.username, tt.favoritedBy, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("ArticleStore.GetArticles() returned %d articles, want %d", len(got), tt.wantLen)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled mock expectations: %v", err)
			}

			t.Logf("Successfully completed test case: %s", tt.name)
		})
	}
}

