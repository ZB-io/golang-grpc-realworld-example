package model

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"

	"your_project/pb" // Import your protobuf package
)

// Define the ISO8601 constant
const ISO8601 = "2006-01-02T15:04:05Z07:00"

// Define the Comment struct
type Comment struct {
	gorm.Model
	Body      string
	UserID    uint
	ArticleID uint
}

func (c *Comment) ProtoComment() *pb.Comment {
	return &pb.Comment{
		Id:        fmt.Sprintf("%d", c.ID),
		Body:      c.Body,
		CreatedAt: c.CreatedAt.Format(ISO8601),
		UpdatedAt: c.UpdatedAt.Format(ISO8601),
	}
}

func (c *Comment) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Body, validation.Required, validation.Length(1, 1000)),
	)
}

func TestProtoComment(t *testing.T) {
	// ... (rest of the TestProtoComment function remains the same)
}

func TestValidate(t *testing.T) {
	// ... (rest of the TestValidate function remains the same)
}
