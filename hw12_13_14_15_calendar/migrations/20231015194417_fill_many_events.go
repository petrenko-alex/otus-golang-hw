package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFillManyEvents, downFillManyEvents)
}

func upFillManyEvents(ctx context.Context, tx *sql.Tx) error {
	query := `INSERT INTO event 
    (title, datetime, user_id) 
VALUES 
    ('event 1', '2023-10-15 10:00:00.000000', 1), 
    ('event 2', '2023-10-15 15:00:00.000000', 1), 
    ('event 3', '2023-10-15 23:00:00.000000', 1), 
    ('event 4', '2023-10-16 10:00:00.000000', 1),
    ('event 5', '2023-10-17 14:00:00.000000', 1),
    ('event 6', '2023-10-18 15:00:00.000000', 1),
    ('event 7', '2023-11-02 17:00:00.000000', 1),
    ('event 8', '2023-11-21 19:00:00.000000', 1)
;`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downFillManyEvents(ctx context.Context, tx *sql.Tx) error {
	query := `DELETE FROM event 
    WHERE title IN ('event 1', 'event 2', 'event 3', 'event 4', 'event 5', 'event 6', 'event 7', 'event 8') ;`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
