package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	field.Int("id").Positive().Unique()
	field.String("text").MaxLen(255).NotEmpty()
	field.Time("created_at").Default(time.Now)
	field.String("user_name").MaxLen(20).NotEmpty()
	return nil
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return nil
}
