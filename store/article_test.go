package store

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"database/sql"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"fmt"
	"time"
	"github.com/stretchr/testify/suite"
)


var dummyTags = []model.Tag{

var mockDBError = errors.New("mock db error")
/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92


 */
func TestNewArticleStore(t *testing.T) {
	type args struct {
		db *gorm.DB
	}

	tests := []struct {
		name string
		args args
	}{
		{
			"Normal Scenario - Valid Input",
			args{
				db: createMockDB(t),
			},
		},
		{
			"Edge Case - Nil Input",
			args{
				db: nil,
			},
		},
		{
			"Validating Error Handling - Initialized DB",
			args{
				db: createMockDBWithInitializedValues(t),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)
			as := NewArticleStore(tt.args.db)
			assert.IsType(t, &ArticleStore{}, as)
		})
	}
}

func createMockDB(t *testing.T) *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	mock.ExpectPing()
	mockDB, _ := gorm.Open("sqlmock", db)
	return mockDB
}

func createMockDBWithInitializedValues(t *testing.T) *gorm.DB {
	mockDB := createMockDB(t)

	return mockDB
}


/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestCreate(t *testing.T) {

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while setting up mock database: %v", err)
	}

	gormDB, _ := gorm.Open("postgres", db)
	defer gormDB.Close()

	tests := []struct {
		name            string
		expectedArticle model.Article
		expectedError   error
		finderFunc      func() *store.ArticleStore
	}{
		{
			name: "Scenario 1: Successful Article Creation",
			expectedArticle: model.Article{
				Title:  "new article",
				Body:   "new body",
				UserID: 1,
			},
			expectedError: nil,
			finderFunc: func() *store.ArticleStore {
				return store.NewArticleStore(gormDB)
			},
		},
		{
			name: "Scenario 2: Attempt to Create Article with Missing Required Fields",
			expectedArticle: model.Article{
				Title:  "",
				Body:   "",
				UserID: 0,
			},
			expectedError: gorm.ErrRecordNotFound,
			finderFunc: func() *store.ArticleStore {
				return store.NewArticleStore(gormDB)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.finderFunc()
			err := as.Create(&tt.expectedArticle)

			if err != tt.expectedError {
				t.Errorf("expected error = %v, got = %v", tt.expectedError, err)
				return
			}

			if tt.expectedError == nil {
				var fetched model.Article
				as.db.First(&fetched)

				if fetched.Title != tt.expectedArticle.Title {
					t.Errorf("expected title = %v, got = %v", tt.expectedArticle.Title, fetched.Title)
				}
				if fetched.Body != tt.expectedArticle.Body {
					t.Errorf("expected body = %v, got = %v", tt.expectedArticle.Body, fetched.Body)
				}
				if fetched.UserID != tt.expectedArticle.UserID {
					t.Errorf("expected userID = %v, got = %v", tt.expectedArticle.UserID, fetched.UserID)
				}
			}
		})
	}

	t.Run("Scenario 3: Attempt to Create Article when Database is Unreachable", func(t *testing.T) {
		gormDB.Close()

		as := store.NewArticleStore(gormDB)
		err := as.Create(&model.Article{
			Title:  "new article",
			Body:   "new body",
			UserID: 1,
		})

		if err == nil {
			t.Errorf("expected an error, got nil")
		}
	})
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestArticleStoreCreateComment(t *testing.T) {

	var (
		mock         sqlmock.Sqlmock
		db           *sql.DB
		articleStore *ArticleStore
		err          error
	)

	testCases := []struct {
		name    string
		comment *model.Comment
		mock    func()
		isError bool
	}{
		{

			name: "When a new comment is created successfully, the function should return without any error",
			comment: &model.Comment{
				Body:      "This is a sample comment",
				UserID:    1,
				ArticleID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "comments"`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			isError: false,
		},
		{

			name:    "When the comment instance is null, the function should return an error",
			comment: nil,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "comments"`).WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectRollback()
			},
			isError: true,
		},
		{

			name: "When there's a database error while creating a comment, the function should return an error",
			comment: &model.Comment{
				Body:      "This is a sample comment",
				UserID:    1,
				ArticleID: 1,
			},
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO "comments"`).WillReturnError(gorm.ErrCantStartTransaction)
				mock.ExpectRollback()
			},
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, mock, _ = sqlmock.New()
			gormDB, _ := gorm.Open("mysql", db)

			articleStore = &ArticleStore{gormDB}
			tc.mock()

			err = articleStore.CreateComment(tc.comment)

			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())

		})
	}
}


/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func TestDelete(t *testing.T) {
	testCases := []struct {
		desc        string
		id          int64
		expectedErr string
	}{
		{
			desc: "Delete an existing article",
			id:   1,
		},
		{
			desc:        "Delete a non-existing article",
			id:          2,
			expectedErr: gorm.ErrRecordNotFound.Error(),
		},
		{
			desc:        "Delete operation when the database connection is down",
			id:          1,
			expectedErr: "failed to connect to database",
		},
		{
			desc: "Delete operation on an article written by a specific author",
			id:   3,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			gormdb, _ := gorm.Open("sqlmock", db)
			defer func() {
				_ = gormdb.Close()
			}()

			switch tC.desc {
			case "Delete an existing article", "Delete operation on an article written by a specific author":
				mock.ExpectExec("DELETE").WithArgs(tC.id).WillReturnResult(sqlmock.NewResult(1, 1))
			case "Delete a non-existing article":
				mock.ExpectExec("DELETE").WithArgs(tC.id).WillReturnError(gorm.ErrRecordNotFound)
			case "Delete operation when the database connection is down":
				_ = db.Close()
			}

			store := &ArticleStore{gormdb}
			article := &model.Article{Model: gorm.Model{ID: uint(tC.id)}}
			err := store.Delete(article)

			if err != nil && err.Error() != tC.expectedErr {
				t.Errorf("expected %v, but got %v", tC.expectedErr, err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {
	t.Run("successful comment deletion", func(t *testing.T) {
		store, mock := newMockArticleStore()
		comment := newMockComment()

		mock.ExpectBegin()
		mock.ExpectExec(
			fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" = $1", "comments", "id")).WithArgs(comment.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := store.DeleteComment(comment)

		assert.NoError(t, err)
	})

	t.Run("deleting non-existing comment", func(t *testing.T) {
		store, mock := newMockArticleStore()
		comment := newMockComment()

		mock.ExpectBegin()
		mock.ExpectExec(
			fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" = $1", "comments", "id")).WithArgs(comment.ID).
			WillReturnResult(sqlmock.NewResult(1, 0))
		mock.ExpectCommit()

		err := store.DeleteComment(comment)

		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("deleting comment with invalid input", func(t *testing.T) {
		store, mock := newMockArticleStore()
		comment := newMockComment()
		comment.ID = 0

		mock.ExpectBegin()
		mock.ExpectExec(
			fmt.Sprintf("DELETE FROM \"%s\" WHERE \"%s\" = $1", "comments", "id")).WithArgs(comment.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := store.DeleteComment(comment)

		assert.Error(t, err)
	})
}

func newMockArticleStore() (*ArticleStore, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("postgres", db)
	return &ArticleStore{
		db: gormDB,
	}, mock
}

func newMockComment() *model.Comment {
	return &model.Comment{
		Model: gorm.Model{
			ID:        uint(1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Body:      "mock_comment",
		UserID:    uint(1),
		ArticleID: uint(1),
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func TestArticleStoreGetCommentByID(t *testing.T) {

	testScenarios := []struct {
		description     string
		id              uint
		mockError       error
		expectedComment *model.Comment
		expectedError   error
	}{
		{
			description: "Get Comment by Valid ID",
			id:          1,
			mockError:   nil,
			expectedComment: &model.Comment{
				Body:      "Test Body",
				UserID:    1,
				ArticleID: 1,
			},
			expectedError: nil,
		},
		{
			description:     "Get Comment by Nonexistent ID",
			id:              2,
			mockError:       gorm.ErrRecordNotFound,
			expectedComment: nil,
			expectedError:   gorm.ErrRecordNotFound,
		},
		{
			description:     "Get Comment by Invalid ID (negative or zero)",
			id:              0,
			mockError:       nil,
			expectedComment: nil,
			expectedError:   errors.New("invalid ID"),
		},
		{
			description:     "DB Connection Failure",
			id:              1,
			mockError:       errors.New("connection error"),
			expectedComment: nil,
			expectedError:   errors.New("connection error"),
		},
	}

	for _, ts := range testScenarios {
		t.Run(ts.description, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to open mock db: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("Failed to open gorm db: %s", err)
			}
			as := &ArticleStore{db: gormDB}

			mock.ExpectQuery("^SELECT (.+) FROM \"comments\"*").
				WithArgs(ts.id).
				WillReturnError(ts.mockError)
			if ts.mockError == nil {
				mock.ExpectQuery("^SELECT (.+) FROM \"comments\"*").
					WithArgs(ts.id).
					WillReturnRows(sqlmock.NewRows([]string{"Body", "UserID", "ArticleID"}).
						AddRow(ts.expectedComment.Body, ts.expectedComment.UserID, ts.expectedComment.ArticleID))
			}

			resultComment, err := as.GetCommentByID(ts.id)

			assert.Equal(t, ts.expectedError, err)
			assert.Equal(t, ts.expectedComment, resultComment)
		})
	}
}


/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func TestArticleStoreGetTags(t *testing.T) {
	tests := []struct {
		name     string
		mock     func()
		wantTags []model.Tag
		wantErr  error
	}{
		{
			name: "retrieve tags successfully",
			mock: func() {
				rows := sqlmock.NewRows([]string{"Name"}).
					AddRow(dummyTags[0].Name).
					AddRow(dummyTags[1].Name).
					AddRow(dummyTags[2].Name)

				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"$").WillReturnRows(rows)
			},
			wantTags: dummyTags,
			wantErr:  nil,
		},
		{
			name: "db returns an error",
			mock: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"$").WillReturnError(mockDBError)
			},
			wantTags: []model.Tag{},
			wantErr:  mockDBError,
		},
		{
			name: "empty database",
			mock: func() {
				rows := sqlmock.NewRows([]string{"Name"})

				mock.ExpectQuery("^SELECT (.+) FROM \"tags\"$").WillReturnRows(rows)
			},
			wantTags: []model.Tag{},
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.mock()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			articleStore := &ArticleStore{gormDB}
			tags, err := articleStore.GetTags()

			if tt.wantErr != err {
				t.Errorf("GetTags() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalTags(tags, tt.wantTags) {
				t.Errorf("GetTags() = %v, want %v", tags, tt.wantTags)
			}
		})
	}
}

func equalTags(a, b []model.Tag) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
	}
	return true
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestArticleStoreGetByID(t *testing.T) {
   
    predefinedArticle := &model.Article {
        Title:        "Test article",
        Description:  "Test description",
        Body:         "Test body",
        UserID:       1,
    }

   
    testCases := []struct {
        name               string
        id                 uint
        predefinedDatabase map[uint]*model.Article
        expectedResult     *model.Article
        expectedError      error
    }{
        {
            "Retrieve an article by its ID",
            1,
            map[uint]*model.Article{1: predefinedArticle},
            predefinedArticle,
            nil,
        },
        {
            "Request ID does not exist in the database",
            2,
            map[uint]*model.Article{1: predefinedArticle},
            nil,
            gorm.ErrRecordNotFound,
        },
        {
            "Database is empty or not initialized",
            1,
            map[uint]*model.Article{},
            nil,
            gorm.ErrRecordNotFound,
        },
       
    }

    for _, testCase := range testCases {
        t.Run(testCase.name, func(t *testing.T) {
           
            db, mock, _ := sqlmock.New()
            gormDB, _ := gorm.Open("postgres", db)

           
            articleStore := ArticleStore{
                db: gormDB,
            }

           
            if len(testCase.predefinedDatabase) > 0 { 
                for id, article := range testCase.predefinedDatabase {
                    mock.ExpectQuery("^SELECT (.+) FROM \"model\" WHERE \"model\".\"id\" = $1").
                        WithArgs(id).
                        WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "body"}).
                            AddRow(id, article.Title, article.Description, article.Body))
                }
            }

           
            result, err := articleStore.GetByID(testCase.id)

           
            if testCase.expectedResult != nil && 
                (result.Title != testCase.expectedResult.Title && 
                result.Description != testCase.expectedResult.Description && 
                result.Body != testCase.expectedResult.Body) {
                t.Errorf("The returned article does not match the predefined article")
            }

            if testCase.expectedError != nil && fmt.Sprintf("%T", err) != fmt.Sprintf("%T", testCase.expectedError) {
                t.Errorf("Expected error %v but got %v", testCase.expectedError, err)
            }

            if testCase.expectedError == nil && err != nil {
                t.Errorf("Got unexpected error: %v", err)
            }
        })
    }
}


/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func TestArticleStoreUpdate(t *testing.T) {

	testCases := []struct {
		testName            string
		article             model.Article
		mockFunc            func(mock sqlmock.Sqlmock)
		expectedErr         bool
	} {
		{
			testName: "Successful Update of an Article",
			article: model.Article{
				Title: "First Blog",
				Description: "First Description",
				Body: "First Body",
				UserID: 1,
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE(.*)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: false,
		},
		{
			testName: "Update an Article with Incomplete Data",
			article: model.Article{
				Title: "",
				Body: "Incomplete Body",
				UserID: 1,
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE(.*)").WillReturnError(gorm.ErrRecordNotFound)
				mock.ExpectCommit()
			},
			expectedErr: true,
		},
		{
			testName: "Update an Article when Database is Unreachable",
			article: model.Article{
				Title: "Unreachable DB",
				Body: "Unreachable Body",
				UserID: 1,
			},
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE(.*)").WillReturnError(gorm.ErrCantStartTransaction)
				mock.ExpectCommit()
			},
			expectedErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
		
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error while mocking sql: %v", err)
			}
			gormDB, _ := gorm.Open("postgres", db)
			defer gormDB.Close()
			
			tt.mockFunc(mock)

			as := ArticleStore{db: gormDB}
			err = as.Update(&tt.article)

		
			if tt.expectedErr {
				t.Log("Expected Error")
				assert.Error(t, err)
			} else {
				t.Log("No error expected")
				assert.NoError(t, err)
			}
                })
	}
}


/*
ROOST_METHOD_HASH=GetComments_e24a0f1b73
ROOST_METHOD_SIG_HASH=GetComments_fa6661983e


 */
func TestArticleStoreGetComments(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(mock sqlmock.Sqlmock)
		expectedList []model.Comment
		expectError  bool
	}{
		{
			name: "Scenario 1: Function retrieves comments successfully",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "body", "user_id"}).AddRow(1, "Mock Comment", 1)
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \$1\).*`).WillReturnRows(rows)
			},
			expectedList: []model.Comment{
				{
					Model: gorm.Model{ID: 1},
					Body:  "Mock Comment",
					UserID: 1,
				},
			},
			expectError: false,
		},
		{
			name: "Scenario 2: Function handles an article without comments gracefully",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \$1\).*`).WillReturnRows(sqlmock.NewRows([]string{}))
			},
			expectedList: []model.Comment{},
			expectError:  false,
		},
		{
			name: "Scenario 3: Function returns an error when making an invalid database query",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "comments" WHERE \(article_id = \$1\).*`).WillReturnError(errors.New("mock error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.name)

			db, mock, _ := sqlmock.New()
			defer db.Close()
			gormDB, _ := gorm.Open("postgres", db)
			tt.setupMock(mock)

			store := &ArticleStore{db: gormDB}
			m := &model.Article{Model: gorm.Model{ID: 1}}

			comments, err := store.GetComments(m)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedList, comments)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func TestArticleStoreIsFavorited(t *testing.T) {
	type args struct {
		article *model.Article
		user    *model.User
		error   error
	}

	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "Scenario 1: Article and User both are not nil",
			args: args{
				article: &model.Article{},
				user:    &model.User{},
				error:   nil,
			},
			expected: true,
		},
		{
			name: "Scenario 2: Article is nil",
			args: args{
				article: nil,
				user:    &model.User{},
				error:   nil,
			},
			expected: false,
		},
		{
			name: "Scenario 3: User is nil",
			args: args{
				article: &model.Article{},
				user:    nil,
				error:   nil,
			},
			expected: false,
		},
		{
			name: "Scenario 4: Scenario of a database error",
			args: args{
				article: &model.Article{},
				user:    &model.User{},
				error:   errors.New("database error"),
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("error '%s' was not expected when opening gorm database", err)
			}

			s := &ArticleStore{db: gormDB}

			mock.ExpectQuery("SELECT count(*) FROM \"favorite_articles\" WHERE \"favorite_articles\".\"deleted_at\" IS NULL AND ((article_id = .* AND user_id = .*))").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

			if tt.args.error != nil {
				mock.ExpectQuery("SELECT count(*) FROM \"favorite_articles\" WHERE \"favorite_articles\".\"deleted_at\" IS NULL AND ((article_id = .* AND user_id = .*))").
					WillReturnError(tt.args.error)
			}

			val, err := s.IsFavorited(tt.args.article, tt.args.user)
			if tt.args.error != nil && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.expected != val {
				t.Errorf("expected %v, but got %v", tt.expected, val)
			}

		})
	}
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func TestRunSuiteArticleStoreGetFeedArticles(t *testing.T) {
	suite.Run(t, new(TestArticleStoreGetFeedArticles))
}


/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func TestArticleStoreAddFavorite(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("mysql", db)

	scenarios := []struct {
		name        string
		initial     *dbMock
		shouldError bool
	}{
		{
			name: "Standard AddFavorite invocation",
			initial: &dbMock{
				db:       gormDB,
				article:  &model.Article{},
				user:     &model.User{},
				favorite: 1,
			},
			shouldError: false,
		},
		{
			name: "Adding a favorite with a database error during association",
			initial: &dbMock{
				db:       gormDB,
				article:  &model.Article{},
				user:     &model.User{},
				favorite: 1,
				errorAssociation: errors.New("unexpected db error"),
			},
			shouldError: true,
		},
		{
			name: "Adding a favorite with a database error during updating the favorite count",
			initial: &dbMock{
				db:       gormDB,
				article:  &model.Article{},
				user:     &model.User{},
				favorite: 1,
				errorUpdate: errors.New("unexpected db error during updating"),
			},
			shouldError: true,
		},
		{
			name: "Providing a nil Article or User",
			initial: &dbMock{
				db:       gormDB,
				article:  nil,
				user:     nil,
				favorite: 1,
			},
			shouldError: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			store := ArticleStore{
				db: s.initial.db,
			}
			mock.ExpectBegin()

			err := store.AddFavorite(s.initial.article, s.initial.user)
			if s.shouldError {
				assert.Error(t, err)
				mock.ExpectRollback()
				t.Logf("Scenario: %s, fails as expected with error: %v", s.name, err)
			} else {
				assert.NoError(t, err)
				mock.ExpectCommit()
				t.Logf("Scenario: %s, succeeds as expected", s.name)
			}
		})
	}
}

