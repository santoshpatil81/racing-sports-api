package db

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

func getTrue() *bool {
	t := true
	return &t
}

func getFalse() *bool {
	t := false
	return &t
}

func getAdvertisedStartTimeFieldName() string {
	s := "advertised_start_time"
	return s
}

func getStatusFieldName() string {
	s := "status"
	return s
}

func Test_racesRepo_applyFilter(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		query     string
		filter    *racing.ListRacesRequestFilter
		configVal *Configuration
	}
	if os.Getenv("CONFIG_FILE") == "" {
		os.Setenv("CONFIG_FILE", "../config.json")
	}
	advertisedStartTime := getAdvertisedStartTimeFieldName()
	status := getStatusFieldName()
	configValues := getConfigValue(os.Getenv("CONFIG_FILE"))
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "TEST001: No filter is passed",
			fields: fields{},
			args: args{
				query:     "select * from races",
				filter:    &racing.ListRacesRequestFilter{},
				configVal: configValues,
			},
			want:  "select * from races",
			want1: nil,
		},
		{
			name:   "TEST002: meeting_ids is 6,8",
			fields: fields{},
			args: args{
				query:     "select * from races where ",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					MeetingIds: []int64{6, 8},
				},
			},
			want: "select * from races where  WHERE meeting_id IN (?,?)",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8))
			}(),
		},
		{
			name:   "TEST003: Visible is true",
			fields: fields{},
			args: args{
				query:     "select * from races where ",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					Visible: getTrue(),
				},
			},
			want: "select * from races where  WHERE visible=true",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, true)
			}(),
		},
		{
			name:   "TEST004: Visible is false",
			fields: fields{},
			args: args{
				query:     "select * from races where ",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					Visible: getFalse(),
				},
			},
			want: "select * from races where  WHERE visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, false)
			}(),
		},
		{
			name:   "TEST005: Visible is false and meeting_ids is 3",
			fields: fields{},
			args: args{
				query:     "select * from races where ",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					MeetingIds: []int64{3},
					Visible:    getFalse(),
				},
			},
			want: "select * from races where  WHERE meeting_id IN (?) AND visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST006: sort by advertised start time",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
				},
			},
			want: "select * from races order by advertised_start_time",
		},
		{
			name:   "TEST007: sort by start time is true and meeting Ids is 6,8 and visible is true",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
					MeetingIds:  []int64{6, 8},
					Visible:     getTrue(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?,?) AND visible=true order by advertised_start_time",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), true)
			}(),
		},
		{
			name:   "TEST008: sort by start time is empty and meeting Ids is 3 and visible is true",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: nil,
					MeetingIds:  []int64{3},
					Visible:     getTrue(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?) AND visible=true",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), true)
			}(),
		},
		{
			name:   "TEST009: sort by start time is true and meeting Ids is 6,8 and visible is false",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
					MeetingIds:  []int64{6, 8},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?,?) AND visible=false order by advertised_start_time",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), false)
			}(),
		},
		{
			name:   "TEST010: sort by start time is empty and meeting Ids is 3 and visible is false",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: nil,
					MeetingIds:  []int64{3},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?) AND visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST011: sort by advertised start time and order by ascending",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
					MeetingIds:  []int64{3},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?) AND visible=false order by advertised_start_time",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST012: sort by advertised start time and order by descending",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
					MeetingIds:  []int64{3},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?) AND visible=false order by advertised_start_time",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST013: sort by start time and order by ascending and meeting Ids is 6,8 and visible is false",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &advertisedStartTime,
					MeetingIds:  []int64{6, 8},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?,?) AND visible=false order by advertised_start_time",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), false)
			}(),
		},
		{
			name:   "TEST014: sort by status",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &status,
				},
			},
			want: "select * from races order by status",
		},
		{
			name:   "TEST015: sort by status and order by descending",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &status,
					MeetingIds:  []int64{3},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?) AND visible=false order by status",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST016: sort by status and order by ascending and meeting Ids is 6,8 and visible is false",
			fields: fields{},
			args: args{
				query:     "select * from races",
				configVal: configValues,
				filter: &racing.ListRacesRequestFilter{
					SortByField: &status,
					MeetingIds:  []int64{6, 8},
					Visible:     getFalse(),
				},
			},
			want: "select * from races WHERE meeting_id IN (?,?) AND visible=false order by status",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), false)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
				db: tt.fields.db,
			}
			got, got1 := r.applyFilter(tt.args.query, tt.args.filter, tt.args.configVal)
			if got != tt.want {
				t.Errorf("racesRepo.applyFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("racesRepo.applyFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_racesRepo_applyGetRaceDetailsFilter(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		query string
		id    int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  []interface{}
	}{
		{
			name:   "TEST001: No id is passed",
			fields: fields{},
			args: args{
				query: "select * from races",
				id:    int64(0),
			},
			want:  "select * from races",
			want1: nil,
		},
		{
			name:   "TEST002: id is 1",
			fields: fields{},
			args: args{
				query: "select * from races",
				id:    int64(1),
			},
			want: "select * from races WHERE id=1",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(1))
			}(),
		},
		{
			name:   "TEST003: id is -100",
			fields: fields{},
			args: args{
				query: "select * from races",
				id:    int64(-100),
			},
			want:  "select * from races",
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
				db: tt.fields.db,
			}
			got, got1 := r.applyGetRaceDetailsFilter(tt.args.query, tt.args.id)
			if got != tt.want {
				t.Errorf("racesRepo.applyGetRaceDetailsFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("racesRepo.applyGetRaceDetailsFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
