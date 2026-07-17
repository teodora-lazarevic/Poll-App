package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Poll holds the schema definition for the Poll entity.
type Poll struct {
	ent.Schema
}

// Fields of the Poll.
func (Poll) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("description").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the Poll.
func (Poll) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("creator", User.Type).Ref("polls").Unique().Required(),                 // Each poll has one creator
		edge.To("options", PollOption.Type).Annotations(entsql.OnDelete(entsql.Cascade)), // A poll has many options, delete options when poll is deleted
		edge.To("votes", Vote.Type).Annotations(entsql.OnDelete(entsql.Cascade)),         // A poll has many votes, delete votes when poll is deleted
	}
}
