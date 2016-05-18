# psql

psql is a PostgreSQL query builder for Go.

It is intended to be an addition to, rather than a replacement for, the
[database/sql](https://golang.org/pkg/database/sql/) package in the
standard library.

psql offers a convenient API for generating safe SQL queries at runtime.
It uses Go's type system to ensure both identifiers (such as table and
column names) and parameters (user input) are correctly escaped.

psql is not an ORM; however, it may be used as a more lightweight
alternative to a full-blown ORM.
