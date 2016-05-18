package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/leocassarani/psql"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := psql.Select(
		psql.TableColumn("users", "name"),
		psql.TableColumn("users", "email"),
	).OrderBy(
		psql.Descending(psql.TableColumn("users", "height")),
	)

	rows, err := db.Query(query.ToSQL())
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, email string
		if err := rows.Scan(&name, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", name, email)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
