package database

import "time"

// Our MySQL driver only keeps track up to the second
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
