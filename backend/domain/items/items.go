package items

import (
	"time"
	"unicode/utf8"

	"github.com/oklog/ulid"
)

type Items struct {
	id        string
	title     string
	content   string
	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func newItem(id string, title string, content string) *Items {
	// ulid validation

	// title validation
	if utf8.RuneCountInString(title) < MinTitleLength || utf8.RuneCountInString(title) > MaxTitleLength {
		return nil
	}

	// content validation
	if utf8.RuneCountInString(content) < MinContentLength || utf8.RuneCountInString(content) > MaxContentLength {
		return nil
	}

	return &Items{
		id:        id,
		title:     title,
		content:   content,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		deletedAt: nil,
	}
}

func NewItem(title string, content string) *Items {
	id := ulid.MustNew(ulid.Now(), nil)
	return newItem(id.String(), title, content)
}

const (
	MinTitleLength   = 1
	MaxTitleLength   = 100
	MinContentLength = 1
	MaxContentLength = 10000
)
