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
		sel:    selectClause{exprs},
		params: newParams(),
	}
}

// A SelectQuery represents a SELECT query with all its clauses.
type SelectQuery struct {
	sel     selectClause
	orderBy orderByClause
	where   whereClause
	groupBy groupByClause

	params *Params
}

// OrderBy returns a copy of the SelectQuery s with an additional ORDER BY
// clause containing the order expressions provided. If an ORDER BY clause
// was already present, this method will overwrite it.
func (s SelectQuery) OrderBy(exprs ...OrderExpression) SelectQuery {
	s.orderBy = orderByClause{exprs}
	return s
}

// GroupBy returns a copy of the SelectQuery s with an additional GROUP BY
// clause containing the Expressions provided. If a GROUP BY clause was
// already present, this method will overwrite it.
func (s SelectQuery) GroupBy(exprs ...Expression) SelectQuery {
	s.groupBy = groupByClause{exprs}
	return s
}

// Where returns a copy of the SelectQuery s with an additional WHERE clause
// containing the BooleanExpressions provided. If a WHERE clause was already
// present, this method will overwrite it.
func (s SelectQuery) Where(exprs ...BooleanExpression) SelectQuery {
	s.where = whereClause{exprs}
	return s
}

// ToSQL returns a string containing the full SQL query version of the
// SelectQuery. If the query is empty, an empty string is returned.
func (s SelectQuery) ToSQL() string {
	var parts []string

	// Since ToSQL() may be called multiple times, we must reset the params
	// so that we always start with a blank slate.
	s.params.Reset()

	for _, clause := range s.clauses() {
		if part := clause.ToSQLClause(s.params); part != "" {
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
		s.where,
		s.groupBy,
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
	rels = append(rels, s.where.Relations()...)
	rels = append(rels, s.groupBy.Relations()...)
	rels = append(rels, s.orderBy.Relations()...)
	return rels
}

// Bindings returns a slice of arguments that can be unpacked and passed
// into the Query and QueryRow methods of the database/sql package.
//
// Any variadic arguments passed into Bindings will be used to replace
// user-supplied parameters in the SELECT query, in the same order as
// they appear in the query.
func (s SelectQuery) Bindings(inputs ...interface{}) []interface{} {
	return s.params.Values(inputs)
}

// Clause is the interface that represents the individual components of
// an SQL query.
//
// ToSQLClause converts the clause to a string representation. If an empty
// string is returned, the clause will be omitted from the final SQL query.
type Clause interface {
	ToSQLClause(p *Params) string
}

type selectClause struct {
	exprs []Expression
}

func (s selectClause) ToSQLClause(p *Params) string {
	if len(s.exprs) == 0 {
		return ""
	}

	args := make([]string, len(s.exprs))
	for i, expr := range s.exprs {
		args[i] = expr.ToSQLExpr(p)
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

func (f fromClause) ToSQLClause(*Params) string {
	if len(f.rels) == 0 {
		return ""
	}

	return fmt.Sprintf("FROM %s", strings.Join(f.rels, ", "))
}

type whereClause struct {
	exprs []BooleanExpression
}

func (w whereClause) ToSQLClause(p *Params) string {
	if len(w.exprs) == 0 {
		return ""
	}

	conds := make([]string, len(w.exprs))
	for i, expr := range w.exprs {
		conds[i] = expr.ToSQLBoolean(p)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(conds, " AND "))
}

func (w whereClause) Relations() []string {
	var rels []string
	for _, expr := range w.exprs {
		rels = append(rels, expr.Relations()...)
	}
	return rels
}

// BooleanExpression is the interface representing an SQL expression that returns
// a boolean value.
//
// ToSQLBoolean converts the BooleanExpression to an SQL representation. In most
// cases, the implementation will be identical to ToSQLExpr().
type BooleanExpression interface {
	Expression
	ToSQLBoolean(*Params) string
}

type groupByClause struct {
	exprs []Expression
}

func (g groupByClause) ToSQLClause(p *Params) string {
	if len(g.exprs) == 0 {
		return ""
	}

	parts := make([]string, len(g.exprs))
	for i, expr := range g.exprs {
		parts[i] = expr.ToSQLExpr(p)
	}

	return fmt.Sprintf("GROUP BY %s", strings.Join(parts, ", "))
}

func (g groupByClause) Relations() []string {
	var rels []string
	for _, expr := range g.exprs {
		rels = append(rels, expr.Relations()...)
	}
	return rels
}

type orderByClause struct {
	exprs []OrderExpression
}

func (o orderByClause) ToSQLClause(p *Params) string {
	if len(o.exprs) == 0 {
		return ""
	}

	parts := make([]string, len(o.exprs))
	for i, expr := range o.exprs {
		parts[i] = expr.ToSQLOrder(p)
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

// An OrderExpression is each individual component of a SELECT query's
// ORDER BY clause, specifying that the results of the query must be
// sorted by a given SQL expression.
type OrderExpression struct {
	expr      Expression
	direction orderDirection
}

func (o OrderExpression) ToSQLOrder(p *Params) string {
	return fmt.Sprintf("%s %s", o.expr.ToSQLExpr(p), o.direction)
}

func (o OrderExpression) Relations() []string {
	return o.expr.Relations()
}

type orderDirection int

const (
	asc orderDirection = iota
	desc
)

func (o orderDirection) String() string {
	switch o {
	case asc:
		return "ASC"
	case desc:
		return "DESC"
	default:
		panic("unknown orderDirection")
	}
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
	ToSQLExpr(*Params) string
	Relations() []string
}

// AllColumns returns an Expression representing all columns in the table.
func AllColumns(table string) allColumns {
	return allColumns{table}
}

type allColumns struct {
	table string
}

func (ac allColumns) ToSQLExpr(*Params) string {
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

func (tc tableColumn) ToSQLExpr(*Params) string {
	return pq.QuoteIdentifier(tc.column)
}

func (tc tableColumn) Relations() []string {
	return []string{
		pq.QuoteIdentifier(tc.table),
	}
}

func newParams() *Params {
	return &Params{
		values: make(map[int]interface{}),
	}
}

type Params struct {
	counter int
	values  map[int]interface{}
}

func (p *Params) Add(value interface{}) string {
	marker := p.next()
	p.values[p.counter] = value
	return marker
}

func (p *Params) New() string {
	return p.next()
}

func (p *Params) next() string {
	p.counter++
	return fmt.Sprintf("$%d", p.counter)
}

func (p *Params) Values(inputs []interface{}) []interface{} {
	var values []interface{}

	for i := 1; i <= p.counter; i++ {
		if val, ok := p.values[i]; ok {
			values = append(values, val)
		} else if len(inputs) > 0 {
			values = append(values, inputs[0])
			inputs = inputs[1:]
		}
	}

	return values
}

func (p *Params) Reset() {
	if p.counter > 0 {
		p.counter = 0
		p.values = make(map[int]interface{})
	}
}

// StringLiteral returns a text literal that will be replaced with str
// when the query is executed. For security reasons, the contents of str
// are not directly interpolated into the query's SQL representation.
// Instead, a marker such as $1 is used, and the value of str is injected
// into the appropriate position when the Bindings() method is called.
func StringLiteral(str string) stringLiteral {
	return stringLiteral(str)
}

type stringLiteral string

func (s stringLiteral) ToSQLExpr(params *Params) string {
	marker := params.Add(string(s))
	return fmt.Sprintf("%s::%s", marker, textType)
}

func (stringLiteral) Relations() []string {
	return nil
}

type dataType int

const (
	textType dataType = iota
)

func (d dataType) String() string {
	switch d {
	case textType:
		return "text"
	default:
		panic("unknown dataType")
	}
}

// StringParam returns a free (unbound) parameter using the "text" type.
func StringParam() freeParam {
	return freeParam{textType}
}

type freeParam struct {
	dataType dataType
}

func (p freeParam) ToSQLExpr(params *Params) string {
	return fmt.Sprintf("%s::%s", params.New(), p.dataType)
}

func (p freeParam) Relations() []string {
	return nil
}
