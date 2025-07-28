package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Video struct {
	ent.Schema
}

func (Video) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.String("path").Immutable(),
	}
}

func (Video) Edges() []ent.Edge {
	return nil
}
