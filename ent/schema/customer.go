package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Customer holds the schema definition for the Customer entity.
type Customer struct {
	ent.Schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).StorageKey("oid"),
		field.Time("created_at").Default(time.Now()),
		field.Time("updated_at"),
		field.Int("phone").Unique(),
		field.String("first_name").Default(""),
		field.String("last_name").Optional(),
		field.String("email").Optional(),
		field.Time("date_of_birth").Optional(),
		field.Bool("is_new").Default(true),
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return nil
}
