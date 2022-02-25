package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// SportsRepo provides repository access to sports events.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of sports events.
	List(filter *sports.ListSportsRequestFilter) ([]*sports.SportsEvent, error)
	GetSportsEventDetails(id int64) ([]*sports.SportsEvent, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sports repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy sports events.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) GetSportsEventDetails(id int64) ([]*sports.SportsEvent, error) {
	var (
		err   error
		query string
		args  []interface{}
	)
	if id == 0 {
		return nil, err
	}
	query = getSportsQueries()[eventsList]

	query, args = s.applyGetEventDetailsFilter(query, id)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

func (s *sportsRepo) List(filter *sports.ListSportsRequestFilter) ([]*sports.SportsEvent, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getSportsQueries()[eventsList]

	query, args = s.applyFilter(query, filter)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListSportsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	order_clauses := ""
	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	if filter.Visible != nil {
		clauses = append(clauses, "visible="+strconv.FormatBool(*filter.Visible))
		args = append(args, *filter.Visible)
	}

	if filter.OrderByStartTime != nil {
		if *filter.OrderByStartTime {
			order_clauses += " order by advertised_start_time asc"
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	if len(order_clauses) != 0 {
		query += order_clauses
	}

	return query, args
}

func (s *sportsRepo) applyGetEventDetailsFilter(query string, id int64) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if id <= 0 {
		return query, args
	}

	clauses = append(clauses, "id="+strconv.FormatInt(id, 10))
	args = append(args, id)

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (s *sportsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.SportsEvent, error) {
	var events []*sports.SportsEvent

	for rows.Next() {
		var event sports.SportsEvent
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.MeetingId, &event.Name, &event.Number, &event.Visible, &advertisedStart, &event.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		event.AdvertisedStartTime = timestamppb.New(advertisedStart)

		events = append(events, &event)
	}

	return events, nil
}
