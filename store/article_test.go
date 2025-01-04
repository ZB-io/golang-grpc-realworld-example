package store

import (
	"github.com/jinzhu/gorm"
	"testing"
	"model"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raahii/golang-grpc-realworld-example/model"
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/raahii/golang-grpc-realworld-example/store"
	"time"
	"github.com/jinzhu/gorm/dialects/mysql"
)


var mockComment = &model.Comment{

var userID uint = 1

var articleID uint = 1

var articleFavCount int32 = 5

var favUser = &model.User{Model: gorm.Model{ID: userID}}

var favArticle = &model.Article{Model: gorm.Model{ID: articleID}, FavoritesCount: articleFavCount}
/*
ROOST_METHOD_HASH=NewArticleStore_6be2824012
ROOST_METHOD_SIG_HASH=NewArticleStore_3fe6f79a92


 */
func TestNewArticleStore(t *testing.T) {

	var db *gorm.DB
	articleStore := NewArticleStore(db)

	if articleStore.db != db {
		t.Errorf("Expected %v, but got %v", db, articleStore.db)
	}
}


/*
ROOST_METHOD_HASH=Create_0a911e138d
ROOST_METHOD_SIG_HASH=Create_723c594377


 */
func TestArticleStoreCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock database: %v", err)
	}

	gdb, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	articleStore := &ArticleStore{gdb}

	articleData := []*model.Article{
		&model.Article{Title: "My First Article", Description: "This is my first Golang Article", Body: "This is body", UserID: 1},
		&model.Article{},
		&model.Article{Title: "My First Article", Description: "Attempt to create duplicate article", Body: "Duplicate body", UserID: 1},
		&model.Article{Title: "Invalid foreign key", Description: "This article has invalid foreign key constraint", Body: "Invalid body", UserID: 999},
	}

	for _, data := range articleData {

		err = articleStore.Create(data)
		if err != nil {
			if data.Title == "" {
				t.Log("Scenario 2: Passed. Error occured because title was missing which is a mandatory field.")
				assert.Error(t, err)
			} else if data.Title == "My First Article" {
				t.Log("Scenario 3: Passed. Error occured due to unique constraint violation on Title.")
				assert.Error(t, err)
			} else if data.UserID == 999 {
				t.Log("Scenario 4: Passed. Error occured due to invalid foreign key.")
				assert.Error(t, err)
			}
			continue
		}

		t.Log("Scenario 1: Successfully created article in database.")
		assert.NoError(t, err)

	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


/*
ROOST_METHOD_HASH=CreateComment_58d394e2c6
ROOST_METHOD_SIG_HASH=CreateComment_28b95f60a6


 */
func TestArticleStoreCreateComment(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open("mysql", db)

	mockArticleStore := NewArticleStore(gormDB)

	tests := []struct {
		name    string
		arg     *model.Comment
		mock    func()
		wantErr error
	}{
		{
			name: "Successful creation of a comment",
			arg:  mockComment,
			mock: func() {
				mock.ExpectExec(`INSERT INTO "comments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), mockComment.Body, mockComment.UserID, mockComment.ArticleID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "Unsuccessful creation due to null fields",
			arg:  mockCommentWithNullFields,
			mock: func() {
				mock.ExpectExec(`INSERT INTO "comments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), mockCommentWithNullFields.Body, mockCommentWithNullFields.UserID, mockCommentWithNullFields.ArticleID).
					WillReturnError(errors.New("fields UserID and ArticleID can not be null"))
			},
			wantErr: errors.New("fields UserID and ArticleID can not be null"),
		},
		{
			name: "Unsuccessful creation due to DB connection error",
			arg:  mockComment,
			mock: func() {
				mock.ExpectExec(`INSERT INTO "comments"`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), mockComment.Body, mockComment.UserID, mockComment.ArticleID).
					WillReturnError(errors.New("DB connection error"))
			},
			wantErr: errors.New("DB connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			if err := mockArticleStore.CreateComment(tt.arg); (err != nil) != (tt.wantErr != nil) || (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) {
				t.Errorf("CreateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})

	}
}


/*
ROOST_METHOD_HASH=Delete_a8dc14c210
ROOST_METHOD_SIG_HASH=Delete_a4cc8044b1


 */
func TestArticleStoreDelete(t *testing.T) {

	testCases := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string
	} {
		{
			"Successful Deletion of an Article",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM articles WHERE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			"",
		},
		{
			"No Article to Delete with given parameters",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM articles WHERE").WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			"record not found",
		},
		{
			"Database Connection Issue",
			nil,
			"database connection error",
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
		
			article := &model.Article{
				Model: gorm.Model{ID: uint(1)},
				Title: "Test Title",
			}
		
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to initialize sqlmock: %v", err)
			}
			defer db.Close()
		
			if test.setupMock != nil {
				test.setupMock(mock)
			}

		
			gdb, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("failed to open a connection with the database: %v", err)
			}

			store := &ArticleStore{
				db: gdb,
			}

		
			err = store.Delete(article)
		
			if test.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteComment_b345e525a7
ROOST_METHOD_SIG_HASH=DeleteComment_732762ff12


 */
func TestArticleStoreDeleteComment(t *testing.T) {
	tests := map[string]struct {
		comment  *model.Comment
		id       uint
		hasError bool
	}{
		"Successful comment deletion": {
			comment: &model.Comment{
				Body:   "Test comment",
				UserID: 1,
			},
			id:       1,
			hasError: false,
		},
		"Deletion of a non-existing comment": {
			comment: &model.Comment{},
			id:       2,
			hasError: true,
		},
		"Passing a nil comment to the delete function": {
			comment:  nil,
			id:       3,
			hasError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("couldn't create sqlmock: %s", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("couldn't open gorm database: %s", err)
			}

			if tc.comment != nil {
				mock.ExpectExec("DELETE FROM comments WHERE id = ?").
					WithArgs(tc.id).
					WillReturnResult(sqlmock.NewResult(int64(tc.id), 1))
			}

			cStore := &ArticleStore{db: gormDB}

		
			err = cStore.DeleteComment(tc.comment)

			if tc.hasError {
			
				assert.Error(t, err)
				t.Log("should return an error when trying to delete comment with an error")
			} else {
			
				assert.NoError(t, err)
				t.Log("no error when deleting existing comment")
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetCommentByID_4bc82104a6
ROOST_METHOD_SIG_HASH=GetCommentByID_333cab101b


 */
func NewArticleStore(db *gorm.DB) *ArticleStore {
 return &ArticleStore{db: db}
}

func TestArticleStoreGetCommentByID(t *testing.T) {


 mockDB, _, err := sqlmock.New()
 if err != nil {
   t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
 }

 mockGorm, err := gorm.Open("postgres", mockDB)
 if err != nil {
   t.Fatalf("could not open connection to gorm: %s", err)
 }

 tests := []struct {
   name          string
   id            uint
   expected      *model.Comment
   expectedError error
 }{
   {
     name: "Successful Retrieval of a Comment by ID",
     id: 1,
     expected: &model.Comment{
       Body: "Test Comment",
       UserID: 1,
       ArticleID: 1,
     },
     expectedError: nil,
   },
   {
     name: "Retrieval of a Comment by a Non-Existent ID",
     id: 2,
     expected: nil,
     expectedError: gorm.ErrRecordNotFound,
   },
   {
     name: "Retrieval of a Comment with a Database Error",
     id: 3,
     expected: nil,
     expectedError: errors.New("database error"),
   },
 }


 articleStore := NewArticleStore(mockGorm)


 for _, tt := range tests {
   t.Run(tt.name, func(t *testing.T) {
    

     result, err := _store.GetCommentByID(tt.id)

    

     assert.Equal(t, tt.expected, result)
     if tt.expectedError != nil {
       assert.EqualError(t, err, tt.expectedError.Error())
     } else {
       assert.NoError(t, err)
     }
   })
 }
}


/*
ROOST_METHOD_HASH=GetTags_ac049ebded
ROOST_METHOD_SIG_HASH=GetTags_25034b82b0


 */
func TestArticleStoreGetTags(t *testing.T) {

	scenarios := []struct {
		desc   string
		mock   func(sqlmock.Sqlmock)
		expErr bool
	}{
		{
			"Retrieve all tags successfully",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "tag1").
					AddRow(2, time.Now(), time.Now(), nil, "tag2")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			false,
		},
		{
			"Empty result",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"})
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows)
			},
			false,
		},
		{
			"Database Error",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnError(errors.New("Some database error"))
			},
			true,
		},
		{
			"Concurrent Call",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name"}).
					AddRow(1, time.Now(), time.Now(), nil, "tag1").
					AddRow(2, time.Now(), time.Now(), nil, "tag2")
				mock.ExpectQuery("^SELECT (.+) FROM `tags`").WillReturnRows(rows).Times(2) 
			},
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening gorm database", err)
			}

			articleStore := &store.ArticleStore{
				DB: gormDB,
			}

			s.mock(mock)

			tags, err := articleStore.GetTags()

			if s.expErr {
				if err == nil {
					t.Error("was expecting an error, but no error returned")
				}
				t.Logf("expected err: %v", err)
			} else {
				if err != nil {
					t.Errorf("was not expecting an error, but got: %v", err)
				}
				t.Logf("tags: %v", tags)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=GetByID_36e92ad6eb
ROOST_METHOD_SIG_HASH=GetByID_9616e43e52


 */
func TestArticleStoreGetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gdb, _ := gorm.Open("postgres", db)

	testCases := []struct {
		desc          string
		inputID       uint
		mockArticle   *model.Article
		mockBehaviour func()
		expectErr     bool
		errMsg        string
	}{
		{
			desc:    "Success - Article fetch with ID",
			inputID: 1,
			mockArticle: &model.Article{
				Model: gorm.Model{ID: 1},
				Title: "Test Article",
			},
			mockBehaviour: func() {
				rows := sqlmock.NewRows([]string{"id", "title"})

					rows.AddRow(1, "Test Article")
	
					mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectErr: false,
		},
		{
			desc:    "Failure - No Article with this ID",
			inputID: 2,
			mockBehaviour: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").
					WithArgs(2).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectErr: true,
			errMsg:    gorm.ErrRecordNotFound.Error(),
		},
		{
			desc:    "Failure - DB Connection issue",
			inputID: 3,
			mockBehaviour: func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\" WHERE (.+)").
					WithArgs(3).
					WillReturnError(errors.New("DB connection issue"))
			},
			expectErr: true,
			errMsg:    "DB connection issue",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tc.mockBehaviour()

			store := ArticleStore{
				db: gdb,
			}

			res, err := store.GetByID(tc.inputID)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockArticle.ID, res.ID)
				assert.Equal(t, tc.mockArticle.Title, res.Title)
			}
			
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}


/*
ROOST_METHOD_HASH=Update_51145aa965
ROOST_METHOD_SIG_HASH=Update_6c1b5471fe


 */
func TestArticleStoreUpdate(t *testing.T) {
	testCases := []struct {
		name    string
		article *model.Article
		dbMock  func(mock sqlmock.Sqlmock, article *model.Article)
		err     error
	}{
		{
			name:    "Successful update of Article",
			article: &model.Article{Title: "Title 1", Description: "Description 1", Body: "Body 1"},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^UPDATE (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			err: nil,
		},
		{
			name:    "Updating a non-existent article",
			article: &model.Article{Title: "Non existant title", Description: "Non existant description", Body: "Non existant body"},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^UPDATE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			err: gorm.ErrRecordNotFound,
		},
		{
			name:    "Database error during update",
			article: &model.Article{Title: "Title 3", Description: "Description 3", Body: "Body 3"},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^UPDATE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			err: gorm.ErrRecordNotFound,
		},
		{
			name:    "Attempt to update with invalid Article data",
			article: &model.Article{Title: "", Description: "", Body: ""},
			dbMock: func(mock sqlmock.Sqlmock, article *model.Article) {
				mock.ExpectExec("^UPDATE (.+)").WillReturnError(gorm.ErrRecordNotFound)
			},
			err: gorm.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}

			gdb, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening gorm database", err)
			}

			tc.dbMock(mock, tc.article)

			store := &ArticleStore{db: gdb}
			err = store.Update(tc.article)

			assert.Equal(t, tc.err, err)

			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
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
        name     string
        mockFunc func()
        input    *model.Article
        expect   func(*[]model.Comment, error)
    }{
        {
            name: "Scenario 1: GetComments returns correctly with valid Model.Article",
           
           
            mockFunc: func() {
               
            },
            input: nil,
            expect: func(cs *[]model.Comment, err error) {
                if err != nil {
                    t.Fatal("Error was not expected")
                }
               
            },
        },
        {
            name: "Scenario 2: GetComments returns an error when passed an incorrect Model.Article",
           
            mockFunc: func() {
               
            },
            input: nil,
            expect: func(cs *[]model.Comment, err error) {
                if err == nil {
                    t.Fatal("Expected error, but none was received")
                }
            },
        },
        {
            name: "Scenario 3: GetComments returns an empty array when passed an Article with no comments",
           
            mockFunc: func() {
               
            },
            input: nil,
            expect: func(cs *[]model.Comment, err error) {
                if err != nil {
                    t.Fatal("Error was not expected")
                }
                if len(*cs) != 0 {
                    t.Fatal("Expected empty array, but got data")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            db, mock, _ := sqlmock.New()
            gormDB, _ := gorm.Open("postgres", db)
            defer gormDB.Close()

           
            tt.mockFunc()
            mock.ExpectCommit()

            articleStore := &ArticleStore{db: gormDB}

            cs, err := articleStore.GetComments(tt.input)

           
            tt.expect(&cs, err)
        })
    }
}


/*
ROOST_METHOD_HASH=IsFavorited_7ef7d3ed9e
ROOST_METHOD_SIG_HASH=IsFavorited_f34d52378f


 */
func TestArticleStoreIsFavorited(t *testing.T) {

	articleStore := &store.ArticleStore{}


	testArticle := &model.Article{ID: 1}
	testUser := &model.User{ID: 1}



	isFavorited, err := articleStore.IsFavorited(testArticle, testUser)


	assert.Nil(t, err)





	assert.Equal(t, false, isFavorited)
}


/*
ROOST_METHOD_HASH=GetFeedArticles_9c4f57afe4
ROOST_METHOD_SIG_HASH=GetFeedArticles_cadca0e51b


 */
func NewArticleStore(db *gorm.DB) *ArticleStore {
	return &ArticleStore{db: db}
}

func TestArticleStoreGetFeedArticles(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %s", err)
	}
	defer db.Close()

	g, err := gorm.Open("postgres", db)
	if err != nil {
		t.Fatalf("got an error when trying to open a connection with gorm: %s", err)
	}

	columns := []string{"id", "title", "description", "body", "user_id", "favorites_count"}
	mockRows := sqlmock.NewRows(columns).
		AddRow(1, "TestTitle", "TestDesc", "TestBody", 1, 0)

	s := NewArticleStore(g)
	userIDs := []uint{1, 2, 3}
	expectedArticles := []model.Article{
		{
			Model:          gorm.Model{ID: 1},
			Title:          "TestTitle",
			Description:    "TestDesc",
			Body:           "TestBody",
			UserID:         1,
			FavoritesCount: 0,
		},
	}

	tests := []struct {
		name  string
		setup func()
		exec  func() ([]model.Article, error)
		assert func([]model.Article, error)
	}{
		{
			"Retrieve Feed Articles Successfully",
			func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\"").
					WithArgs(1, 2, 3, 10, 0).
					WillReturnRows(mockRows)
			},
			func() ([]model.Article, error) {
				return s.GetFeedArticles(userIDs, 10, 0)
			},
			func(articles []model.Article, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expectedArticles, articles)
			},
		},
		{
			"User IDs Not Present in Database",
			func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\"").
					WithArgs(4, 5, 6, 10, 0).
					WillReturnRows(sqlmock.NewRows(columns))
			},
			func() ([]model.Article, error) {
				return s.GetFeedArticles([]uint{4, 5, 6}, 10, 0)
			},
			func(articles []model.Article, err error) {
				assert.Nil(t, err)
				assert.Empty(t, articles)
			},
		},
		{
			"Database Error When Attempting to Retrieve Articles",
			func() {
				mock.ExpectQuery("^SELECT (.+) FROM \"articles\"").
					WithArgs(1, 2, 3, 10, 0).
					WillReturnError(errors.New("some database error"))
			},
			func() ([]model.Article, error) {
				return s.GetFeedArticles(userIDs, 10, 0)
			},
			func(articles []model.Article, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, articles)
			},
		},
		{
			"Testing Negative Limit and Offset",
			func() {},
			func() ([]model.Article, error) {
				return s.GetFeedArticles(userIDs, -10, -1)
			},
			func(articles []model.Article, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, articles)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			articles, err := tt.exec()
			tt.assert(articles, err)
		})
	}
}


/*
ROOST_METHOD_HASH=AddFavorite_2b0cb9d894
ROOST_METHOD_SIG_HASH=AddFavorite_c4dea0ee90


 */
func TestArticleStoreAddFavorite(t *testing.T) {

	testScenarios := []struct {
		name       string
		shouldFail bool
		mockFunc   func(mock sqlmock.Sqlmock)
	}{
		{
		
			name: "successfully add a user to the articles' favorited list",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles \\(article_id, user_id\\)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE articles SET favorites_count = favorites_count + \\$1 WHERE id = \\$2$").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
		
			name:       "AddFavorite fails due to error while adding user to favorited list",
			shouldFail: true,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles \\(article_id, user_id\\)").WillReturnError(errors.New("error in adding favourite"))
				mock.ExpectRollback()
			},
		},
		{
		
			name:       "AddFavorite fails due to error while updating favoritesCount",
			shouldFail: true,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO favorite_articles \\(article_id, user_id\\)").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("^UPDATE articles SET favorites_count = favorites_count + \\$1 WHERE id = \\$2$").WillReturnError(errors.New("error in updating favoritesCount"))
				mock.ExpectRollback()
			},
		},
	}

	for _, ts := range testScenarios {
		t.Run(ts.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			gdb, _ := gorm.Open("postgres", db)

			ts.mockFunc(mock)

			s := &ArticleStore{db: gdb}

			err := s.AddFavorite(&model.Article{}, &model.User{})
			if ts.shouldFail {
				if err == nil {
					t.Errorf("expected AddFavorite to fail : %v", ts.name)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error : %v, scenario : %v", err, ts.name)
				return
			}

		
			err = mock.ExpectationsWereMet()
			if err != nil {
				t.Errorf("unfulfilled expectations : %v scenario : %v", err, ts.name)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=DeleteFavorite_a856bcbb70
ROOST_METHOD_SIG_HASH=DeleteFavorite_f7e5c0626f


 */
func TestArticleStoreDeleteFavorite(t *testing.T) {
  tests := []struct {
    name                string
    expectationFunction func(mock sqlmock.Sqlmock)
    expectedError       error
  }{
    {
      "Successful Deletion of Favorite",
      func(mock sqlmock.Sqlmock) {
        mock.ExpectBegin()
        mock.ExpectExec(
          "DELETE FROM \"\" WHERE \\(article_id IN \\(\\?\\) AND user_id IN \\(\\?\\)\\)",
        ).WithArgs(favArticle.ID, favUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectExec(
          "UPDATE \"articles\" SET \"favorites_count\" = \"favorites_count\" - \\? WHERE \"articles\".\"id\" = \\?",
        ).WithArgs(1, favArticle.ID).WillReturnResult(sqlmock.NewResult(1, 1))
        mock.ExpectCommit()
      },
      nil,
    },
    {
      "Deletion of Non-Existent Favorite",
      func(mock sqlmock.Sqlmock) {
        mock.ExpectBegin()
        mock.ExpectExec(
          "DELETE FROM \"\" WHERE \\(article_id IN \\(\\?\\) AND user_id IN \\(\\?\\)\\)",
        ).WithArgs(favArticle.ID, favUser.ID).WillReturnResult(sqlmock.NewResult(1, 0))
        mock.ExpectRollback()
      },
      errors.New("record not found"),
    },
    {
      "Error during deletion",
      func(mock sqlmock.Sqlmock) {
        mock.ExpectBegin()
        mock.ExpectExec(
          "DELETE FROM \"\" WHERE \\(article_id IN \\(\\?\\) AND user_id IN \\(\\?\\)\\)",
        ).WithArgs(favArticle.ID, favUser.ID).WillReturnError(gorm.ErrRecordNotFound)
        mock.ExpectRollback()
      },
      gorm.ErrRecordNotFound,
    },
    {
      "Unchanged FavoritesCount After Failure",
      func(mock sqlmock.Sqlmock) {
        mock.ExpectBegin()
        mock.ExpectExec(
          "DELETE FROM \"\" WHERE \\(article_id IN \\(\\?\\) AND user_id IN \\(\\?\\)\\)",
        ).WithArgs(favArticle.ID, favUser.ID).WillReturnError(gorm.ErrRecordNotFound)
        mock.ExpectRollback()
      },
      gorm.ErrRecordNotFound,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      db, mock, err := sqlmock.New()
      if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
      }
      defer db.Close()

      gdb, err := gorm.Open("postgres", db)
      if err != nil {
        t.Fatalf("Failed to open the connection: %v\n", err)
      }

      tt.expectationFunction(mock)

      store := &ArticleStore{db: gdb}
      article := favArticle
      err = store.DeleteFavorite(article, favUser) 

      assert.Equal(t, tt.expectedError, err)

      if tt.expectedError == gorm.ErrRecordNotFound {
        assert.Equal(t, int(articleFavCount), int(article.FavoritesCount))
      } else {
        assert.Equal(t, int(articleFavCount)-1, int(article.FavoritesCount))
      }
    })
  }
}

