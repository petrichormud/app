package request

import (
	"regexp"

	fiber "github.com/gofiber/fiber/v2"
)

const (
	CommentMinLength = 1
	CommentMaxLength = 500
)

type Comment struct {
	Text           string
	Author         string
	AvatarLink     string
	Replies        []Comment
	ID             int64
	CreatedAt      int64
	VID            int32
	ViewedByAuthor bool
	Current        bool
}

func (c *Comment) Bind() fiber.Map {
	return fiber.Map{
		"Current":        c.Current,
		"ID":             c.ID,
		"VID":            c.VID,
		"Text":           c.Text,
		"Author":         c.Author,
		"AvatarLink":     c.AvatarLink,
		"CreatedAt":      c.CreatedAt,
		"ViewedByAuthor": c.ViewedByAuthor,
	}
}

func SanitizeComment(c string) string {
	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return re.ReplaceAllString(c, "")
}

func IsCommentValid(c string) bool {
	if len(c) < CommentMinLength {
		return false
	}

	if len(c) > CommentMaxLength {
		return false
	}

	re := regexp.MustCompile("[^a-zA-Z, \"'\\-\\.?!()\\r\\n]+")
	return !re.MatchString(c)
}
