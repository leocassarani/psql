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

func (i intLiteral) ToSQLExpr() string {
	return strconv.Itoa(int(i))
}

// Plus returns an Expression representing the addition of Expressions a and b.
func Plus(a, b Expression) binaryOp {
	return binaryOp{a, b, plus}
}

// Minus returns an Expression representing the subtraction of Expression b from a.
func Minus(a, b Expression) binaryOp {
	return binaryOp{a, b, minus}
}

// Times returns an Expression representing the multiplication of Expression a and b.
func Times(a, b Expression) binaryOp {
	return binaryOp{a, b, times}
}

// Divide returns an Expression representing the division of Expression a and b.
func Divide(a, b Expression) binaryOp {
	return binaryOp{a, b, divide}
}

// Modulo returns an Expression representing the modulo of Expression a and b.
func Modulo(a, b Expression) binaryOp {
	return binaryOp{a, b, modulo}
}

// Modulo returns an Expression representing the exponentiation of Expression a and b.
func Pow(a, b Expression) binaryOp {
	return binaryOp{a, b, pow}
}

type binaryOpType string

const (
	plus   binaryOpType = "+"
	minus               = "-"
	times               = "*"
	divide              = "/"
	modulo              = "%"
	pow                 = "^"
)

type binaryOp struct {
	a, b   Expression
	opType binaryOpType
}

func (o binaryOp) ToSQLExpr() string {
	return fmt.Sprintf("(%s %s %s)", o.a.ToSQLExpr(), o.opType, o.b.ToSQLExpr())
}
