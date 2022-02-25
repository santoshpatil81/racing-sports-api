package db

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
    			case when strftime('%s', advertised_start_time) <= strftime('%s', 'now') then 'OPEN'
    				 when strftime('%s', advertised_start_time) > strftime('%s', 'now') then 'CLOSED'
    			end as status
			FROM races
		`,
	}
}
