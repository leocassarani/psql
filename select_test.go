package psql

import "testing"

func TestSelectQuery(t *testing.T) {
	cases := []struct {
		query SelectQuery
		sql   string
	}{
		{
			Select(),
			"",
		},
		{
			Select(
				IntLiteral(123),
			),
			"SELECT 123",
		},
		{
			Select(
				IntLiteral(123), IntLiteral(42),
			),
			"SELECT 123, 42",
		},
	}

	for i, tc := range cases {
		got := tc.query.ToSQL()
		if got != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, got)
		}
	}
}
