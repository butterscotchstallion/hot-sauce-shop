package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Inventory holds the schema definition for the Inventory entity.
type Inventory struct {
	ent.Schema
}

func (Inventory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.String("shortDescription"),
		field.String("slug"),
		field.Float32("price"),
		field.Int("spiceRating").
			Default(3),
		field.Time("createdAt").
			Default(time.Now).
			Optional().
			Nillable(),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

func (Inventory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tags", Tag.Type),
		edge.From("cartItems", CartItems.Type).Ref("inventory"),
	}
}
