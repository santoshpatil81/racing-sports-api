package db

import (
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

type Configuration struct {
	OrderBy []string
	SortBy  []string
}

// getConfigValue reads config values from config value
func getConfigValue(configFileName string) *Configuration {

	file, _ := os.Open(configFileName)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return nil
	}

	return &configuration
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

// List returns a list of matching races based on the criteria specified in the body of the request
func (r *racesRepo) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	configVal := getConfigValue(os.Getenv("CONFIG_FILE"))

	query, args = r.applyFilter(query, filter, configVal)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

// Compare the provided options with the options specifed in the config file
func CheckOptionValidity(options []string, field string) int {

	for key, val := range options {
		if val == field {
			return key
		}
	}

	return -1
}

// Process the filter criteria to return matching races
func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter, configVal *Configuration) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	orderClauses := ""
	if filter == nil {
		return query, args
	}

	// Process "meeting_ids" option in filter. One or more meeting IDs are specified in an array
	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// Process "visible" option in filter
	if filter.Visible != nil {
		clauses = append(clauses, "visible="+strconv.FormatBool(*filter.Visible))
		args = append(args, *filter.Visible)
	}

	// Process "sort_by_field" option specific in the filter
	if filter.SortByField != nil {
		sortFieldIndex := CheckOptionValidity(configVal.SortBy, *filter.SortByField)
		if sortFieldIndex > -1 {
			orderClauses += " order by " + configVal.SortBy[sortFieldIndex]
			// Process "order_by" option specific in the filter
			// "order_by" works in conjunctions with the "sort_by_field" but is optional
			// if "sort_by_field" option in filter is empty then "order_by" has no effect
			if filter.OrderBy != nil {
				orderFieldIndex := CheckOptionValidity(configVal.OrderBy, *filter.OrderBy)
				if orderFieldIndex > -1 {
					orderClauses += " " + configVal.OrderBy[orderFieldIndex]
				}
			}
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	if len(orderClauses) != 0 {
		query += orderClauses
	}

	return query, args
}

// Scan rows returned by database and create a response object to send to ListRaces
func (m *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		race.AdvertisedStartTime = ts

		races = append(races, &race)
	}

	return races, nil
}
