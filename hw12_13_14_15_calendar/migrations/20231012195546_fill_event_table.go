package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFillEventTable, downFillEventTable)
}

func upFillEventTable(ctx context.Context, tx *sql.Tx) error {
	query := `INSERT INTO event 
    (title, description, datetime, duration, remind_time, user_id) 
VALUES 
    ('golang meetup', 'first ever golang meetup in Russia', '2023-11-10 22:00:00.000000', '02:00:00', '01:00:00', 1), 
    ('car repair', NULL, '2023-12-12 22:19:05.000000', '01:00:00', '03:00:00', 1),
    ('swimming pool', NULL, '2023-11-17 22:19:55.000000', '00:45:00', '00:15:00', 2);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downFillEventTable(ctx context.Context, tx *sql.Tx) error {
	query := "TRUNCATE TABLE event"

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
