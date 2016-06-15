package psql

import "fmt"

// Now returns an Expression representing a call to the 0-arity date/time function now().
func Now() fnCall {
	return fnCall{"now"}
}

type fnCall struct {
	name string
}

func (f fnCall) ToSQLExpr(*Params) string {
	return fmt.Sprintf("%s()", f.name)
}

func (f fnCall) Relations() []string {
	return nil
}
