package permission

import "github.com/google/uuid"

/*
 applicaiton should capture all the information needed to provision and
 manage an permission on the bestir network
*/
type permission struct {
	ID   uuid.UUID `db:"id" json:"id"`
	Name string    `db:"name" json:"name"`
}

type Incomingpermission struct {
	Name string `json:"name" required:"true"`
	// IdempotencyKey null.String `json:"-" db:"idempotency_key"`
}

type permissionGroup struct {
}
