package db

import (
	"database/sql"
	"reflect"
	"testing"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

func getTrue() *bool {
	b := true
	return &b
}

func getFalse() *bool {
	b := false
	return &b
}

func Test_sportsRepo_applyFilter(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		query  string
		filter *sports.ListSportsRequestFilter
	}
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
				query:  "select * from sports",
				filter: &sports.ListSportsRequestFilter{},
			},
			want:  "select * from sports",
			want1: nil,
		},
		{
			name:   "TEST002: meeting_ids is 6,8",
			fields: fields{},
			args: args{
				query: "select * from sports where ",
				filter: &sports.ListSportsRequestFilter{
					MeetingIds: []int64{6, 8},
				},
			},
			want: "select * from sports where  WHERE meeting_id IN (?,?)",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8))
			}(),
		},
		{
			name:   "TEST003: Visible is true",
			fields: fields{},
			args: args{
				query: "select * from sports where ",
				filter: &sports.ListSportsRequestFilter{
					Visible: getTrue(),
				},
			},
			want: "select * from sports where  WHERE visible=true",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, true)
			}(),
		},
		{
			name:   "TEST004: Visible is false",
			fields: fields{},
			args: args{
				query: "select * from sports where ",
				filter: &sports.ListSportsRequestFilter{
					Visible: getFalse(),
				},
			},
			want: "select * from sports where  WHERE visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, false)
			}(),
		},
		{
			name:   "TEST005: Visible is false and meeting_ids is 3",
			fields: fields{},
			args: args{
				query: "select * from sports where ",
				filter: &sports.ListSportsRequestFilter{
					MeetingIds: []int64{3},
					Visible:    getFalse(),
				},
			},
			want: "select * from sports where  WHERE meeting_id IN (?) AND visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
		{
			name:   "TEST006: order by start time is true",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getTrue(),
				},
			},
			want: "select * from sports order by advertised_start_time asc",
		},
		{
			name:   "TEST006: order by start time is false",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getFalse(),
				},
			},
			want: "select * from sports",
		},
		{
			name:   "TEST007: order by start time is true and meeting Ids is 6.8 and visible is true",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getTrue(),
					MeetingIds:       []int64{6, 8},
					Visible:          getTrue(),
				},
			},
			want: "select * from sports WHERE meeting_id IN (?,?) AND visible=true order by advertised_start_time asc",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), true)
			}(),
		},
		{
			name:   "TEST008: order by start time is false and meeting Ids is 3 and visible is true",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getFalse(),
					MeetingIds:       []int64{3},
					Visible:          getTrue(),
				},
			},
			want: "select * from sports WHERE meeting_id IN (?) AND visible=true",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), true)
			}(),
		},
		{
			name:   "TEST009: order by start time is true and meeting Ids is 6.8 and visible is false",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getTrue(),
					MeetingIds:       []int64{6, 8},
					Visible:          getFalse(),
				},
			},
			want: "select * from sports WHERE meeting_id IN (?,?) AND visible=false order by advertised_start_time asc",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(6), int64(8), false)
			}(),
		},
		{
			name:   "TEST010: order by start time is false and meeting Ids is 3 and visible is false",
			fields: fields{},
			args: args{
				query: "select * from sports",
				filter: &sports.ListSportsRequestFilter{
					OrderByStartTime: getFalse(),
					MeetingIds:       []int64{3},
					Visible:          getFalse(),
				},
			},
			want: "select * from sports WHERE meeting_id IN (?) AND visible=false",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(3), false)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &sportsRepo{
				db: tt.fields.db,
			}
			got, got1 := r.applyFilter(tt.args.query, tt.args.filter)
			if got != tt.want {
				t.Errorf("sportsRepo.applyFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sportsRepo.applyFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_sportsRepo_applyGetRaceDetailsFilter(t *testing.T) {
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
				query: "select * from sports",
				id:    int64(0),
			},
			want:  "select * from sports",
			want1: nil,
		},
		{
			name:   "TEST002: id is 1",
			fields: fields{},
			args: args{
				query: "select * from sports",
				id:    int64(1),
			},
			want: "select * from sports WHERE id=1",
			want1: func() []interface{} {
				var ret []interface{}
				return append(ret, int64(1))
			}(),
		},
		{
			name:   "TEST003: id is -100",
			fields: fields{},
			args: args{
				query: "select * from sports",
				id:    int64(-100),
			},
			want:  "select * from sports",
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &sportsRepo{
				db: tt.fields.db,
			}
			got, got1 := r.applyGetEventDetailsFilter(tt.args.query, tt.args.id)
			if got != tt.want {
				t.Errorf("sportsRepo.applyGetEventDetailsFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("sportsRepo.applyGetEventDetailsFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
