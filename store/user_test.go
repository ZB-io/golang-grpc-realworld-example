package store

import (
		"reflect"
		"sync"
		"testing"
		"github.com/jinzhu/gorm"
		"github.com/raahii/golang-grpc-realworld-example/model"
		"github.com/stretchr/testify/assert"
		"errors"
		"github.com/stretchr/testify/mock"
		"time"
		"github.com/DATA-DOG/go-sqlmock"
)

type DB struct {
	sync.RWMutex
	Value        interface{}
	Error        error
	RowsAffected int64

	// single db
	db                SQLCommon
	blockGlobalUpdate bool
	logMode           logModeValue
	logger            logger
	search            *search
	values            sync.Map

	// global db
	parent        *DB
	callbacks     *Callback
	dialect       Dialect
	singularTable bool

	// function to be used to override the creating of a new timestamp
	nowFuncOverride func() time.Time
}

type UserStore struct {
	db *gorm.DB
}

type B struct {
	common
	importPath       string // import path of the package containing the benchmark
	context          *benchContext
	N                int
	previousN        int           // number of iterations in the previous run
	previousDuration time.Duration // total duration of the previous run
	benchFunc        func(b *B)
	benchTime        durationOrCountFlag
	bytes            int64
	missingBytes     bool // one of the subbenchmarks does not have bytes set.
	timerOn          bool
	showAllocResult  bool
	result           BenchmarkResult
	parallelism      int // RunParallel creates parallelism*GOMAXPROCS goroutines
	// The initial states of memStats.Mallocs and memStats.TotalAlloc.
	startAllocs uint64
	startBytes  uint64
	// The net total of this test after being run.
	netAllocs uint64
	netBytes  uint64
	// Extra metrics collected by ReportMetric.
	extra map[string]float64
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type B struct {
	common
	importPath       string
	context          *benchContext
	N                int
	previousN        int
	previousDuration time.Duration
	benchFunc        func(b *B)
	benchTime        durationOrCountFlag
	bytes            int64
	missingBytes     bool
	timerOn          bool
	showAllocResult  bool
	result           BenchmarkResult
	parallelism      int

	startAllocs uint64
	startBytes  uint64

	netAllocs uint64
	netBytes  uint64

	extra map[string]float64
}B struct {
	common
	importPath       string
	context          *benchContext
	N                int
	previousN        int
	previousDuration time.Duration
	benchFunc        func(b *B)
	benchTime        durationOrCountFlag
	bytes            int64
	missingBytes     bool
	timerOn          bool
	showAllocResult  bool
	result           BenchmarkResult
	parallelism      int

	startAllocs uint64
	startBytes  uint64

	netAllocs uint64
	netBytes  uint64

	extra map[string]float64
}
type T struct {
	common
	isEnvSet bool
	context  *testContext
}T struct {
	common
	isEnvSet bool
	context  *testContext
}
type ArticleStore struct {
	db *gorm.DB
}

type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)

	// PanicMsg holds msg to be used to mock panic on the function call
	//  if the PanicMsg is set to a non nil string the function call will panic
	// irrespective of other settings
	PanicMsg *string

	// Calls which must be satisfied before this call can be
	requires []*Call
}

type Call struct {
	Parent *Mock

	Method string

	Arguments Arguments

	ReturnArguments Arguments

	callerInfo []string

	Repeatability int

	totalCalls int

	optional bool

	WaitFor <-chan time.Time

	waitTime time.Duration

	RunFn func(Arguments)

	PanicMsg *string

	requires []*Call
}
type MockDB struct {
	mock.Mock
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

type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}

type Article struct {
	gorm.Model
	Title          string `gorm:"not null"`
	Description    string `gorm:"not null"`
	Body           string `gorm:"not null"`
	Tags           []Tag  `gorm:"many2many:article_tags"`
	Author         User   `gorm:"foreignkey:UserID"`
	UserID         uint   `gorm:"not null"`
	FavoritesCount int32  `gorm:"not null;default=0"`
	FavoritedUsers []User `gorm:"many2many:favorite_articles"`
	Comments       []Comment
}
type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}

type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}
type Association struct {
	Error  error
	scope  *Scope
	column string
	field  *Field
}
type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}

type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}

type ExpectedQuery struct {
	queryBasedExpectation
	rows             driver.Rows
	delay            time.Duration
	rowsMustBeClosed bool
	rowsWereClosed   bool
}
type Rows struct {
	converter driver.ValueConverter
	cols      []string
	def       []*Column
	rows      [][]driver.Value
	pos       int
	nextErr   map[int]error
	closeErr  error
}
/*
ROOST_METHOD_HASH=NewUserStore_6a331dd890
ROOST_METHOD_SIG_HASH=NewUserStore_4f0c2dfca9


 */
func BenchmarkNewUserStore(b *testing.B) {
	mockDB := &gorm.DB{}
	for i := 0; i < b.N; i++ {
		NewUserStore(mockDB)
	}
}

func TestNewUserStore(t *testing.T) {
	tests := []struct {
		name string
		db   *gorm.DB
		want *UserStore
	}{
		{
			name: "Create UserStore with valid gorm.DB instance",
			db:   &gorm.DB{},
			want: &UserStore{db: &gorm.DB{}},
		},
		{
			name: "Create UserStore with nil gorm.DB instance",
			db:   nil,
			want: &UserStore{db: nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserStore(tt.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserStoreConcurrency(t *testing.T) {
	mockDB := &gorm.DB{}
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			userStore := NewUserStore(mockDB)
			if userStore == nil || userStore.db != mockDB {
				t.Errorf("NewUserStore() failed in concurrent execution")
			}
		}()
	}

	wg.Wait()
}

func TestNewUserStoreDBFieldSet(t *testing.T) {
	mockDB := &gorm.DB{
		Value: "unique_identifier",
	}
	userStore := NewUserStore(mockDB)

	if userStore.db != mockDB {
		t.Errorf("NewUserStore() db field = %v, want %v", userStore.db, mockDB)
	}
}

func TestNewUserStoreMultipleInstances(t *testing.T) {
	db1 := &gorm.DB{Value: "db1"}
	db2 := &gorm.DB{Value: "db2"}

	userStore1 := NewUserStore(db1)
	userStore2 := NewUserStore(db2)

	if userStore1.db == userStore2.db {
		t.Errorf("NewUserStore() created non-unique instances")
	}
}


/*
ROOST_METHOD_HASH=Create_889fc0fc45
ROOST_METHOD_SIG_HASH=Create_4c48ec3920


 */
func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		dbSetup func(*gorm.DB)
		wantErr bool
	}{
		{
			name: "Successfully Create a New User",
			user: &model.User{
				Username: "newuser",
				Email:    "newuser@example.com",
				Password: "password123",
				Bio:      "New user bio",
				Image:    "https://example.com/newuser.jpg",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: false,
		},
		{
			name: "Attempt to Create a User with a Duplicate Username",
			user: &model.User{
				Username: "existinguser",
				Email:    "newuser@example.com",
				Password: "password123",
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "existinguser", Email: "existing@example.com"})
			},
			wantErr: true,
		},
		{
			name: "Attempt to Create a User with a Duplicate Email",
			user: &model.User{
				Username: "newuser",
				Email:    "existing@example.com",
				Password: "password123",
			},
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "existinguser", Email: "existing@example.com"})
			},
			wantErr: true,
		},
		{
			name: "Create a User with Minimum Required Fields",
			user: &model.User{
				Username: "minuser",
				Email:    "minuser@example.com",
				Password: "password123",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: false,
		},
		{
			name: "Attempt to Create a User with Invalid Email Format",
			user: &model.User{
				Username: "invaliduser",
				Email:    "notanemail",
				Password: "password123",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
		},
		{
			name: "Create a User with Maximum Length Values",
			user: &model.User{
				Username: "maxuser" + string(make([]byte, 250)),
				Email:    "maxuser@example.com",
				Password: string(make([]byte, 255)),
				Bio:      string(make([]byte, 1000)),
				Image:    "https://example.com/" + string(make([]byte, 235)),
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: false,
		},
		{
			name: "Attempt to Create a User with Empty Required Fields",
			user: &model.User{
				Username: "",
				Email:    "emptyuser@example.com",
				Password: "password123",
			},
			dbSetup: func(db *gorm.DB) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err)
			defer db.Close()

			err = db.AutoMigrate(&model.User{}).Error
			assert.NoError(t, err)

			tt.dbSetup(db)

			store := &UserStore{db: db}

			err = store.Create(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				var createdUser model.User
				result := db.Where("username = ?", tt.user.Username).First(&createdUser)
				assert.NoError(t, result.Error)
				assert.Equal(t, tt.user.Email, createdUser.Email)
			}
		})
	}
}

func (m *mockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetByID_bbf946112e
ROOST_METHOD_SIG_HASH=GetByID_728dd55ed1


 */
func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(out, where)
	return args.Get(0).(*gorm.DB)
}

func TestGetByID(t *testing.T) {
	tests := []struct {
		name      string
		id        uint
		mockSetup func(*MockDB)
		want      *model.User
		wantErr   error
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			s := &UserStore{db: mockDB}

			got, err := s.GetByID(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)

			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetByEmail_3574af40e5
ROOST_METHOD_SIG_HASH=GetByEmail_5731b833c1


 */
func (m *MockDB) First(out interface{}, where ...interface{}) *gorm.DB {
	called := m.Called(out, where)
	return called.Get(0).(*gorm.DB)
}

func TestGetByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockSetup     func(*MockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:  "Successfully retrieve a user by email",
			email: "user@example.com",
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Where", "email = ?", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.AnythingOfType("*model.User"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = model.User{
						Model:    gorm.Model{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						Username: "testuser",
						Email:    "user@example.com",
						Password: "hashedpassword",
						Bio:      "Test bio",
						Image:    "https://example.com/image.jpg",
					}
				}).Return(&gorm.DB{Error: nil})
			},
			expectedUser: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "testuser",
				Email:    "user@example.com",
				Password: "hashedpassword",
				Bio:      "Test bio",
				Image:    "https://example.com/image.jpg",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			userStore := &UserStore{db: mockDB}

			user, err := userStore.GetByEmail(tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedUser, user)

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	called := m.Called(query, args)
	return called.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=GetByUsername_f11f114df2
ROOST_METHOD_SIG_HASH=GetByUsername_954d096e24


 */
func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		mockSetup     func(*MockDB)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:     "Successfully retrieve a user by username",
			username: "testuser",
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Where", "username = ?", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.AnythingOfType("*model.User"), mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*model.User)
					*arg = model.User{Username: "testuser", Email: "test@example.com"}
				}).Return(mockDB)
			},
			expectedUser:  &model.User{Username: "testuser", Email: "test@example.com"},
			expectedError: nil,
		},
		{
			name:     "Attempt to retrieve a non-existent user",
			username: "nonexistent",
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Where", "username = ?", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.AnythingOfType("*model.User"), mock.Anything).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:     "Handle database connection error",
			username: "testuser",
			mockSetup: func(mockDB *MockDB) {
				mockDB.On("Where", "username = ?", mock.Anything).Return(mockDB)
				mockDB.On("First", mock.AnythingOfType("*model.User"), mock.Anything).Return(&gorm.DB{Error: errors.New("database connection error")})
			},
			expectedUser:  nil,
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.mockSetup(mockDB)

			userStore := &UserStore{db: mockDB}

			user, err := userStore.GetByUsername(tt.username)

			if tt.expectedUser == nil {
				assert.Nil(t, user)
			} else {
				assert.Equal(t, tt.expectedUser, user)
			}

			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError.Error())
			}

			mockDB.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=Update_68f27dd78a
ROOST_METHOD_SIG_HASH=Update_87150d6435


 */
func TestConcurrentUpdates(t *testing.T) {

}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*gorm.DB)
		input   *model.User
		wantErr bool
		errMsg  string
	}{
		{
			name: "Successfully Update User Information",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "olduser", Email: "old@example.com", Bio: "Old bio"})
			},
			input: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "olduser",
				Email:    "old@example.com",
				Bio:      "Updated bio",
			},
			wantErr: false,
		},
		{
			name:  "Update Non-Existent User",
			setup: func(db *gorm.DB) {},
			input: &model.User{
				Model: gorm.Model{ID: 999},
			},
			wantErr: true,
			errMsg:  "record not found",
		},
		{
			name: "Update User with Duplicate Username",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "user1", Email: "user1@example.com"})
				db.Create(&model.User{Username: "user2", Email: "user2@example.com"})
			},
			input: &model.User{
				Model:    gorm.Model{ID: 2},
				Username: "user1",
				Email:    "user2@example.com",
			},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "Update User with Empty Fields",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "user", Email: "user@example.com"})
			},
			input: &model.User{
				Model:    gorm.Model{ID: 1},
				Username: "",
				Email:    "",
			},
			wantErr: true,
			errMsg:  "not null constraint",
		},
		{
			name: "Update User's Relationships",
			setup: func(db *gorm.DB) {
				db.Create(&model.User{Username: "user", Email: "user@example.com"})
				db.Create(&model.User{Username: "follow", Email: "follow@example.com"})
				db.Create(&model.Article{Title: "Article 1"})
			},
			input: &model.User{
				Model:            gorm.Model{ID: 1},
				Username:         "user",
				Email:            "user@example.com",
				Follows:          []model.User{{Model: gorm.Model{ID: 2}}},
				FavoriteArticles: []model.Article{{Model: gorm.Model{ID: 1}}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			assert.NoError(t, err)
			defer db.Close()

			db.AutoMigrate(&model.User{}, &model.Article{})

			tt.setup(db)

			us := &UserStore{db: db}

			err = us.Update(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)

				var updatedUser model.User
				err = db.First(&updatedUser, tt.input.ID).Error
				assert.NoError(t, err)
				assert.Equal(t, tt.input.Username, updatedUser.Username)
				assert.Equal(t, tt.input.Email, updatedUser.Email)
				assert.Equal(t, tt.input.Bio, updatedUser.Bio)

				if len(tt.input.Follows) > 0 {
					var follows []model.User
					db.Model(&updatedUser).Association("Follows").Find(&follows)
					assert.Equal(t, len(tt.input.Follows), len(follows))
				}
				if len(tt.input.FavoriteArticles) > 0 {
					var favorites []model.Article
					db.Model(&updatedUser).Association("FavoriteArticles").Find(&favorites)
					assert.Equal(t, len(tt.input.FavoriteArticles), len(favorites))
				}
			}
		})
	}
}


/*
ROOST_METHOD_HASH=Follow_48fdf1257b
ROOST_METHOD_SIG_HASH=Follow_8217e61c06


 */
func (m *MockAssociation) Append(values ...interface{}) *gorm.Association {
	args := m.Called(values...)
	return args.Get(0).(*gorm.Association)
}

func (m *MockDB) Association(column string) *gorm.Association {
	args := m.Called(column)
	return args.Get(0).(*gorm.Association)
}

func (m *MockAssociation) Error() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestFollow(t *testing.T) {
	tests := []struct {
		name      string
		follower  *model.User
		followed  *model.User
		mockSetup func(*MockDB, *MockAssociation)
		wantErr   bool
	}{
		{
			name:     "Successfully follow a user",
			follower: &model.User{Username: "userA"},
			followed: &model.User{Username: "userB"},
			mockSetup: func(db *MockDB, assoc *MockAssociation) {
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "Follows").Return(assoc)
				assoc.On("Append", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Attempt to follow a user that is already being followed",
			follower: &model.User{Username: "userA", Follows: []model.User{{Username: "userB"}}},
			followed: &model.User{Username: "userB"},
			mockSetup: func(db *MockDB, assoc *MockAssociation) {
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "Follows").Return(assoc)
				assoc.On("Append", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Follow with a nil follower user",
			follower: nil,
			followed: &model.User{Username: "userB"},
			mockSetup: func(db *MockDB, assoc *MockAssociation) {

			},
			wantErr: true,
		},
		{
			name:     "Follow with a nil user to be followed",
			follower: &model.User{Username: "userA"},
			followed: nil,
			mockSetup: func(db *MockDB, assoc *MockAssociation) {

			},
			wantErr: true,
		},
		{
			name:     "Follow with database connection error",
			follower: &model.User{Username: "userA"},
			followed: &model.User{Username: "userB"},
			mockSetup: func(db *MockDB, assoc *MockAssociation) {
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "Follows").Return(assoc)
				assoc.On("Append", mock.Anything).Return(assoc)
				assoc.On("Error").Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:     "Follow oneself",
			follower: &model.User{Username: "userA"},
			followed: &model.User{Username: "userA"},
			mockSetup: func(db *MockDB, assoc *MockAssociation) {
				db.On("Model", mock.Anything).Return(db)
				db.On("Association", "Follows").Return(assoc)
				assoc.On("Append", mock.Anything).Return(assoc)
				assoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockAssoc := new(MockAssociation)
			tt.mockSetup(mockDB, mockAssoc)

			s := &UserStore{
				db: mockDB,
			}

			err := s.Follow(tt.follower, tt.followed)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockAssoc.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=IsFollowing_f53a5d9cef
ROOST_METHOD_SIG_HASH=IsFollowing_9eba5a0e9c


 */
func (m *MockDB) Count(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Table(name string) *gorm.DB {
	args := m.Called(name)
	return args.Get(0).(*gorm.DB)
}

func TestIsFollowing(t *testing.T) {
	tests := []struct {
		name     string
		userA    *model.User
		userB    *model.User
		mockDB   func(*MockDB)
		expected bool
		err      error
	}{
		{
			name:  "User A is following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func(db *MockDB) {
				db.On("Table", "follows").Return(db)
				db.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*int)
					*arg = 1
				}).Return(db)
			},
			expected: true,
			err:      nil,
		},
		{
			name:  "User A is not following User B",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func(db *MockDB) {
				db.On("Table", "follows").Return(db)
				db.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*int)
					*arg = 0
				}).Return(db)
			},
			expected: false,
			err:      nil,
		},
		{
			name:     "Null User A",
			userA:    nil,
			userB:    &model.User{Model: gorm.Model{ID: 2}},
			mockDB:   func(db *MockDB) {},
			expected: false,
			err:      nil,
		},
		{
			name:     "Null User B",
			userA:    &model.User{Model: gorm.Model{ID: 1}},
			userB:    nil,
			mockDB:   func(db *MockDB) {},
			expected: false,
			err:      nil,
		},
		{
			name:  "Database error",
			userA: &model.User{Model: gorm.Model{ID: 1}},
			userB: &model.User{Model: gorm.Model{ID: 2}},
			mockDB: func(db *MockDB) {
				db.On("Table", "follows").Return(db)
				db.On("Where", "from_user_id = ? AND to_user_id = ?", uint(1), uint(2)).Return(db)
				db.On("Count", mock.AnythingOfType("*int")).Return(&gorm.DB{Error: errors.New("database error")})
			},
			expected: false,
			err:      errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := new(MockDB)
			tt.mockDB(mockDB)

			store := &UserStore{db: mockDB}

			result, err := store.IsFollowing(tt.userA, tt.userB)

			assert.Equal(t, tt.expected, result)
			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}


/*
ROOST_METHOD_HASH=Unfollow_57959a8a53
ROOST_METHOD_SIG_HASH=Unfollow_8bd8e0bc55


 */
func (m *MockDB) Association(column string) *MockAssociation {
	args := m.Called(column)
	return args.Get(0).(*MockAssociation)
}

func (m *MockAssociation) Delete(values ...interface{}) *MockAssociation {
	args := m.Called(values...)
	return args.Get(0).(*MockAssociation)
}

func (m *MockDB) Model(value interface{}) *MockDB {
	args := m.Called(value)
	return args.Get(0).(*MockDB)
}

func TestUnfollow(t *testing.T) {
	tests := []struct {
		name      string
		follower  *model.User
		followee  *model.User
		mockSetup func(*MockDB, *MockAssociation)
		wantErr   bool
	}{
		{
			name:     "Successful Unfollow",
			follower: &model.User{Username: "userA"},
			followee: &model.User{Username: "userB"},
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "Follows").Return(mockAssoc)
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockAssoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "Unfollow User Not Previously Followed",
			follower: &model.User{Username: "userA"},
			followee: &model.User{Username: "userB"},
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "Follows").Return(mockAssoc)
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockAssoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Unfollow with Nil User (Follower)",
			follower:  nil,
			followee:  &model.User{Username: "userB"},
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {},
			wantErr:   true,
		},
		{
			name:      "Unfollow with Nil User (Followee)",
			follower:  &model.User{Username: "userA"},
			followee:  nil,
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {},
			wantErr:   true,
		},
		{
			name:     "Database Error During Unfollow",
			follower: &model.User{Username: "userA"},
			followee: &model.User{Username: "userB"},
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "Follows").Return(mockAssoc)
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockAssoc.On("Error").Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:     "Unfollow Self",
			follower: &model.User{Username: "userA"},
			followee: &model.User{Username: "userA"},
			mockSetup: func(mockDB *MockDB, mockAssoc *MockAssociation) {
				mockDB.On("Model", mock.Anything).Return(mockDB)
				mockDB.On("Association", "Follows").Return(mockAssoc)
				mockAssoc.On("Delete", mock.Anything).Return(mockAssoc)
				mockAssoc.On("Error").Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockAssoc := new(MockAssociation)
			tt.mockSetup(mockDB, mockAssoc)

			s := &UserStore{
				db: mockDB,
			}

			err := s.Unfollow(tt.follower, tt.followee)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockAssoc.AssertExpectations(t)
		})
	}
}


/*
ROOST_METHOD_HASH=GetFollowingUserIDs_ba3670aa2c
ROOST_METHOD_SIG_HASH=GetFollowingUserIDs_55ccc2afd7


 */
func TestGetFollowingUserIDs(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		user          *model.User
		expectedIDs   []uint
		expectedError error
	}{
		{
			name: "Successful retrieval of following user IDs",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"}).
					AddRow(uint(2)).
					AddRow(uint(3)).
					AddRow(uint(4))
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{2, 3, 4},
			expectedError: nil,
		},
		{
			name: "User with no followers",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
		{
			name: "Database error handling",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(uint(1)).
					WillReturnError(errors.New("database error"))
			},
			user:          &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs:   []uint{},
			expectedError: errors.New("database error"),
		},
		{
			name: "Large number of followers",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				for i := 2; i <= 1001; i++ {
					rows.AddRow(uint(i))
				}
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(uint(1)).
					WillReturnRows(rows)
			},
			user: &model.User{Model: gorm.Model{ID: 1}},
			expectedIDs: func() []uint {
				ids := make([]uint, 1000)
				for i := 0; i < 1000; i++ {
					ids[i] = uint(i + 2)
				}
				return ids
			}(),
			expectedError: nil,
		},
		{
			name: "User not found in database",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"to_user_id"})
				mock.ExpectQuery("SELECT to_user_id FROM follows WHERE from_user_id = ?").
					WithArgs(uint(999)).
					WillReturnRows(rows)
			},
			user:          &model.User{Model: gorm.Model{ID: 999}},
			expectedIDs:   []uint{},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			gormDB, err := gorm.Open("mysql", db)
			assert.NoError(t, err)
			defer gormDB.Close()

			tt.setupMock(mock)

			store := &UserStore{db: gormDB}

			ids, err := store.GetFollowingUserIDs(tt.user)

			assert.Equal(t, tt.expectedIDs, ids)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

