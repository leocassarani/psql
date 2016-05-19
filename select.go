package psql

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// Select creates a new SelectQuery using the expressions exprs to populate
// the SELECT list, that is the portion of the query between the key words
// "SELECT" and "FROM".
func Select(exprs ...Expression) SelectQuery {
	return SelectQuery{
		sel: selectClause{exprs},
	}
}

// A SelectQuery represents a SELECT query with all its clauses.
type SelectQuery struct {
	sel     selectClause
	orderBy orderByClause
}

// OrderBy returns a copy of the SelectQuery s with an additional ORDER BY
// clause containing the args provided. If an ORDER BY clause was already
// present, this operation will overwrite it.
func (s SelectQuery) OrderBy(exprs ...OrderExpression) SelectQuery {
	s.orderBy = orderByClause{exprs}
	return s
}

// ToSQL returns a string containing the full SQL query version of the
// SelectQuery. If the query is empty, an empty string is returned.
func (s SelectQuery) ToSQL() string {
	var parts []string

	for _, clause := range s.clauses() {
		if part := clause.ToSQLClause(); part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, " ")
}

func (s SelectQuery) clauses() []Clause {
	return []Clause{
		s.sel,
		s.from(),
		s.orderBy,
	}
}

func (s SelectQuery) from() Clause {
	rels := make([]string, 0)
	set := make(map[string]struct{})

	for _, rel := range s.relations() {
		// Have we seen this relation before?
		if _, ok := set[rel]; !ok {
			set[rel] = struct{}{}
			rels = append(rels, rel)
		}
	}

	return fromClause{rels}
}

func (s SelectQuery) relations() []string {
	var rels []string
	rels = append(rels, s.sel.Relations()...)
	rels = append(rels, s.orderBy.Relations()...)
	return rels
}

// Clause is the interface that represents the individual components of
// an SQL query.
//
// ToSQLClause converts the clause to a string representation. If an empty
// string is returned, the clause will be omitted from the final SQL query.
type Clause interface {
	ToSQLClause() string
}

type selectClause struct {
	exprs []Expression
}

func (s selectClause) ToSQLClause() string {
	if len(s.exprs) == 0 {
		return ""
	}

	args := make([]string, len(s.exprs))
	for i, expr := range s.exprs {
		args[i] = expr.ToSQLExpr()
	}

	return fmt.Sprintf("SELECT %s", strings.Join(args, ", "))
}

func (s selectClause) Relations() []string {
	var rels []string
	for _, expr := range s.exprs {
		rels = append(rels, expr.Relations()...)
	}
	return rels
}

type fromClause struct {
	rels []string
}

func (f fromClause) ToSQLClause() string {
	if len(f.rels) == 0 {
		return ""
	}

	return fmt.Sprintf("FROM %s", strings.Join(f.rels, ", "))
}

type orderByClause struct {
	exprs []OrderExpression
}

func (o orderByClause) ToSQLClause() string {
	if len(o.exprs) == 0 {
		return ""
	}

	parts := make([]string, len(o.exprs))
	for i, expr := range o.exprs {
		parts[i] = expr.ToSQLOrder()
	}

	return fmt.Sprintf("ORDER BY %s", strings.Join(parts, ", "))
}

func (o orderByClause) Relations() []string {
	var rels []string
	for _, expr := range o.exprs {
		rels = append(rels, expr.Relations()...)
	}
	return rels
}

// Ascending returns a new OrderExpression specifying that the results
// of the query must be ordered by the given Expression in ascending order.
func Ascending(expr Expression) OrderExpression {
	return OrderExpression{expr, asc}
}

// Descending returns a new OrderExpression specifying that the results
// of the query must be ordered by the given Expression in descending order.
func Descending(expr Expression) OrderExpression {
	return OrderExpression{expr, desc}
}

type orderDirection string

const (
	asc  orderDirection = "ASC"
	desc                = "DESC"
)

// An OrderExpression is each individual component of a SELECT query's
// ORDER BY clause, specifying that the results of the query must be
// sorted by a given SQL expression.
type OrderExpression struct {
	expr      Expression
	direction orderDirection
}

func (o OrderExpression) ToSQLOrder() string {
	return fmt.Sprintf("%s %s", o.expr.ToSQLExpr(), o.direction)
}

func (o OrderExpression) Relations() []string {
	return o.expr.Relations()
}

// Expression is the interface that represents any SQL expression that
// can be used in the SELECT list of an SQL query.
//
// ToSQLExpr converts the expression into a snippet of SQL that may be
// safely embedded in a query. If the expression needs to be quoted or
// otherwise escaped, ToSQLExpr must return the quoted version.
//
// Relations returns a slice of strings corresponding to the (quoted)
// names of all relations used by the Expression.
type Expression interface {
	ToSQLExpr() string
	Relations() []string
}

// AllColumns returns an Expression representing all columns in the table.
func AllColumns(table string) allColumns {
	return allColumns{table}
}

type allColumns struct {
	table string
}

func (ac allColumns) ToSQLExpr() string {
	return fmt.Sprintf("%s.*", pq.QuoteIdentifier(ac.table))
}

func (ac allColumns) Relations() []string {
	return []string{
		pq.QuoteIdentifier(ac.table),
	}
}

// TableColumn returns an Expression representing the column col of the
// database table with the given name.
func TableColumn(table, col string) tableColumn {
	return tableColumn{table, col}
}

type tableColumn struct {
	table, column string
}

func (tc tableColumn) ToSQLExpr() string {
	return pq.QuoteIdentifier(tc.column)
}

func (tc tableColumn) Relations() []string {
	return []string{
		pq.QuoteIdentifier(tc.table),
	}
}
