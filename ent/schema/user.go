package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username"),
		field.String("password"),
		field.String("avatarFilename"),
		field.Time("createdAt").
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
