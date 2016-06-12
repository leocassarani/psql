package psql

import "fmt"

// Eq returns an Expression representing the equality comparison between a and b.
func Eq(a, b Expression) comparison {
	return comparison{a, b, eq}
}

// Neq returns an Expression representing the inequality comparison between a and b.
func Neq(a, b Expression) comparison {
	return comparison{a, b, neq}
}

// Lt returns an Expression representing the less-than comparison between a and b.
func Lt(a, b Expression) comparison {
	return comparison{a, b, lt}
}

// Lte returns an Expression representing the less-than-or-equal-to comparison between a and b.
func Lte(a, b Expression) comparison {
	return comparison{a, b, lte}
}

// Gt returns an Expression representing the greater-than comparison between a and b.
func Gt(a, b Expression) comparison {
	return comparison{a, b, gt}
}

// Gte returns an Expression representing the greater-than-or-equal-to comparison between a and b.
func Gte(a, b Expression) comparison {
	return comparison{a, b, gte}
}

type comparisonType string

const (
	eq  comparisonType = "="
	neq                = "<>"
	lt                 = "<"
	lte                = "<="
	gt                 = ">"
	gte                = ">="
)

type comparison struct {
	a, b     Expression
	compType comparisonType
}

func (c comparison) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("(%s %s %s)", c.a.ToSQLExpr(p), c.compType, c.b.ToSQLExpr(p))
}

func (c comparison) Relations() []string {
	var rels []string
	rels = append(rels, c.a.Relations()...)
	rels = append(rels, c.b.Relations()...)
	return rels
}
