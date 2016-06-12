# psql [![Build Status](https://travis-ci.org/leocassarani/psql.svg?branch=master)](https://travis-ci.org/leocassarani/psql)

psql is a PostgreSQL query builder for Go.

It is intended to be an addition to, rather than a replacement for, the
[database/sql](https://golang.org/pkg/database/sql/) package in the
standard library.

psql offers a convenient API for generating safe SQL queries at runtime.
It uses Go's type system to ensure both identifiers (such as table and
column names) and parameters (user input) are correctly escaped.

psql is not an ORM; however, it may be used as a more lightweight
alternative to a full-blown ORM.

## Usage

Use psql to compose your query, then call `ToSQL()` to convert it into
an SQL string that you can pass into `DB.Query()` or `DB.QueryRow()`.

```go
query := psql.Select(
  psql.TableColumn("users", "name"),
  psql.TableColumn("users", "email"),
).Where(
  psql.NotEq(psql.TableColumn("users", "name"), psql.StringParam()),
).OrderBy(
  psql.Descending(psql.TableColumn("users", "height")),
)

// SELECT "name", "email" FROM "users" WHERE "name" <> 'Joe' ORDER BY "height"
fmt.Println(query.ToSQL(), query.Bindings("Joe")...)

// Run the query on a database connection.
db.Query(query.ToSQL())
```
