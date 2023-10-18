package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq" // driver import
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
)

type PgStorage struct {
	db  *sql.DB
	ctx context.Context
}

// dto to process null values.
type sqlEvent struct {
	ID          string
	Title       string
	DateTime    time.Time
	Description sql.NullString
	Duration    sql.NullString
	RemindTime  sql.NullString

	UserID int
}

var ErrConnectFailed = errors.New("error connecting to db")

const (
	tableName          = "event"
	tableColumnsRead   = "id,title,description,datetime,duration,remind_time,user_id"
	tableColumnsInsert = "title,description,datetime,duration,remind_time,user_id"
)

func (s *PgStorage) Create(event entity.Event) (string, error) {
	err := s.db.QueryRowContext(
		s.ctx,
		fmt.Sprintf("INSERT INTO %s(%s) VALUES($1,$2,$3,$4,$5,$6) RETURNING id", tableName, tableColumnsInsert),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration,
		event.RemindTime,
		strconv.Itoa(event.UserID),
	).Scan(&event.ID)
	if err != nil {
		return "", err
	}

	return event.ID, nil
}

func (s *PgStorage) GetByID(id string) (*entity.Event, error) {
	event := sqlEvent{}

	err := s.db.QueryRowContext(
		s.ctx,
		fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", tableColumnsRead, tableName),
		id,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.DateTime,
		&event.Duration,
		&event.RemindTime,
		&event.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = entity.ErrEventNotFound
		}

		return nil, err
	}

	return s.sqlEventToEvent(&event), nil
}

func (s *PgStorage) GetAll() (*entity.Events, error) {
	events := entity.Events{}

	rows, err := s.db.QueryContext(s.ctx, fmt.Sprintf("SELECT %s FROM %s", tableColumnsRead, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := sqlEvent{}
		err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.DateTime,
			&event.Duration,
			&event.RemindTime,
			&event.UserID,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, *s.sqlEventToEvent(&event))
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func (s *PgStorage) Update(event entity.Event) error {
	_, err := s.db.ExecContext(
		s.ctx,
		fmt.Sprintf(
			"UPDATE %s SET title=$1, description=$2, datetime=$3, duration=$4, remind_time=$5, user_id=$6 WHERE id=$7",
			tableName,
		),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration,
		event.RemindTime,
		event.UserID,
		event.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Delete(id string) error {
	_, err := s.db.ExecContext(s.ctx, fmt.Sprintf("DELETE FROM %s where id=$1", tableName), id)
	if err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) GetForPeriod(start time.Time, end time.Time) (*entity.Events, error) {
	events := entity.Events{}

	rows, err := s.db.QueryContext(
		s.ctx,
		fmt.Sprintf("SELECT %s FROM %s WHERE datetime BETWEEN $1 AND $2", tableColumnsRead, tableName),
		start,
		end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := sqlEvent{}
		err = rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.DateTime,
			&event.Duration,
			&event.RemindTime,
			&event.UserID,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, *s.sqlEventToEvent(&event))
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func (s *PgStorage) GetForTime(t time.Time) (*entity.Event, error) {
	event := sqlEvent{}

	err := s.db.QueryRowContext(
		s.ctx,
		fmt.Sprintf("SELECT %s FROM %s WHERE datetime=$1", tableColumnsRead, tableName), t,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.DateTime,
		&event.Duration,
		&event.RemindTime,
		&event.UserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = entity.ErrEventNotFound
		}

		return nil, err
	}

	return s.sqlEventToEvent(&event), nil
}

func New() *PgStorage {
	return &PgStorage{}
}

func (s *PgStorage) Connect(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return ErrConnectFailed
	}

	db, openErr := sql.Open("postgres", cfg.DB.Dsn)
	if openErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", openErr)
	}

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", pingErr)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *PgStorage) Close(_ context.Context) error {
	closeErr := s.db.Close()
	if closeErr != nil {
		return closeErr
	}

	s.ctx = nil

	return nil
}

func (s *PgStorage) sqlEventToEvent(sqlEvent *sqlEvent) *entity.Event {
	event := entity.Event{}

	event.ID = sqlEvent.ID
	event.Title = sqlEvent.Title
	event.DateTime = sqlEvent.DateTime
	event.UserID = sqlEvent.UserID

	if sqlEvent.Description.Valid {
		event.Description = sqlEvent.Description.String
	}

	if sqlEvent.Duration.Valid {
		event.Duration = sqlEvent.Duration.String
	}

	if sqlEvent.RemindTime.Valid {
		event.RemindTime = sqlEvent.RemindTime.String
	}

	return &event
}
