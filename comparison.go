package psql

import "fmt"

// Eq returns an Expression representing the equality comparison between a and b.
func Eq(a, b Expression) comparison {
	return comparison{a, b, eq}
}

// NotEq returns an Expression representing the inequality comparison between a and b.
func NotEq(a, b Expression) comparison {
	return comparison{a, b, neq}
}

// LessThan returns an Expression representing the less-than comparison between a and b.
func LessThan(a, b Expression) comparison {
	return comparison{a, b, lt}
}

// LessThanOrEq returns an Expression representing the less-than-or-equal-to
// comparison between a and b.
func LessThanOrEq(a, b Expression) comparison {
	return comparison{a, b, lte}
}

// GreaterThan returns an Expression representing the greater-than comparison
// between a and b.
func GreaterThan(a, b Expression) comparison {
	return comparison{a, b, gt}
}

// GreaterThanOrEq returns an Expression representing the greater-than-or-equal-to
// comparison between a and b.
func GreaterThanOrEq(a, b Expression) comparison {
	return comparison{a, b, gte}
}

type comparison struct {
	a, b     Expression
	compType comparisonType
}

func (c comparison) ToSQLBoolean(p *Params) string {
	return c.ToSQLExpr(p)
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

type comparisonType int

const (
	eq comparisonType = iota
	neq
	lt
	lte
	gt
	gte
)

func (c comparisonType) String() string {
	switch c {
	case eq:
		return "="
	case neq:
		return "<>"
	case lt:
		return "<"
	case lte:
		return "<="
	case gt:
		return ">"
	case gte:
		return ">="
	default:
		panic("unknown comparisonType")
	}
}

// IsNull returns an Expression comparing expr and NULL for equality.
func IsNull(expr Expression) nullCheck {
	return nullCheck{expr, false}
}

// IsNotNull returns an Expression comparing expr and NULL for inequality.
func IsNotNull(expr Expression) nullCheck {
	return nullCheck{expr, true}
}

type nullCheck struct {
	expr    Expression
	negated bool
}

func (c nullCheck) ToSQLBoolean(p *Params) string {
	return c.ToSQLExpr(p)
}

func (c nullCheck) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("%s %s", c.expr.ToSQLExpr(p), c.operator())
}

func (c nullCheck) operator() string {
	if c.negated {
		return "IS NOT NULL"
	}
	return "IS NULL"
}

func (c nullCheck) Relations() []string {
	return c.expr.Relations()
}
