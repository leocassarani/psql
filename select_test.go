package psql

import "testing"

func TestSelectQuery(t *testing.T) {
	cases := []struct {
		query SelectQuery
		sql   string
	}{}

	for i, tc := range cases {
		got := tc.query.ToSQL()
		if got != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, got)
		}
	}
}
