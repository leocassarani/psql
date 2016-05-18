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
		{
			Select(
				Plus(IntLiteral(9), IntLiteral(33)),
				Minus(IntLiteral(123), IntLiteral(81)),
				Times(IntLiteral(14), IntLiteral(3)),
				Divide(IntLiteral(714), IntLiteral(17)),
			),
			"SELECT 9 + 33, 123 - 81, 14 * 3, 714 / 17",
		},
	}

	for i, tc := range cases {
		got := tc.query.ToSQL()
		if got != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, got)
		}
	}
}
