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
		selectClause{exprs},
		fromClause{exprs},
	}
}

// A SelectQuery represents a SELECT query with all its clauses.
type SelectQuery struct {
	selectClause selectClause
	fromClause   fromClause
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
		s.selectClause,
		s.fromClause,
	}
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

type fromClause struct {
	exprs []Expression
}

func (f fromClause) ToSQLClause() string {
	rels := make([]string, 0)
	set := make(map[string]struct{})

	for _, expr := range f.exprs {
		for _, rel := range expr.Relations() {
			// Have we seen this relation before?
			if _, ok := set[rel]; !ok {
				set[rel] = struct{}{}
				rels = append(rels, rel)
			}
		}
	}

	if len(rels) == 0 {
		return ""
	}

	return fmt.Sprintf("FROM %s", strings.Join(rels, ", "))
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
