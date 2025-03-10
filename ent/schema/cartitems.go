package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CartItems holds the schema definition for the CartItems entity.
type CartItems struct {
	ent.Schema
}

// Fields of the CartItems.
func (CartItems) Fields() []ent.Field {
	return []ent.Field{
		field.Int8("quantity").Default(1),
		field.Time("createdAt").
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the CartItems.
func (CartItems) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("inventory", Inventory.Type),
	}
}
