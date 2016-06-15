package psql

import (
	"reflect"
	"testing"
)

func TestSelectQuerySQL(t *testing.T) {
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
				Avg(TableColumn("users", "age")),
				Min(TableColumn("animals", "weight")),
				Max(TableColumn("users", "height")),
				Sum(TableColumn("animals", "paws")),
			),
			`SELECT AVG("age"), MIN("weight"), MAX("height"), SUM("paws") FROM "users", "animals"`,
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
		{
			Select(
				Avg(TableColumn("users", "height")),
				TableColumn("users", "name"),
			).GroupBy(
				TableColumn("users", "name"),
			),
			`SELECT AVG("height"), "name" FROM "users" GROUP BY "name"`,
		},
		{
			Select(
				IsNull(TableColumn("users", "name")),
				IsNotNull(TableColumn("users", "email")),
			),
			`SELECT "name" IS NULL, "email" IS NOT NULL FROM "users"`,
		},
		{
			Select(
				Eq(IntLiteral(42), IntLiteral(42)),
				NotEq(IntLiteral(42), IntLiteral(21)),
				LessThan(IntLiteral(21), IntLiteral(42)),
				LessThanOrEq(IntLiteral(42), IntLiteral(42)),
				GreaterThan(IntLiteral(42), IntLiteral(21)),
				GreaterThanOrEq(IntLiteral(42), IntLiteral(42)),
			),
			`SELECT (42 = 42), (42 <> 21), (21 < 42), (42 <= 42), (42 > 21), (42 >= 42)`,
		},
		{
			Select(
				TableColumn("users", "name"),
				GreaterThanOrEq(TableColumn("users", "height"), IntLiteral(180)),
			),
			`SELECT "name", ("height" >= 180) FROM "users"`,
		},
		{
			Select(
				TableColumn("users", "height"),
			).Where(
				Eq(TableColumn("users", "name"), StringParam()),
				NotEq(TableColumn("users", "city"), StringParam()),
				IsNotNull(TableColumn("users", "height")),
			),
			`SELECT "height" FROM "users" WHERE ("name" = $1::text) AND ("city" <> $2::text) AND "height" IS NOT NULL`,
		},
		{
			Select(
				Now(),
			),
			`SELECT now()`,
		},
		{
			Select(
				DatePart(DayField, TableColumn("animals", "birthdate")),
				DatePart(MonthField, TableColumn("animals", "birthdate")),
				DatePart(YearField, Now()),
			),
			`SELECT date_part('day', "birthdate"), date_part('month', "birthdate"), date_part('year', now()) FROM "animals"`,
		},
		{
			Select(
				DateTrunc(DayPrecision, TableColumn("users", "signup_date")),
				Max(TableColumn("users", "height")),
			).GroupBy(
				DateTrunc(DayPrecision, TableColumn("users", "signup_date")),
			),
			`SELECT date_trunc('day', "signup_date"), MAX("height") FROM "users" GROUP BY date_trunc('day', "signup_date")`,
		},
	}

	for i, tc := range cases {
		got := tc.query.ToSQL()
		if got != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, got)
		}
	}
}

func TestSelectQueryBindings(t *testing.T) {
	cases := []struct {
		query  SelectQuery
		inputs []interface{}

		sql      string
		bindings []interface{}
	}{
		{
			Select(
				StringLiteral("Hello"),
				StringLiteral("World"),
			),
			[]interface{}{},

			`SELECT $1::text, $2::text`,
			[]interface{}{"Hello", "World"},
		},
		{
			Select(
				StringLiteral("Hello"),
				StringParam(),
			),
			[]interface{}{"Joe"},

			`SELECT $1::text, $2::text`,
			[]interface{}{"Hello", "Joe"},
		},
	}

	for i, tc := range cases {
		sql := tc.query.ToSQL()
		if sql != tc.sql {
			t.Errorf("test case %d: expected %q, got %q", i+1, tc.sql, sql)
		}

		bindings := tc.query.Bindings(tc.inputs...)
		if !reflect.DeepEqual(bindings, tc.bindings) {
			t.Errorf("test case %d: expected %v, got %v", i+1, tc.bindings, bindings)
		}
	}
}

func TestSelectQueryIdempotence(t *testing.T) {
	cases := []SelectQuery{
		Select(
			TableColumn("users", "name"),
			TableColumn("users", "email"),
		),

		Select(
			StringLiteral("Hello"),
			StringLiteral("World"),
		),

		Select(
			StringLiteral("Hello"),
			StringParam(),
		).Where(
			Eq(TableColumn("users", "name"), StringParam()),
		),
	}

	for i, query := range cases {
		// Check that calling ToSQL() repeatedly on the same query returns the same string.
		want, got := query.ToSQL(), query.ToSQL()
		if want != got {
			t.Errorf("test case %d: expected %q, got %q", i+1, want, got)
		}
	}
}
