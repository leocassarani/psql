package psql

import "fmt"

type aggregationType string

const (
	avg aggregationType = "AVG"
	max                 = "MAX"
	min                 = "MIN"
	sum                 = "SUM"
)

// Avg returns an Expression representing a call to the AVG aggregate
// function with the table column col as an argument.
func Avg(col tableColumn) aggregateFunc {
	return aggregateFunc{col, avg}
}

// Max returns an Expression representing a call to the MAX aggregate
// function with the table column col as an argument.
func Max(col tableColumn) aggregateFunc {
	return aggregateFunc{col, max}
}

// Min returns an Expression representing a call to the MIN aggregate
// function with the table column col as an argument.
func Min(col tableColumn) aggregateFunc {
	return aggregateFunc{col, min}
}

// Sum returns an Expression representing a call to the SUM aggregate
// function with the table column col as an argument.
func Sum(col tableColumn) aggregateFunc {
	return aggregateFunc{col, sum}
}

type aggregateFunc struct {
	column tableColumn
	fnType aggregationType
}

func (f aggregateFunc) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("%s(%s)", f.fnType, f.column.ToSQLExpr(p))
}

func (f aggregateFunc) Relations() []string {
	return f.column.Relations()
}
