package services

import (
	"context"
	"github.com/go-pg/pg/v10"
	"log"
)

type PostgresLogger struct{}

func (d PostgresLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d PostgresLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	query, _ := q.FormattedQuery()
	log.Printf("%s\n", query)
	return nil
}
