package psql

import (
	"fmt"
	"strconv"
)

// IntLiteral wraps an integer n in an Expression that can be used to
// build an SQL query.
func IntLiteral(n int) intLiteral {
	return intLiteral(n)
}

type intLiteral int

func (i intLiteral) ToSQLExpr(*Params) string {
	return strconv.Itoa(int(i))
}

func (i intLiteral) Relations() []string {
	return nil
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

func (o binaryOp) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("(%s %s %s)", o.a.ToSQLExpr(p), o.opType, o.b.ToSQLExpr(p))
}

func (o binaryOp) Relations() []string {
	var rels []string
	rels = append(rels, o.a.Relations()...)
	rels = append(rels, o.b.Relations()...)
	return rels
}
