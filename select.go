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
	return SelectQuery{exprs}
}

// A SelectQuery represents a SELECT query with all its clauses.
type SelectQuery struct {
	exprs []Expression
}

// ToSQL returns a string containing the full SQL query version of the
// SelectQuery. If the query is empty, an empty string is returned.
func (s SelectQuery) ToSQL() string {
	if len(s.exprs) == 0 {
		return ""
	}

	args := make([]string, len(s.exprs))
	rels := []string{}

	for i, expr := range s.exprs {
		args[i] = expr.ToSQLExpr()
		rels = append(rels, expr.Relations()...)
	}

	selectList := fmt.Sprintf("SELECT %s", strings.Join(args, ", "))
	if len(rels) == 0 {
		return selectList
	}

	fromClause := fmt.Sprintf("FROM %s", strings.Join(rels, ","))
	return fmt.Sprintf("%s %s", selectList, fromClause)
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
