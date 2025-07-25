package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Url struct {
	ent.Schema
}

func (Url) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Immutable(),
		field.Bool("is_used").Default(false),
		field.Bool("is_image").Default(false),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("expires_at").Default(func() time.Time { return time.Now().Add(time.Hour * 3) }).Immutable(),
	}
}

func (Url) Edges() []ent.Edge {
	return nil
}
