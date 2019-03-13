# QueryBuilder

[QueryBuilder](https://github.com/squioc/querybuilder) aims to help to generate dynamic sql queries for `database/sql` and [github.com/jmoiron/sqlx](http://jmoiron.github.io/sqlx/).

## Documentation

API documentation can be find on [godoc.org](https://godoc.org/github.com/squioc/querybuilder)

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/squioc/querybuilder"

    _ "github.com/lib/pq"
    "github.com/jmoiron/sqlx" 
)

type Row struct {
    Name string `db:"name"`
    Id  int `db:"id"`
}

func main() {
    db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

    queryBuilder := builder.NewQueryBuilder("select * from table")
    queryBuilder.appendCriterion("name", "myRow")
    queryBuilder.appendCriterion("id", 6)
    query, values = queryBuilder.Build()

    fmt.Printf("query: %s", query)
    // query: select * from table where name = $1 and id = $2
    fmt.Printf("values: %+v", values)
    // values: [myRow 6]


    err = db.Get(&row, query, values...)
}
```
