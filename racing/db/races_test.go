package db

import (
	"database/sql"
	"reflect"
	"sync"
	"testing"

	"git.neds.sh/matty/entain/racing/proto/racing"
	_ "github.com/mattn/go-sqlite3"
)

func getTrue() *bool {
	b := true
	return &b
}

func getFalse() *bool {
	b := false
	return &b
}

func Test_racesRepo_applyFilter(t *testing.T) {
	type fields struct {
		db   *sql.DB
		init sync.Once
	}
	type args struct {
		query  string
		filter *racing.ListRacesRequestFilter
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
				query:  "select * from races",
				filter: &racing.ListRacesRequestFilter{},
			},
			want:  "select * from races",
			want1: nil,
		},
		{
			name:   "TEST002: meeting_ids is 6,8",
			fields: fields{},
			args: args{
				query: "select * from races where ",
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
				query: "select * from races where ",
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
				query: "select * from races where ",
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
				query: "select * from races where ",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &racesRepo{
				db:   tt.fields.db,
				init: tt.fields.init,
			}
			got, got1 := r.applyFilter(tt.args.query, tt.args.filter)
			if got != tt.want {
				t.Errorf("racesRepo.applyFilter() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("racesRepo.applyFilter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
