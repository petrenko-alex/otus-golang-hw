package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/entity"
	"strconv"
	"time"
)

// TODO: Query with context?

type PgStorage struct {
	db *sql.DB
}

var (
	ErrNotFound      = errors.New("event with specified id was not found")
	ErrConnectFailed = errors.New("error connecting to db")
)

const (
	tableName          = "event"
	tableColumnsRead   = "id,title,description,datetime,duration,remind_time,user_id"
	tableColumnsInsert = "title,description,datetime,duration,remind_time,user_id"
)

func (s *PgStorage) Create(event entity.Event) (string, error) {
	err := s.db.QueryRow(
		fmt.Sprintf("INSERT INTO %s(%s) VALUES($1,$2,$3,$4,$5,$6) RETURNING id", tableName, tableColumnsInsert),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration,
		event.RemindTime,
		strconv.Itoa(event.UserId),
	).Scan(&event.ID)

	if err != nil {
		return "", err
	}

	return event.ID, nil
}

func (s *PgStorage) ReadOne(id string) (*entity.Event, error) {
	event := entity.Event{}

	err := s.db.QueryRow(
		fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", tableColumnsRead, tableName), id,
	).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.DateTime,
		&event.Duration,
		&event.RemindTime,
		&event.UserId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNotFound
		}

		return nil, err
	}

	return &event, nil
}

func (s *PgStorage) ReadAll() (*entity.Events, error) {
	events := entity.Events{}

	rows, err := s.db.Query(fmt.Sprintf("SELECT %s FROM %s", tableColumnsRead, tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := entity.Event{}
		err = rows.Scan(&event.ID, &event.Title, &event.Description, &event.DateTime, &event.Duration, &event.RemindTime, &event.UserId)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &events, nil
}

func (s *PgStorage) Update(event entity.Event) error {
	_, err := s.db.Exec(
		fmt.Sprintf("UPDATE %s SET title=$1, description=$2, datetime=$3, duration=$4, remind_time=$5, user_id=$6 WHERE id=$7", tableName),
		event.Title,
		event.Description,
		event.DateTime.Format(time.RFC3339),
		event.Duration,
		event.RemindTime,
		event.UserId,
		event.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) Delete(id string) error {
	_, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s where id=$1", tableName), id)
	if err != nil {
		return err
	}

	return nil
}

func New() *PgStorage {
	return &PgStorage{}
}

func (s *PgStorage) Connect(ctx context.Context) error {
	cfg, ok := ctx.Value("config").(*config.Config)
	if !ok {
		return ErrConnectFailed
	}

	db, openErr := sql.Open("postgres", cfg.Db.Dsn)
	if openErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", openErr)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", pingErr)
	}

	s.db = db

	return nil
}

func (s *PgStorage) Close(ctx context.Context) error {
	closeErr := s.db.Close()
	if closeErr != nil {
		return closeErr
	}

	return nil
}
