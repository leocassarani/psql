package psql

import (
	"fmt"
	"strconv"
	"strings"
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
	for i, expr := range s.exprs {
		args[i] = expr.ToSQLExpr()
	}

	return fmt.Sprintf("SELECT %s", strings.Join(args, ", "))
}

// An Expression can be used in the SELECT list of an SQL query.
type Expression interface {
	ToSQLExpr() string
}

// IntLiteral wraps an integer n in an Expression that can be used to
// build an SQL query.
func IntLiteral(n int) intLiteral {
	return intLiteral(n)
}

type intLiteral int

// ToSQLExpr returns the SQL string representation of the integer literal.
func (i intLiteral) ToSQLExpr() string {
	return strconv.Itoa(int(i))
}
