package store

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"github.com/jinzhu/gorm"
	"github.com/raahii/golang-grpc-realworld-example/model"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}
func TestGetTags(t *testing.T) {
	tests := []struct {
		name    string
		dbSetup func(*gorm.DB)
		want    []model.Tag
		wantErr bool
	}{
		{
			name: "Successfully Retrieve All Tags",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
				db.Create(&model.Tag{Name: "tag3"})
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag2"},
				{Model: gorm.Model{ID: 3}, Name: "tag3"},
			},
			wantErr: false,
		},
		{
			name:    "Empty Tag List",
			dbSetup: func(db *gorm.DB) {},
			want:    []model.Tag{},
			wantErr: false,
		},
		{
			name: "Database Connection Error",
			dbSetup: func(db *gorm.DB) {
				db.AddError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Large Number of Tags",
			dbSetup: func(db *gorm.DB) {
				for i := 1; i <= 1000; i++ {
					db.Create(&model.Tag{Name: fmt.Sprintf("tag%d", i)})
				}
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Duplicate Tag Names",
			dbSetup: func(db *gorm.DB) {
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag1"})
				db.Create(&model.Tag{Name: "tag2"})
			},
			want: []model.Tag{
				{Model: gorm.Model{ID: 1}, Name: "tag1"},
				{Model: gorm.Model{ID: 2}, Name: "tag1"},
				{Model: gorm.Model{ID: 3}, Name: "tag2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open("sqlite3", ":memory:")
			if err != nil {
				t.Fatalf("failed to open database: %v", err)
			}
			defer db.Close()

			db.AutoMigrate(&model.Tag{})

			tt.dbSetup(db)

			s := &ArticleStore{db: db}

			got, err := s.GetTags()

			if (err != nil) != tt.wantErr {
				t.Errorf("ArticleStore.GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.name == "Large Number of Tags" {
				if len(got) != 1000 {
					t.Errorf("ArticleStore.GetTags() returned %d tags, want 1000", len(got))
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArticleStore.GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
