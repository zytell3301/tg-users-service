package domain

import (
	"github.com/gocql/gocql"
	"time"
)

type User struct {
	Id            gocql.UUID
	Username      string
	Phone         string
	Online_status bool
	Created_at    time.Time
}
