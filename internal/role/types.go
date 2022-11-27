package role

import "github.com/google/uuid"

/*
 applicaiton should capture all the information needed to provision and
 manage an permission on the bestir network
*/
type Role struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}
