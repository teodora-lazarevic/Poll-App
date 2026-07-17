package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
)

// Vote holds the schema definition for the Vote entity.
type Vote struct {
	ent.Schema
}

// Fields of the Vote.
func (Vote) Fields() []ent.Field {
	return []ent.Field{}
}

// Edges of the Vote.
func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("votes").Unique().Required(),         // Each vote belongs to one user
		edge.From("option", PollOption.Type).Ref("votes").Unique().Required(), // Each vote belongs to one option
		edge.From("poll", Poll.Type).Ref("votes").Unique().Required(),         // Each vote belongs to one poll
	}
}

func (Vote) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("user", "poll", "option").Unique(),
	}
}
