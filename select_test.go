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
				Modulo(IntLiteral(1149), IntLiteral(123)),
				Pow(IntLiteral(42), IntLiteral(1)),
			),
			"SELECT (9 + 33), (123 - 81), (14 * 3), (714 / 17), (1149 % 123), (42 ^ 1)",
		},
		{
			Select(
				Plus(
					IntLiteral(7), Times(
						Plus(IntLiteral(1), IntLiteral(10)),
						Plus(IntLiteral(25), IntLiteral(50)),
					),
				),
			),
			"SELECT (7 + ((1 + 10) * (25 + 50)))",
		},
		{
			Select(
				TableColumn("users", "name"),
				TableColumn("users", "email"),
				Divide(TableColumn("users", "height"), IntLiteral(100)),
			),
			`SELECT "name", "email", ("height" / 100) FROM "users"`,
		},
		{
			Select(
				TableColumn("users", "name"),
				TableColumn("animals", "species"),
			),
			`SELECT "name", "species" FROM "users", "animals"`,
		},
		{
			Select(
				AllColumns("users"),
				AllColumns("animals"),
			),
			`SELECT "users".*, "animals".* FROM "users", "animals"`,
		},
		{
			Select(
				TableColumn("users", "name"),
			).OrderBy(
				Descending(TableColumn("users", "height")),
				Ascending(TableColumn("users", "name")),
			),
			`SELECT "name" FROM "users" ORDER BY "height" DESC, "name" ASC`,
		},
		{
			Select(
				TableColumn("users", "name"),
			).OrderBy(
				Ascending(TableColumn("animals", "weight")),
			),
			`SELECT "name" FROM "users", "animals" ORDER BY "weight" ASC`,
		},
		{
			Select(
				TableColumn("users", "name"),
			).OrderBy(
				Descending(AllColumns("users")),
			),
			`SELECT "name" FROM "users" ORDER BY "users".* DESC`,
		},
		{
			Select(
				TableColumn("users", "name"),
			).OrderBy(
				Descending(Divide(IntLiteral(10), IntLiteral(5))),
			),
			`SELECT "name" FROM "users" ORDER BY (10 / 5) DESC`,
		},
	}

	for i, tc := range cases {
		got := tc.query.ToSQL()
		if got != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, got)
		}
	}
}
