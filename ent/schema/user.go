package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MinLen(3).MaxLen(32).NotEmpty(),
		field.Int("age").Min(18).Max(130),
	}
}

func (User) Edges() []ent.Edge {
	return nil
}
