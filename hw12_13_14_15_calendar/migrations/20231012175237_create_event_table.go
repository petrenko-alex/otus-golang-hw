package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateEventTable, downCreateEventTable)
}

func upCreateEventTable(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE event(
    id          uuid default gen_random_uuid() not null primary key,
    title       varchar(255)                   not null,
    description text,
    datetime    timestamp                      not null,
    duration    varchar(255),
    remind_time varchar(255),
    user_id     integer                        not null
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downCreateEventTable(ctx context.Context, tx *sql.Tx) error {
	query := "DROP TABLE IF EXISTS event"

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
