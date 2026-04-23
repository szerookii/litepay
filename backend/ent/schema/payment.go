package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type Payment struct {
	ent.Schema
}

func (Payment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("wallet_address"),
		field.String("sol_reference").Optional().Nillable(),
		field.Int("address_index").Default(0),
		field.Float("amount_crypto"),
		field.String("currency_crypto"),
		field.Float("amount_fiat"),
		field.String("currency_fiat"),
		field.Float("received_amount").Optional().Nillable(),
		field.String("transaction_hash").Optional().Nillable(),
		field.Enum("status").Values("PENDING", "CONFIRMING", "PAID", "EXPIRED", "REFUNDED", "CASHED_OUT").Default("PENDING"),
		field.Time("expires_at").Default(func() time.Time {
			return time.Now().Add(time.Hour)
		}),
	}
}

func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("payments").
			Field("user_id").
			Required().
			Unique(),
	}
}

func (Payment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("wallet_address"),
		index.Fields("user_id"),
		index.Fields("sol_reference"),
	}
}
