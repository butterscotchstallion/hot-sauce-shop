package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Inventory holds the schema definition for the Inventory entity.
type Inventory struct {
	ent.Schema
}

// Fields of the Inventory.
func (Inventory) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
		field.String("shortDescription"),
		field.String("slug"),
		field.Float32("price"),
		field.Time("createdAt").
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now).
			Optional().
			Nillable(),
	}
}

// Edges of the Inventory.
func (Inventory) Edges() []ent.Edge {
	return nil
}
