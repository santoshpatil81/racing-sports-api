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
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error)
	// Get race details returns details of race based on ID
	GetRaceDetails(id int64) (*racing.Race, error)
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

// GetRaceDetails returns the details of a race based on the ID provided
func (r *racesRepo) GetRaceDetails(id int64) (*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	if id <= 0 {
		var race racing.Race
		return &race, err
	}

	query = getRaceQueries()[racesList]

	query, args = r.applyGetRaceDetailsFilter(query, id)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		var race racing.Race
		return &race, err
	}

	var raceDetails racing.Race
	id, meetingId, name, number, visible, advertisedStartTime, status, err := r.scanRace(rows)
	raceDetails.Id = id
	raceDetails.MeetingId = meetingId
	raceDetails.Name = name
	raceDetails.Number = number
	raceDetails.Visible = visible
	raceDetails.AdvertisedStartTime = timestamppb.New(advertisedStartTime)
	raceDetails.Status = status

	if err != nil {
		var race racing.Race
		return &race, err
	}

	return &raceDetails, nil
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

// Search data based on Race ID
func (r *racesRepo) applyGetRaceDetailsFilter(query string, id int64) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if id > 0 {
		clauses = append(clauses, "id="+strconv.FormatInt(id, 10))
		args = append(args, id)

		if len(clauses) != 0 {
			query += " WHERE " + strings.Join(clauses, " AND ")
		}
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

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		race.AdvertisedStartTime = timestamppb.New(advertisedStart)

		races = append(races, &race)
	}

	return races, nil
}

// Scan rows returned by database and create a response object to send to GetRaceDetails
func (m *racesRepo) scanRace(
	rows *sql.Rows,
) (int64, int64, string, int64, bool, time.Time, string, error) {

	var race racing.Race
	var advertisedStart time.Time
	for rows.Next() {
		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &race.Status); err != nil {
			var t time.Time
			if err == sql.ErrNoRows {
				return -1, -1, "", -1, false, t, "", err
			}
			return -1, -1, "", -1, false, t, "", err
		}
	}

	return race.Id, race.MeetingId, race.Name, race.Number, race.Visible, advertisedStart, race.Status, nil
}
